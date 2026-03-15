# Blockchain Wallet

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![gRPC](https://img.shields.io/badge/gRPC-Protocol_Buffers-244C5A?style=for-the-badge&logo=google&logoColor=white)](https://grpc.io/)
[![License](https://img.shields.io/badge/License-MIT-blue?style=for-the-badge)](LICENSE)

A production-grade, multi-chain blockchain wallet infrastructure supporting **18+ blockchains** across both Account-based (ETH, SOL, TRON, etc.) and UTXO-based (BTC, LTC, DASH, etc.) models. Built with Go, designed for high throughput and horizontal scalability.

---

## Architecture Overview

The system is composed of **7 independent microservices** organized into four layers:

| Layer | Service | Responsibility |
|:------|:--------|:---------------|
| **Protocol** | `blockchain-proto` | Protobuf definitions for Account & UTXO gRPC service interfaces |
| **Wallet** | `blockchain-wallet-account` | Multi-chain node abstraction for 18 Account-model chains |
| | `blockchain-wallet-utxo` | Node abstraction for UTXO-model chains (BTC, LTC, etc.) |
| **Sync** | `blockchain-sync-sol` | Real-time Solana on-chain data synchronization |
| | `blockchain-sync-account` | EVM-compatible chain synchronization |
| | `blockchain-sync-btc` | Bitcoin UTXO synchronization |
| **Data** | `chain-explorer-api` | Unified wrapper for Etherscan / SolScan / OKLink APIs |

### Data Flow

```
                      ┌─────────────────────────────────┐
                      │       Application Layer          │
                      │       (Wallet API / DApp)        │
                      └──────────────┬──────────────────┘
                                     │ gRPC
            ┌────────────────────────┼──────────────────────────┐
            ▼                        ▼                          ▼
  ┌──────────────────┐   ┌────────────────────┐   ┌──────────────────┐
  │  sync-account    │   │  sync-sol          │   │  sync-btc        │
  │  EVM Chain Sync  │   │  Solana Sync       │   │  Bitcoin Sync    │
  │  ChannelBank     │   │  ChannelBank       │   │  ChannelBank     │
  │  AddressCache    │   │  AddressCache      │   │  AddressCache    │
  └────────┬─────────┘   └────────┬───────────┘   └────────┬─────────┘
           │ gRPC                 │ gRPC                    │ gRPC
           ▼                      ▼                         ▼
  ┌────────────────────────────────────┐   ┌─────────────────────────┐
  │    blockchain-wallet-account       │   │ blockchain-wallet-utxo  │
  │    ChainDispatcher (18 chains)     │   │ ChainDispatcher (UTXO)  │
  │    ┌─────┬─────┬──────┬────┐      │   │ ┌─────┬─────┬─────┐    │
  │    │ ETH │ SOL │ TRON │ .. │      │   │ │ BTC │ LTC │DASH │    │
  │    └──┬──┴──┬──┴──┬───┴──┬─┘      │   │ └──┬──┴──┬──┴──┬──┘    │
  └───────┼─────┼─────┼──────┼────────┘   └────┼─────┼─────┼───────┘
          │     │     │      │                  │     │     │
          ▼     ▼     ▼      ▼                  ▼     ▼     ▼
    ┌────────────────────────────────────────────────────────────────┐
    │                  RPC Providers / Self-hosted Nodes             │
    │          Alchemy · Helius · QuickNode · GetBlock · Ankr       │
    └────────────────────────────────────────────────────────────────┘
```

## Core Design Patterns

### 1. ChainDispatcher — Strategy-based Multi-chain Routing

All chains implement a unified `IChainAdaptor` interface. The dispatcher routes incoming gRPC requests to the appropriate chain adaptor at runtime. Adding a new chain requires only implementing the interface and registering it — zero changes to the routing layer.

```
gRPC Request → ChainDispatcher → registry[chainName] → Chain Adaptor
                                       ├── Ethereum Adaptor
                                       ├── Solana Adaptor
                                       ├── Tron Adaptor
                                       └── ... (18 chains)
```

### 2. Worker Pipeline — Transaction Classification Engine

Transactions are automatically classified into five categories and dispatched to specialized workers:

```
BaseSynchronizer ──→ ChannelBank ──→ Deposit Worker   (incoming transfers)
   (block scanner)   (ordered)   ──→ Withdraw Worker  (outgoing transfers)
                                 ──→ Internal Worker  (hot/cold wallet mgmt)
                                 ──→ FallBack Worker  (failure recovery)
```

Transaction types: **Deposit** · **Withdrawal** · **Aggregation** · **Hot-to-Cold** · **Cold-to-Hot**

### 3. ChannelBank — Concurrent Scanning with Ordered Consumption

Inspired by [Mantle's derivation pipeline](https://github.com/mantlenetworkio/mantle), ChannelBank uses a **min-heap** to enable concurrent block fetching while guaranteeing strictly ordered output:

```
Multiple goroutines fetch blocks concurrently (unordered)
                    │
                    ▼
            ┌───────────────┐
            │  ChannelBank  │
            │  MinHeap Sort │  ← Sort on arrival
            │ nextExpected  │  ← Strict sequential guarantee
            └───────┬───────┘
                    │ Ordered output
                    ▼
          Consumed in block-number order
```

### 4. AddressCache — Two-layer Probabilistic Filtering

A high-performance address lookup mechanism combining a Bloom filter (~12MB for 10M addresses) with an exact-match hash map. This eliminates database queries for the vast majority of transactions:

```
Tx Address → Bloom Filter (O(1), ~12MB / 10M addresses)
                 │
                 ├── Definitely absent → Skip (majority filtered here)
                 └── Possibly present → Exact HashMap lookup
                                           │
                                           ├── Not found → False positive, skip
                                           └── Found → Return BusinessId + AddressType
```

## Supported Chains

### Account Model (blockchain-wallet-account)

| Chain | Status | Chain | Status |
|:------|:------:|:------|:------:|
| Ethereum | ✅ | Arbitrum | ✅ |
| Solana | ✅ | Optimism | ✅ |
| Tron | ✅ | Linea | ✅ |
| BSC | ✅ | Scroll | ✅ |
| Polygon | ✅ | zkSync Era | ✅ |
| Mantle | ✅ | BTT | ✅ |
| Aptos | ✅ | Sui | ✅ |
| Cosmos | ✅ | TON | ✅ |
| Stellar (XLM) | ✅ | ICP | ✅ |

### UTXO Model (blockchain-wallet-utxo)

| Chain | Status | Chain | Status |
|:------|:------:|:------|:------:|
| Bitcoin | ✅ | Bitcoin Cash | ✅ |
| Litecoin | ✅ | Horizen (ZEN) | ✅ |
| Dash | ✅ | | |

## Getting Started

### Prerequisites

- Go 1.22+
- MySQL 8.0+
- Protocol Buffers compiler (`protoc`)

### Run Wallet Node Service

```bash
# Account-model chains
cd blockchain-wallet-account
cp config.yml.example config.yml   # Configure RPC endpoints
go run main.go

# UTXO-model chains
cd blockchain-wallet-utxo
cp config.yml.example config.yml
go run main.go
```

### Run Sync Service

```bash
# Solana sync
cd blockchain-sync-sol
go run cmd/multichain-sync/main.go \
  --chain-name=Solana \
  --rpc-url=http://localhost:8189 \
  --starting-height=0

# EVM chain sync
cd blockchain-sync-account
go run cmd/multichain-sync/main.go \
  --chain-name=Ethereum \
  --rpc-url=http://localhost:8189
```

## Key Features

| Feature | Description |
|:--------|:------------|
| **Unified gRPC Interface** | All chains accessed through a single protobuf-defined API |
| **Real-time Block Scanning** | Finalized confirmation support with automatic reorg detection |
| **Transaction Builder** | Unsigned transaction construction, signature verification, broadcasting |
| **Address Management** | Generation, validation, and cross-format conversion |
| **Fee Estimation** | Three-tier fee estimation (slow/normal/fast), Solana priority fee support |
| **Explorer API Integration** | Etherscan, SolScan, OKLink unified under a single interface |
| **Prometheus Metrics** | Sync latency, RPC call stats, cache hit rates |
| **Graceful Shutdown** | Parallel worker shutdown with configurable timeout |

## Project Details

<details>
<summary><b>blockchain-proto</b> — Protocol Definitions</summary>

Defines two gRPC service contracts:
- `WalletAccountService` — Account-model chain interface (transfers, balances, blocks, transactions)
- `WalletUtxoService` — UTXO-model chain interface (UTXO queries, transaction construction)
</details>

<details>
<summary><b>blockchain-wallet-account</b> — Account Chain Node Service</summary>

Routes gRPC requests via `ChainDispatcher` to 18 chain adaptors, each implementing `IChainAdaptor`:
- `GetBlockByNumber` / `GetBlockByHash` — Block queries
- `GetTxByHash` / `GetTxByAddress` — Transaction queries
- `GetAccount` / `GetBalance` — Account & balance lookups
- `BuildUnSignTransaction` / `SendTx` — Transaction construction & broadcast
</details>

<details>
<summary><b>blockchain-sync-sol</b> — Solana Sync Engine</summary>

Core components:
- **BaseSynchronizer** — Concurrent block scanner with ChannelBank-based ordered consumption
- **AddressCache** — Bloom filter + HashMap two-layer address filtering
- **Workers** — Deposit / Withdraw / Internal / FallBack transaction processors
- **Notifier** — Webhook-based transaction notification with fallback retry
</details>

<details>
<summary><b>chain-explorer-api</b> — Block Explorer API Abstraction</summary>

Unified wrapper over multiple blockchain explorer APIs:
- **Etherscan** — ETH/EVM chain transactions, accounts, gas estimation
- **SolScan** — Solana transactions, tokens, account data
- **OKLink** — Multi-chain transactions, UTXO queries
</details>

## License

This project is licensed under the [MIT License](LICENSE).
