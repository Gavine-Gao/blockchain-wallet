FROM golang:1.22.11-alpine3.21 as builder

RUN apk add --no-cache make ca-certificates gcc musl-dev linux-headers git jq bash

COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum

WORKDIR /app

RUN go mod download

ARG CONFIG=config.yml

# build blockchain-wallet-utxo with the shared go.mod & go.sum files
COPY . /app/blockchain-wallet-utxo

WORKDIR /app/blockchain-wallet-utxo

RUN make

FROM alpine:3.18

COPY --from=builder /app/blockchain-wallet-utxo/blockchain-wallet-utxo /usr/local/bin
COPY --from=builder /app/blockchain-wallet-utxo/${CONFIG} /etc/blockchain-wallet-utxo/

WORKDIR /app

ENTRYPOINT ["blockchain-wallet-utxo"]
CMD ["-c", "/etc/blockchain-wallet-utxo/config.yml"]
