<!--
parent:
  order: false
-->

<div align="center">
  <h1> blockchain-wallet-utxo repo </h1>
</div>

<div align="center">
  <a href="https://github.com/Gavine-Gao/blockchain-wallet-utxo/releases/latest">
    <img alt="Version" src="https://img.shields.io/github/tag/Gavine-Gao/blockchain-wallet-utxo.svg" />
  </a>
  <a href="https://github.com/Gavine-Gao/blockchain-wallet-utxo/blob/main/LICENSE">
    <img alt="License: Apache-2.0" src="https://img.shields.io/github/license/Gavine-Gao/blockchain-wallet-utxo.svg" />
  </a>
  <a href="https://pkg.go.dev/github.com/Gavine-Gao/blockchain-wallet-utxo">
    <img alt="GoDoc" src="https://godoc.org/github.com/Gavine-Gao/blockchain-wallet-utxo?status.svg" />
  </a>
</div>

This repo is utxo chains rpc service gateway. currently support `Bitcoin`, `Bitcoincash`, `Dash`, `Dogecoin`, `Litecoin`, written in golang, provides grpc interface for upper-layer service access

**Tips**: need [Go 1.22+](https://golang.org/dl/)

## Install

### Install dependencies
```bash
go mod tidy
```
### build
```bash
go build or go install blockchain-wallet-utxo
```

### start
```bash
./blockchain-wallet-utxo -c ./config.yml
```

### Start the RPC interface test interface

```bash
grpcui -plaintext 127.0.0.1:8389
```

## Contribute

### 1.fork repo

fork blockchain-wallet-utxo to your github

### 2.clone repo

```bash
git@github.com:guoshijiang/blockchain-wallet-utxo.git
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

Have a pr on your github and submit it to the blockchain-wallet-utxo repository

### 5.review

After the blockchain-wallet-utxo code maintainer has passed the review, the code will be merged into the blockchain-wallet-utxo library. At this point, your PR submission is complete

### 6.Disclaimer

This code has not yet been audited, and should not be used in any production systems.
