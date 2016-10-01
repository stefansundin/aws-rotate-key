# aws-rotate-key

Easily rotate all of the AWS access keys defined in your local 
[aws credentials file](http://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html#cli-multiple-profiles).

By running `aws-rotate-key`, the program will find out what AWS key you have,
rotate it for you, and then update your credentials file.

The program will use the AWS API to check which API access keys you
have created for the provided profile. If you only have a single access
key, then it will deactivate the key currently configured in your 
credentials file and update your credentials file to use a newly 
generated key. The old key will only be deactivated (**not** deleted),
so that if you later find out you use the old key elsewhere, you
can open the AWS console and reactivate it. If you already have two
access keys for the provided profile, the program will ask you if
you want to delete the key not currently configured in your credentials
file. Then, it will perform the key rotation logic now that you have
a single access key.


## Usage
Usage of aws-rotate-key:
```
  -profile string
    	The profile to use. (default "default")
  -version
    	Print version number (1.0.0)
  -y	Automatic "yes" to prompts.
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

## IAM policy

Make sure your users have permissions to update their access keys. Here is an
example policy you can create:

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
