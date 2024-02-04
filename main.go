package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/user"
	"regexp"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
)

const version = "1.1.0"

var defaultProfile = "default"

func init() {
	// Respect AWS_PROFILE if it is set
	if v, ok := os.LookupEnv("AWS_PROFILE"); ok {
		defaultProfile = v
	}
}

func main() {
	var yesFlag bool
	var mfaFlag bool
	var profileFlag string
	var authProfileFlag string
	var mfaSerialNumber string
	var versionFlag bool
	var deleteFlag bool
	flag.BoolVar(&yesFlag, "y", false, `Automatic "yes" to prompts.`)
	flag.BoolVar(&mfaFlag, "mfa", false, "Use MFA.")
	flag.BoolVar(&deleteFlag, "d", false, "Delete the old key instead of deactivating it.")
	flag.StringVar(&profileFlag, "profile", defaultProfile, "The profile to use.")
	flag.StringVar(&authProfileFlag, "auth-profile", "", "Use a different profile when calling AWS.")
	flag.StringVar(&mfaSerialNumber, "mfa-serial-number", "", "Specify the MFA device to use. (optional)")
	flag.BoolVar(&versionFlag, "version", false, "Print version number")
	flag.Parse()

	if versionFlag {
		fmt.Println(version)
		os.Exit(0)
	}

	// Get credentials
	credentialsPath := os.Getenv("AWS_SHARED_CREDENTIALS_FILE")
	if len(credentialsPath) == 0 {
		usr, err := user.Current()
		if err != nil {
			fmt.Println("Error: Could not locate your home directory. Please set the AWS_SHARED_CREDENTIALS_FILE environment variable.")
			os.Exit(1)
		}
		credentialsPath = fmt.Sprintf("%s/.aws/credentials", usr.HomeDir)
	}

	credentialsProvider := credentials.NewSharedCredentials(credentialsPath, profileFlag)
	creds, err := credentialsProvider.Get()
	check(err)
	fmt.Printf("Using access key %s from profile \"%s\".\n", creds.AccessKeyID, profileFlag)

	if authProfileFlag != "" {
		profileFlag = authProfileFlag
	}

	// Read credentials file
	bytes, err := os.ReadFile(credentialsPath)
	check(err)
	credentialsText := string(bytes)
	// Check if we can find the credentials in the file
	// It's better to detect a malformed file now than after we have created the new key
	re_aws_access_key_id := regexp.MustCompile(fmt.Sprintf(`(?m)^aws_access_key_id *= *%s`, regexp.QuoteMeta(creds.AccessKeyID)))
	re_aws_secret_access_key := regexp.MustCompile(fmt.Sprintf(`(?m)^aws_secret_access_key *= *%s`, regexp.QuoteMeta(creds.SecretAccessKey)))
	if !re_aws_access_key_id.MatchString(credentialsText) || !re_aws_secret_access_key.MatchString(credentialsText) {
		fmt.Println()
		fmt.Printf("Unable to find your credentials in %s\n", credentialsPath)
		fmt.Println("Please make sure your file is formatted like the following:")
		fmt.Println()
		fmt.Printf("aws_access_key_id=%s\n", creds.AccessKeyID)
		fmt.Println("aws_secret_access_key=...")
		fmt.Println()
		os.Exit(1)
	}

	// Create session
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           profileFlag,
	}))

	// sts get-caller-identity
	stsClient := sts.New(sess)
	respGetCallerIdentity, err := stsClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		fmt.Println("Error getting caller identity. Is the key disabled?")
		fmt.Println()
		check(err)
	}
	fmt.Printf("Your user ARN is: %s\n", *respGetCallerIdentity.Arn)

	// mfa
	if mfaFlag || mfaSerialNumber != "" {
		if mfaSerialNumber == "" {
			// iam list-mfa-devices
			// If the UserName field is not specified, the UserName is determined implicitly based on the AWS access key ID used to sign the request.
			iamClient := iam.New(sess)
			respMFADevices, err := iamClient.ListMFADevices(&iam.ListMFADevicesInput{})
			check(err)
			if len(respMFADevices.MFADevices) == 0 {
				fmt.Println("You do not have any MFA devices assigned to your user.")
				os.Exit(1)
			}

			supportedSerialNumbers := make([]*string, 0, len(respMFADevices.MFADevices))
			for _, device := range respMFADevices.MFADevices {
				if !isU2F(*device.SerialNumber) {
					supportedSerialNumbers = append(supportedSerialNumbers, device.SerialNumber)
				}
			}

			if len(supportedSerialNumbers) == 0 {
				fmt.Println()
				fmt.Println("You have an U2F MFA device assigned to your user. These are not supported.")
				fmt.Println("Please add another MFA to your user.")
				os.Exit(1)
			} else if len(supportedSerialNumbers) == 1 {
				mfaSerialNumber = *supportedSerialNumbers[0]
				fmt.Printf("Your MFA serial number is: %s\n\n", mfaSerialNumber)
			} else {
				fmt.Println()
				fmt.Println("You have multiple MFA devices assigned to your user.")
				if len(supportedSerialNumbers) != len(respMFADevices.MFADevices) {
					fmt.Println("Note: You have U2F MFA devices assigned to your user. These are not supported and are not in this list.")
				}
				fmt.Println()
				for i, serialNumber := range supportedSerialNumbers {
					fmt.Printf("%d: %s\n", i+1, *serialNumber)
				}
				fmt.Println()
				if yesFlag {
					mfaSerialNumber = *supportedSerialNumbers[0]
					fmt.Println("Because you used -y, the first MFA device was automatically chosen. You can use -mfa-serial-number to pick a different device.")
				} else {
					var input string
					fmt.Println("Which MFA device do you want to use?")
					fmt.Print("Enter a number from the list above or the full serial number: ")
					_, err = fmt.Scanln(&input)
					check(err)
					if isNumeric(input) {
						i, err := strconv.Atoi(input)
						check(err)
						if i < 1 || i > len(supportedSerialNumbers) {
							fmt.Println("Invalid selection!")
							os.Exit(1)
						}
						mfaSerialNumber = *supportedSerialNumbers[i-1]
					} else {
						mfaSerialNumber = input
					}
				}
			}
		}

		// I have no idea how much work it would be to support U2F
		if isU2F(mfaSerialNumber) {
			fmt.Println("Sorry, U2F MFA devices are not supported. Please use another MFA.")
			os.Exit(1)
		}

		// Prompt for the code
		var code string
		fmt.Print("MFA token code: ")
		_, err = fmt.Scanln(&code)
		check(err)

		// Get the new credentials
		respSessionToken, err := stsClient.GetSessionToken(&sts.GetSessionTokenInput{
			DurationSeconds: aws.Int64(900), // valid for 15 minutes (the minimum)
			SerialNumber:    aws.String(mfaSerialNumber),
			TokenCode:       aws.String(code),
		})
		check(err)

		// Create a new session that use the new credentials
		c := respSessionToken.Credentials
		mfaCreds := credentials.NewStaticCredentials(*c.AccessKeyId, *c.SecretAccessKey, *c.SessionToken)
		sess = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
			Profile:           profileFlag,
			Config:            aws.Config{Credentials: mfaCreds},
		}))
	}
	fmt.Println()

	// iam list-access-keys
	// If the UserName field is not specified, the UserName is determined implicitly based on the AWS access key ID used to sign the request.
	iamClient := iam.New(sess)
	respListAccessKeys, err := iamClient.ListAccessKeys(&iam.ListAccessKeysInput{})
	check(err)

	// Print key information
	fmt.Printf("Your user has %d access key%s:\n", len(respListAccessKeys.AccessKeyMetadata), pluralize(len(respListAccessKeys.AccessKeyMetadata)))
	for _, key := range respListAccessKeys.AccessKeyMetadata {
		respAccessKeyLastUsed, err2 := iamClient.GetAccessKeyLastUsed(&iam.GetAccessKeyLastUsedInput{
			AccessKeyId: key.AccessKeyId,
		})
		if err2 != nil {
			fmt.Printf("- %s (%s, created %s)\n", *key.AccessKeyId, *key.Status, key.CreateDate)
		} else if respAccessKeyLastUsed.AccessKeyLastUsed.LastUsedDate == nil {
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

		if !yesFlag {
			fmt.Println("You have two access keys, which is the maximum number of access keys allowed.")
			fmt.Printf("Do you want to delete %s and create a new key? [yN] ", *respListAccessKeys.AccessKeyMetadata[keyIndex].AccessKeyId)
			if *respListAccessKeys.AccessKeyMetadata[keyIndex].Status == "Active" {
				fmt.Printf("\nWARNING: This key is currently Active! ")
			}
			reader := bufio.NewReader(os.Stdin)
			yn, err2 := reader.ReadString('\n')
			check(err2)
			if yn[0] != 'y' && yn[0] != 'Y' {
				os.Exit(1)
			}
		}

		_, err2 := iamClient.DeleteAccessKey(&iam.DeleteAccessKeyInput{
			AccessKeyId: respListAccessKeys.AccessKeyMetadata[keyIndex].AccessKeyId,
		})
		check(err2)
		fmt.Printf("Deleted access key: %s\n", *respListAccessKeys.AccessKeyMetadata[keyIndex].AccessKeyId)
	} else if !yesFlag {
		cleanupAction := "deactivate"
		if deleteFlag {
			cleanupAction = "delete"
		}
		fmt.Printf("Do you want to create a new key and %s %s? [yN] ", cleanupAction, *respListAccessKeys.AccessKeyMetadata[0].AccessKeyId)
		reader := bufio.NewReader(os.Stdin)
		yn, err2 := reader.ReadString('\n')
		check(err2)
		if yn[0] != 'y' && yn[0] != 'Y' {
			os.Exit(1)
		}
	}

	// Create the new access key
	// If you do not specify a user name, IAM determines the user name implicitly based on the AWS access key ID signing the request.
	respCreateAccessKey, err := iamClient.CreateAccessKey(&iam.CreateAccessKeyInput{})
	check(err)
	newAccessKeyId := *respCreateAccessKey.AccessKey.AccessKeyId
	newSecretAccessKey := *respCreateAccessKey.AccessKey.SecretAccessKey
	fmt.Printf("Created access key: %s\n", newAccessKeyId)

	// Replace key pair in credentials file
	// This search & replace does not limit itself to the specified profile, which is useful if the user is using the same key in multiple profiles
	credentialsText = re_aws_access_key_id.ReplaceAllString(credentialsText, `aws_access_key_id=`+newAccessKeyId)
	credentialsText = re_aws_secret_access_key.ReplaceAllString(credentialsText, `aws_secret_access_key=`+newSecretAccessKey)

	// Verify that the regexp actually replaced something
	if !strings.Contains(credentialsText, newAccessKeyId) || !strings.Contains(credentialsText, newSecretAccessKey) {
		fmt.Println("Error: Failed to replace the old access key.")
		fmt.Printf("Please verify that the file %s is formatted correctly.\n", credentialsPath)
		// Delete the key we created
		_, err2 := iamClient.DeleteAccessKey(&iam.DeleteAccessKeyInput{
			AccessKeyId: aws.String(newAccessKeyId),
		})
		check(err2)
		fmt.Printf("Deleted access key: %s\n", newAccessKeyId)
		os.Exit(1)
	}

	// Write new file
	err = os.WriteFile(credentialsPath, []byte(credentialsText), 0600)
	check(err)
	fmt.Printf("Wrote new key pair to %s\n", credentialsPath)

	// Delete the old key if flag is set, otherwise deactivate it
	if deleteFlag {
		_, err := iamClient.DeleteAccessKey(&iam.DeleteAccessKeyInput{
			AccessKeyId: aws.String(creds.AccessKeyID),
		})
		check(err)
		fmt.Printf("Deleted old access key: %s\n", creds.AccessKeyID)
	} else {
		_, err = iamClient.UpdateAccessKey(&iam.UpdateAccessKeyInput{
			AccessKeyId: aws.String(creds.AccessKeyID),
			Status:      aws.String("Inactive"),
		})
		check(err)
		fmt.Printf("Deactivated old access key: %s\n", creds.AccessKeyID)
		fmt.Println("Please make sure this key is not used elsewhere.")
	}
	fmt.Println("Please note that it may take a minute for your new access key to propagate in the AWS control plane.")
}
