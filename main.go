package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
)

const version = "1.0.4"

func main() {
	var yesFlag bool
	var profileFlag string
	var versionFlag bool
	var deleteFlag bool
	flag.BoolVar(&yesFlag, "y", false, `Automatic "yes" to prompts.`)
	flag.BoolVar(&deleteFlag, "d", false, "Delete old key without deactivation.")
	flag.StringVar(&profileFlag, "profile", "default", "The profile to use.")
	flag.BoolVar(&versionFlag, "version", false, "Print version number ("+version+")")
	flag.Parse()

	if versionFlag {
		fmt.Println(version)
		os.Exit(0)
	}

	// Get credentials
	usr, _ := user.Current()

	credentialsPath := os.Getenv("AWS_SHARED_CREDENTIALS_FILE")
	if len(credentialsPath) == 0 {
		credentialsPath = fmt.Sprintf("%s/.aws/credentials", usr.HomeDir)
	}

	credentialsProvider := credentials.NewSharedCredentials(credentialsPath, profileFlag)
	creds, err := credentialsProvider.Get()
	check(err)
	fmt.Printf("Using access key %s from profile \"%s\".\n", creds.AccessKeyID, profileFlag)

	// Create session
	sess, err := session.NewSession(&aws.Config{Credentials: credentialsProvider})
	check(err)

	// sts get-caller-identity
	stsClient := sts.New(sess)
	respGetCallerIdentity, err := stsClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		fmt.Println("Error getting caller identity. Is the key disabled?")
		fmt.Println()
		check(err)
	}
	fmt.Printf("Your user arn is: %s\n\n", *respGetCallerIdentity.Arn)

	// iam list-access-keys
	// If the UserName field is not specified, the UserName is determined implicitly based on the AWS access key ID used to sign the request.
	iamClient := iam.New(sess)
	respListAccessKeys, err := iamClient.ListAccessKeys(&iam.ListAccessKeysInput{})
	check(err)

	// Print key information
	fmt.Printf("You have %d access key%v associated with your user:\n", len(respListAccessKeys.AccessKeyMetadata), pluralize(len(respListAccessKeys.AccessKeyMetadata)))
	for _, key := range respListAccessKeys.AccessKeyMetadata {
		respAccessKeyLastUsed, err2 := iamClient.GetAccessKeyLastUsed(&iam.GetAccessKeyLastUsedInput{
			AccessKeyId: key.AccessKeyId,
		})
		check(err2)
		if respAccessKeyLastUsed.AccessKeyLastUsed.LastUsedDate == nil {
			fmt.Printf("- %s (%s, created %s, never used)\n", *key.AccessKeyId, *key.Status, key.CreateDate)
		} else {
			fmt.Printf("- %s (%s, created %s, last used %s for service %s in %s)\n", *key.AccessKeyId, *key.Status, key.CreateDate, respAccessKeyLastUsed.AccessKeyLastUsed.LastUsedDate, *respAccessKeyLastUsed.AccessKeyLastUsed.ServiceName, *respAccessKeyLastUsed.AccessKeyLastUsed.Region)
		}
	}
	fmt.Println()

	if len(respListAccessKeys.AccessKeyMetadata) == 2 {
		keyIndex := 0
		if *respListAccessKeys.AccessKeyMetadata[0].AccessKeyId == creds.AccessKeyID {
			keyIndex = 1
		}

		if yesFlag == false {
			fmt.Println("You have two access keys, which is the max number of access keys.")
			fmt.Printf("Do you want to delete %s and create a new key? [yN] ", *respListAccessKeys.AccessKeyMetadata[keyIndex].AccessKeyId)
			if *respListAccessKeys.AccessKeyMetadata[keyIndex].Status == "Active" {
				fmt.Printf("\nWARNING: This key is currently Active! ")
			}
			reader := bufio.NewReader(os.Stdin)
			yn, err2 := reader.ReadString('\n')
			check(err2)
			if yn[0] != 'y' && yn[0] != 'Y' {
				return
			}
		}

		_, err2 := iamClient.DeleteAccessKey(&iam.DeleteAccessKeyInput{
			AccessKeyId: respListAccessKeys.AccessKeyMetadata[keyIndex].AccessKeyId,
		})
		check(err2)
		fmt.Printf("Deleted access key %s.\n", *respListAccessKeys.AccessKeyMetadata[keyIndex].AccessKeyId)
	} else if yesFlag == false {
		fmt.Printf("Do you want to create a new key and deactivate %s? [yN] ", *respListAccessKeys.AccessKeyMetadata[0].AccessKeyId)
		reader := bufio.NewReader(os.Stdin)
		yn, err2 := reader.ReadString('\n')
		check(err2)
		if yn[0] != 'y' && yn[0] != 'Y' {
			return
		}
	}

	// Create the new access key
	// If you do not specify a user name, IAM determines the user name implicitly based on the AWS access key ID signing the request.
	respCreateAccessKey, err := iamClient.CreateAccessKey(&iam.CreateAccessKeyInput{})
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	fmt.Printf("Created access key %s.\n", *respCreateAccessKey.AccessKey.AccessKeyId)

	// Read credentials file
	bytes, err := ioutil.ReadFile(credentialsPath)
	check(err)
	text := string(bytes)

	// Replace key pair in credentials file
	// This search & replace does not limit itself to the specified profile, which may be useful if the user is using the same key in multiple profiles
	re := regexp.MustCompile(fmt.Sprintf(`(?m)^aws_access_key_id ?= ?%s`, regexp.QuoteMeta(creds.AccessKeyID)))
	text = re.ReplaceAllString(text, `aws_access_key_id=`+*respCreateAccessKey.AccessKey.AccessKeyId)
	re = regexp.MustCompile(fmt.Sprintf(`(?m)^aws_secret_access_key ?= ?%s`, regexp.QuoteMeta(creds.SecretAccessKey)))
	text = re.ReplaceAllString(text, `aws_secret_access_key=`+*respCreateAccessKey.AccessKey.SecretAccessKey)

	// Verify that the regexp actually replaced something
	if !strings.Contains(text, *respCreateAccessKey.AccessKey.AccessKeyId) || !strings.Contains(text, *respCreateAccessKey.AccessKey.SecretAccessKey) {
		fmt.Println("Failed to replace old access key. Aborting.")
		fmt.Printf("Please verify that the file %s is formatted correctly.\n", credentialsPath)
		// Delete the key we created
		_, err2 := iamClient.DeleteAccessKey(&iam.DeleteAccessKeyInput{
			AccessKeyId: respCreateAccessKey.AccessKey.AccessKeyId,
		})
		check(err2)
		fmt.Printf("Deleted access key %s.\n", *respCreateAccessKey.AccessKey.AccessKeyId)
		os.Exit(1)
	}

	// Write new file
	err = ioutil.WriteFile(credentialsPath, []byte(text), 0600)
	check(err)
	fmt.Printf("Wrote new key pair to %s\n", credentialsPath)

	// Deleting the key if flag is set, otherwise only deactivating
	if deleteFlag {
		_, err := iamClient.DeleteAccessKey(&iam.DeleteAccessKeyInput{
			AccessKeyId: &creds.AccessKeyID,
		})
		check(err)
		fmt.Printf("Deleted old access key %s.\n", creds.AccessKeyID)
	} else {
		_, err = iamClient.UpdateAccessKey(&iam.UpdateAccessKeyInput{
			AccessKeyId: &creds.AccessKeyID,
			Status:      aws.String("Inactive"),
		})
		check(err)
		fmt.Printf("Deactivated old access key %s.\n", creds.AccessKeyID)
		fmt.Println("Please make sure this key is not used elsewhere.")
	}
}

func pluralize(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
