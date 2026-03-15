SHELL := /bin/bash

GITCOMMIT := $(shell git rev-parse HEAD)
GITDATE := $(shell git show -s --format='%ct')

LDFLAGSSTRING +=-X main.GitCommit=$(GITCOMMIT)
LDFLAGSSTRING +=-X main.GitDate=$(GITDATE)
LDFLAGSSTRING +=-X main.GitVersion=$(GITVERSION)
LDFLAGS :=-ldflags "$(LDFLAGSSTRING)"

blockchain-wallet-utxo:
	env GO111MODULE=on go build $(LDFLAGS)
.PHONY: blockchain-wallet-utxo

clean:
	rm blockchain-wallet-utxo

test:
	go test -v ./...

lint:
	golangci-lint run ./...