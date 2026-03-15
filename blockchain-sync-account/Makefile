GITCOMMIT := $(shell git rev-parse HEAD)
GITDATE := $(shell git show -s --format='%ct')

LDFLAGSSTRING +=-X main.GitCommit=$(GITCOMMIT)
LDFLAGSSTRING +=-X main.GitDate=$(GITDATE)
LDFLAGS := -ldflags "$(LDFLAGSSTRING)"

blockchain-sync:
	env GO111MODULE=on go build -v $(LDFLAGS) ./cmd/blockchain-sync

clean:
	rm blockchain-sync

test:
	go test -v ./...

protogo:
	sh ./sh/go_compile.sh

lint:
	golangci-lint run ./...

.PHONY: \
	blockchain-sync \
	clean \
	test \
	protogo \
	lint