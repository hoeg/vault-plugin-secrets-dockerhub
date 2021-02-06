GOARCH = amd64

UNAME = $(shell uname -s)

ifndef OS
	ifeq ($(UNAME), Linux)
		OS = linux
	else ifeq ($(UNAME), Darwin)
		OS = darwin
	endif
endif

.DEFAULT_GOAL := all

PLUGIN_NAME=vault-plugin-secrets-dockerhub
BIN=./vault/plugins/$(PLUGIN_NAME)

all: fmt build start

build:
	GOOS=$(OS) GOARCH="$(GOARCH)" go build -o $(BIN) cmd/vault-plugin-secrets-dockerhub/main.go

start:
	vault server -dev -dev-root-token-id=root -dev-plugin-dir=./vault/plugins

register: build
	vault write sys/plugins/catalog/$(PLUGIN_NAME) sha_256=$(shell shasum -a 256 $(BIN) | cut -d " " -f1) command="$(PLUGIN_NAME)"

enable: register
	vault secrets enable -path=dockerhub ${PLUGIN_NAME}

clean:
	rm -f $(BIN)

fmt:
	go fmt $$(go list ./...)

.PHONY: build clean fmt start enable register