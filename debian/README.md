```
sudo apt install debhelper dh-golang golang golang-github-aws-aws-sdk-go-dev

# in the repository root:
git clean -fdX
tar cvJf ../aws-rotate-key_1.0.7.orig.tar.xz --exclude=debian *
debuild -i -us -uc -b
```
