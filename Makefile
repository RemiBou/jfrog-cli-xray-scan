GOOS = $(shell go env GOOS)
GOARCH = amd64

build: clean
	go build -o "xray-scan-${GOOS}-${GOARCH}"

build-install:: build
	mv xray-scan-${GOOS}-${GOARCH} ~/.jfrog/plugins/xray-scan

clean:
	rm -f xray-scan-${GOOS}-${GOARCH}

all: build-linux-amd64 build-darwin-amd64

test: build
	go test ./... -coverprofile=c.out

itest: build
	go test ./... --tags=itest -coverprofile=c.out

build-linux-amd64:
	@$(call echoDebug,"")
	@GOOS="linux" GOARCH="amd64" $(MAKE) build

build-darwin-amd64:
	@$(call echoDebug,"")
	@GOOS="darwin" GOARCH="amd64" $(MAKE) build