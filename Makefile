APP := $(shell basename $(shell git remote get-url origin))
REGISTRY := ghcr.io/otrokov
VERSION=$(shell git describe --tags --abbrev=0)-$(shell git rev-parse --short HEAD)
TARGETOS=linux
TARGETARCH =amd64
IMAGE := ghcr.io/otrokov/currencybot:${VERSION}-${TARGETOS}-${TARGETARCH}
NAME := currencybot

format:
	gofmt -s -w ./


lint:
	golint

test:
	go test -v
get:
	go get


build: format
	CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH}	go build -v -o currencybot -ldflags "-X="github.com/otrokov/currencybot/cmd.appVersion=${VERSION}

image:
	docker build . -t ${REGISTRY}/${NAME}:${VERSION}-${TARGETOS}-${TARGETARCH}
push:
	docker push ${REGISTRY}/${NAME}:${VERSION}-${TARGETOS}-${TARGETARCH}



clean:
	rm -rf currencybot
	docker rmi ${REGISTRY}/${NAME}:${VERSION}-${TARGETOS}-${TARGETARCH}