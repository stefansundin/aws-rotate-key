# aws-rotate-key

As a security best practice, AWS recommends that users periodically
[regenerate their API access keys](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_access-keys.html#Using_RotateAccessKey).
This tool simplifies the rotation of access keys defined in your
[credentials file](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html#cli-multiple-profiles).

When run, the program will list the current access keys associated with your
IAM user, and print the steps it has to perform to rotate them.
It will then wait for your confirmation before continuing.

## Usage

```
$ aws-rotate-key --help
Usage of aws-rotate-key:
  -d	Delete old key without deactivation.
  -mfa
    	Use MFA.
  -profile string
    	The profile to use. (default "default")
  -version
    	Print version number
  -y	Automatic "yes" to prompts.
```

## Example

```
$ aws-rotate-key --profile work
Using access key AKIAJMIGD6UPCXCFWVOA from profile "work".
Your user ARN is: arn:aws:iam::123456789012:user/your_username

You have 2 access keys associated with your user:
- AKIAI3KI7UC6BPI4O57A (Inactive, created 2018-11-22 21:47:46 +0000 UTC, last used 2018-11-30 20:35:41 +0000 UTC for service s3 in us-west-2)
- AKIAJMIGD6UPCXCFWVOA (Active, created 2018-11-30 21:55:57 +0000 UTC, last used 2018-12-20 12:14:10 +0000 UTC for service s3 in us-west-2)

You have two access keys, which is the max number of access keys.
Do you want to delete AKIAI3KI7UC6BPI4O57A and create a new key? [yN] y
Deleted access key AKIAI3KI7UC6BPI4O57A.
Created access key AKIAIX46CKYT7E5I3KVQ.
Wrote new key pair to /Users/your_username/.aws/credentials
Deactivated old access key AKIAJMIGD6UPCXCFWVOA.
Please make sure this key is not used elsewhere.
Please note that it may take a minute for your new access key to propagate in the AWS control plane.
```

## Install

You can download binaries from [the releases section](https://github.com/stefansundin/aws-rotate-key/releases/latest).

You can use Homebrew to install on macOS:

```
$ brew install stefansundin/tap/aws-rotate-key
```

## Setup

Make sure your users have permissions to update their own access keys.
The following AWS documentation page explains the required permissions:
https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_delegate-permissions_examples.html#creds-policies-credentials.

The following IAM policy is enough for aws-rotate-key:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "iam:ListAccessKeys",
                "iam:GetAccessKeyLastUsed",
                "iam:DeleteAccessKey",
                "iam:CreateAccessKey",
                "iam:UpdateAccessKey"
            ],
            "Resource": [
                "arn:aws:iam::AWS_ACCOUNT_ID:user/${aws:username}"
            ]
        }
    ]
}
```

Replace `AWS_ACCOUNT_ID` with your AWS account id.

### Require MFA

You can require MFA by adding a `Condition` clause. Please note that you
have to use the `-mfa` option when running the program.

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "iam:ListMFADevices"
            ],
            "Resource": [
                "arn:aws:iam::AWS_ACCOUNT_ID:user/${aws:username}"
            ]
        },
        {
            "Effect": "Allow",
            "Action": [
                "iam:ListAccessKeys",
                "iam:GetAccessKeyLastUsed",
                "iam:DeleteAccessKey",
                "iam:CreateAccessKey",
                "iam:UpdateAccessKey"
            ],
            "Resource": [
                "arn:aws:iam::AWS_ACCOUNT_ID:user/${aws:username}"
            ],
            "Condition": {
                "Bool": {
                    "aws:MultiFactorAuthPresent": true
                }
            }
        }
    ]
}
```

Note that this makes it harder to rotate access keys using aws-cli commands,
as it only supports MFA when assuming roles. You will still be able to use
the AWS management console.

## Contribute

To download and hack on the source code, run:

```
$ git clone https://github.com/stefansundin/aws-rotate-key.git
$ cd aws-rotate-key
$ go build
```
