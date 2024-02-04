```shell
sudo apt install devscripts debhelper dh-golang golang golang-any golang-github-aws-aws-sdk-go-v2-dev

# in the repository root:
git clean -fdX
tar cvJf ../aws-rotate-key_1.2.0.orig.tar.xz --exclude=debian *

# build deb file:
debuild -i -us -uc -b

# publish to ppa:
debuild -S
cd ..
dput ppa:stefansundin/aws-rotate-key aws-rotate-key_1.2.0-1_source.changes
```
