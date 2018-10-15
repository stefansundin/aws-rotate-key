# aws-rotate-key

As a security best practice, AWS recommends that administrators require
IAM users to periodically [regenerate their API access keys](http://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_access-keys.html#Using_RotateAccessKey).
This `aws-rotate-key` tool allows users to easily rotate all of the AWS access keys defined in their local 
[aws credentials file](http://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html#cli-multiple-profiles).

The program will use the AWS API to check which access keys exist
for the provided profile. If only one access key exists, then it will
deactivate that key and update your credentials file to use a newly 
generated key. The old key will only be deactivated (**not** deleted),
so that if you later find out you use the old key elsewhere, you
can open the AWS console and reactivate it. If two access keys exist,
then you will be asked whether you want to delete the key which is
not currently configured in your credentials file to create an empty
slot for the key rotation. Then, it will perform the same key rotation 
logic on the remaining key.


## Usage
Usage of aws-rotate-key:
```
  -profile string
    	The profile to use. (default "default")
  -version
    	Print version number (1.0.4)
  -y
        Automatic "yes" to prompts.
  -d
        Delete old key without deactivation.

```

## Example

```
someone$ aws-rotate-key --profile primary
Using access key A123 from profile "primary".
Your user arn is: arn:aws:iam::123456789012:user/someone@example.com

You have 2 access keys associated with your user:
-A123 (Inactive, created 2015-01-01 02:55:00 +0000 UTC, last used 2016-01-01 00:02:00 +0000 UTC for service sts in us-east-1)
- B123 (Active, created 2016-01-01 00:02:47 +0000 UTC, last used 2016-01-01 00:03:00 +0000 UTC for service s3 in N/A)

You have two access keys, which is the max number of access keys.
Do you want to delete A123 and create a new key? [yN] y
Deleted access key A123.
Created access key C123.
Wrote new key pair to /Users/someone/.aws/credentials
Deactivated old access key B123.
Please make sure this key is not used elsewhere.
```

## Install

You can download the 64-bit binaries from
[the releases section](https://github.com/Fullscreen/aws-rotate-key/releases/latest)
of this repository.

Or, you can use our homebrew tap on OSX:

```
brew tap fullscreen/tap
brew install aws-rotate-key
aws-rotate-key
```

## Setup

Make sure your users have permissions to update their own access keys via the CLI. The AWS
documentation [here](http://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_delegate-permissions_examples.html#creds-policies-credentials)
explains the required permissions and the following IAM profile should get you setup:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": [
                "iam:ListAccessKeys",
                "iam:GetAccessKeyLastUsed",
                "iam:DeleteAccessKey",
                "iam:CreateAccessKey",
                "iam:UpdateAccessKey"
            ],
            "Effect": "Allow",
            "Resource": [
                "arn:aws:iam::AWS_ACCOUNT_ID:user/${aws:username}"
            ]
        }
    ]
}
```

Replace `AWS_ACCOUNT_ID` with your AWS account id.

## Contribute

To download and hack on the source code, run:
```
$ go get -u github.com/Fullscreen/aws-rotate-key
$ cd $GOPATH/src/github.com/Fullscreen/aws-rotate-key
$ go build
```
