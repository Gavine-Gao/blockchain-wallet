<!--
parent:
  order: false
-->

<div align="center">
  <h1> blockchain-wallet-account repo </h1>
</div>

<div align="center">
  <a href="https://github.com/Gavine-Gao/blockchain-wallet-account/releases/latest">
    <img alt="Version" src="https://img.shields.io/github/tag/blockchain-wallet- account/blockchain-wallet-account.svg" />
  </a>
  <a href="https://github.com/blockchain-wallet- account/blockchain-wallet-account/blob/main/LICENSE">
    <img alt="License: Apache-2.0" src="https://img.shields.io/github/license/blockchain-wallet- account/blockchain-wallet-account.svg" />
  </a>
  <a href="https://pkg.go.dev/github.com/blockchain-wallet- account/blockchain-wallet-account">
    <img alt="GoDoc" src="https://godoc.org/github.com/blockchain-wallet- account/blockchain-wallet-account?status.svg" />
  </a>
</div>

This repo is account model chains rpc service gateway. currently support `Ethereum`, `Aptos`, `Cosmos`, `Sui`, `Solana` etc, written in golang, provides grpc interface for upper-layer service access

**Tips**: need [Go 1.22+](https://golang.org/dl/)

## Install

### Install dependencies
```bash
go mod tidy
```
### build
```bash
go build or go install blockchain-wallet-account
```

### start
```bash
./blockchain-wallet-account -c ./config.yml
```

### Start the RPC interface test interface

```bash
grpcui -plaintext 127.0.0.1:8189
```

## Contribute

### 1.fork repo

fork blockchain-wallet-account to your github

### 2.clone repo

```bash
git@github.com:gavine/blockchain-wallet-account.git
```

### 3. create new branch and commit code

```bash
git branch -C xxx
git checkout xxx

coding

git add .
git commit -m "xxx"
git push origin xxx
```

### 4.commit PR

Have a pr on your github and submit it to the blockchain-wallet-account repository

### 5.review

After the blockchain-wallet-account code maintainer has passed the review, the code will be merged into the blockchain-wallet-account library. At this point, your PR submission is complete

### 6.Disclaimer

This code has not yet been audited, and should not be used in any production systems.
