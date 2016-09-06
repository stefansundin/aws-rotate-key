package main

import (
	"flag"
	"fmt"
	"os/user"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
)

func main() {
	var yesFlag bool
	var profileFlag string
	flag.BoolVar(&yesFlag, "y", false, `Automatic "yes" to prompts.`)
	flag.StringVar(&profileFlag, "profile", "default", "The profile to use.")
	flag.Parse()

	// Get credentials
	usr, _ := user.Current()
	credentialsPath := fmt.Sprintf("%s/.aws/credentials", usr.HomeDir)
	credentialsProvider := credentials.NewSharedCredentials(credentialsPath, profileFlag)
	creds, err := credentialsProvider.Get()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("Using access key %s from profile \"%s\".\n", creds.AccessKeyID, profileFlag)

	// Create session
	sess, err := session.NewSession(&aws.Config{Credentials: credentialsProvider})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// sts get-caller-identity
	stsClient := sts.New(sess)
	respGetCallerIdentity, err := stsClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("Your user arn is: %s\n", *respGetCallerIdentity.Arn)

	// iam list-access-keys
	// If the UserName field is not specified, the UserName is determined implicitly based on the AWS access key ID used to sign the request.
	iamClient := iam.New(sess)
	respListAccessKeys, err := iamClient.ListAccessKeys(&iam.ListAccessKeysInput{})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Print key information
	fmt.Printf("You have %d access key%v associated with your user:\n", len(respListAccessKeys.AccessKeyMetadata), pluralize(len(respListAccessKeys.AccessKeyMetadata)))
	for _, key := range respListAccessKeys.AccessKeyMetadata {
		fmt.Printf("- %s (created %s)\n", *key.AccessKeyId, key.CreateDate)
	}
	fmt.Println()
}

func pluralize(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}
