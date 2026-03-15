FROM golang:1.22.11-alpine3.21 as builder

RUN apk add --no-cache make ca-certificates gcc musl-dev linux-headers git jq bash

COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum

WORKDIR /app

RUN go mod download

ARG CONFIG=config.yml

# build blockchain-wallet-account with the shared go.mod & go.sum files
COPY . /app/blockchain-wallet-account

WORKDIR /app/blockchain-wallet-account

RUN make

FROM alpine:3.18

COPY --from=builder /app/blockchain-wallet-account/blockchain-wallet-account /usr/local/bin
COPY --from=builder /app/blockchain-wallet-account/${CONFIG} /etc/blockchain-wallet-account/

WORKDIR /app

ENTRYPOINT ["blockchain-wallet-account"]
CMD ["-c", "/etc/blockchain-wallet-account/config.yml"]
