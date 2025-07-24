GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

build:
	CGO_ENABLED=0 go build -o ./bin/arista-snake-$(GOOS)-$(GOARCH)
