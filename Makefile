VERSION = 1.0.4
LDFLAGS = -ldflags '-s -w'
GOARCH = amd64
linux: export GOOS=linux
darwin: export GOOS=darwin
windows: export GOOS=windows

all: linux darwin windows

linux:
	go build $(LDFLAGS)
	mkdir -p release
	rm -f release/aws-rotate-key-${VERSION}-${GOOS}_${GOARCH}.zip
	zip release/aws-rotate-key-${VERSION}-${GOOS}_${GOARCH}.zip aws-rotate-key

darwin:
	go build $(LDFLAGS)
	mkdir -p release
	rm -f release/aws-rotate-key-${VERSION}-${GOOS}_${GOARCH}.zip
	zip release/aws-rotate-key-${VERSION}-${GOOS}_${GOARCH}.zip aws-rotate-key

windows:
	go build $(LDFLAGS)
	mkdir -p release
	rm -f release/aws-rotate-key-${VERSION}-${GOOS}_${GOARCH}.zip
	zip release/aws-rotate-key-${VERSION}-${GOOS}_${GOARCH}.zip aws-rotate-key.exe

.PHONY: clean
clean:
	rm -rf release
	rm -f aws-rotate-key aws-rotate-key.exe
