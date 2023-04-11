build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o aws-registration-callback main.go

package: build
	zip aws-registration-callback.zip aws-registration-callback

TAG ?= dev
REGISTRY ?= workstation
image:
	docker build -t $(REGISTRY):$(TAG) .