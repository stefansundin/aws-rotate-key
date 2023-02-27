VERSION = 1.1.0
LDFLAGS = -ldflags '-s -w' -gcflags=-trimpath=${PWD} -asmflags=-trimpath=${PWD}

linux_amd64: export GOOS=linux
linux_amd64: export GOARCH=amd64
linux_arm: export GOOS=linux
linux_arm: export GOARCH=arm
linux_arm: export GOARM=6
linux_arm64: export GOOS=linux
linux_arm64: export GOARCH=arm64
darwin_amd64: export GOOS=darwin
darwin_amd64: export GOARCH=amd64
darwin_arm64: export GOOS=darwin
darwin_arm64: export GOARCH=arm64
windows_amd64: export GOOS=windows
windows_amd64: export GOARCH=amd64
windows_arm: export GOOS=windows
windows_arm: export GOARCH=arm
windows_arm64: export GOOS=windows
windows_arm64: export GOARCH=arm64

.PHONY: all linux_amd64 linux_arm linux_arm64 darwin_amd64 darwin_arm64 windows_amd64 windows_arm windows_arm64 clean

all: linux_amd64 linux_arm linux_arm64 darwin_amd64 darwin_arm64 windows_amd64 windows_arm windows_arm64

linux_amd64:
	go build $(LDFLAGS)
	mkdir -p release
	rm -f release/aws-rotate-key-${VERSION}-${GOOS}_${GOARCH}.zip
	zip release/aws-rotate-key-${VERSION}-${GOOS}_${GOARCH}.zip aws-rotate-key

linux_arm:
	go build $(LDFLAGS)
	mkdir -p release
	rm -f release/aws-rotate-key-${VERSION}-${GOOS}_${GOARCH}.zip
	zip release/aws-rotate-key-${VERSION}-${GOOS}_${GOARCH}.zip aws-rotate-key

linux_arm64:
	go build $(LDFLAGS)
	mkdir -p release
	rm -f release/aws-rotate-key-${VERSION}-${GOOS}_${GOARCH}.zip
	zip release/aws-rotate-key-${VERSION}-${GOOS}_${GOARCH}.zip aws-rotate-key

darwin_amd64:
	go build $(LDFLAGS)
	mkdir -p release
	rm -f release/aws-rotate-key-${VERSION}-${GOOS}_${GOARCH}.zip
	zip release/aws-rotate-key-${VERSION}-${GOOS}_${GOARCH}.zip aws-rotate-key

darwin_arm64:
	go build $(LDFLAGS)
	mkdir -p release
	rm -f release/aws-rotate-key-${VERSION}-${GOOS}_${GOARCH}.zip
	zip release/aws-rotate-key-${VERSION}-${GOOS}_${GOARCH}.zip aws-rotate-key

windows_amd64:
	go build $(LDFLAGS)
	mkdir -p release
	rm -f release/aws-rotate-key-${VERSION}-${GOOS}_${GOARCH}.zip
	zip release/aws-rotate-key-${VERSION}-${GOOS}_${GOARCH}.zip aws-rotate-key.exe

windows_arm:
	go build $(LDFLAGS)
	mkdir -p release
	rm -f release/aws-rotate-key-${VERSION}-${GOOS}_${GOARCH}.zip
	zip release/aws-rotate-key-${VERSION}-${GOOS}_${GOARCH}.zip aws-rotate-key.exe

windows_arm64:
	go build $(LDFLAGS)
	mkdir -p release
	rm -f release/aws-rotate-key-${VERSION}-${GOOS}_${GOARCH}.zip
	zip release/aws-rotate-key-${VERSION}-${GOOS}_${GOARCH}.zip aws-rotate-key.exe

clean:
	rm -rf release
	rm -f aws-rotate-key aws-rotate-key.exe
