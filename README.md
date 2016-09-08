# aws-rotate-key

Easily rotate your AWS key.

By running `aws-rotate-key`, the program will find out what AWS key you have,
rotate it for you and then update your credentials file.

The program will automatically deactivate your old key, but it will not delete
it. If you later find out that you use the key elsewhere, you can open the AWS
console and reactivate it.

If you already have two keys, the program will ask you if you want to delete the
one you are not currently using.

## Binaries

There are 64-bit binaries in [the releases section](https://github.com/Fullscreen/aws-rotate-key/releases/latest).

You can use our homebrew tap on Mac:

```
brew tap fullscreen/tap
brew install aws-rotate-key
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
