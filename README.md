# Blockchain Wallet

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue?style=flat-square)](LICENSE)

多链区块链钱包基础设施，支持 **18+ 条公链**，涵盖 Account 模型（ETH/SOL/TRON 等）和 UTXO 模型（BTC/LTC/DASH 等）。

## 📦 项目结构

```
blockchain-wallet/
├── blockchain-proto              # Protocol Buffer 定义（gRPC 接口规范）
├── blockchain-wallet-account     # Account 模型链节点服务
├── blockchain-wallet-utxo        # UTXO 模型链节点服务
├── blockchain-sync-sol           # Solana 链上数据同步服务
├── blockchain-sync-account       # Account 模型链同步服务
├── blockchain-sync-btc           # BTC UTXO 同步服务
└── chain-explorer-api            # 链浏览器 API 封装（Etherscan/SolScan/OKLink）
```

## 🏗️ 系统架构

系统由 **7 个独立微服务** 组成，分为四层：

| 层级 | 项目 | 职责 |
|:-----|:-----|:-----|
| **协议层** | `blockchain-proto` | Protobuf 定义，Account/UTXO 两套 gRPC 服务接口 |
| **钱包层** | `blockchain-wallet-account` | Account 模型多链钱包，支持 18 条链 |
| | `blockchain-wallet-utxo` | UTXO 模型钱包，支持 BTC/LTC/DASH 等 |
| **同步层** | `blockchain-sync-sol` | Solana 链交易同步 |
| | `blockchain-sync-account` | EVM 兼容链交易同步 |
| | `blockchain-sync-btc` | BTC 链交易同步 |
| **数据层** | `chain-explorer-api` | Etherscan/SolScan/OKLink 浏览器 API 统一封装 |

### 整体数据流

```
                    ┌─────────────────────────────┐
                    │       业务层（钱包 API）       │
                    └──────────────┬──────────────┘
                                   │ gRPC
          ┌────────────────────────┼────────────────────────┐
          ▼                        ▼                        ▼
┌─────────────────┐   ┌──────────────────┐   ┌──────────────────┐
│ sync-account    │   │ sync-sol         │   │ sync-btc         │
│ EVM 链同步       │   │ Solana 链同步     │   │ BTC 链同步        │
│ ChannelBank     │   │ ChannelBank      │   │ ChannelBank      │
│ AddressCache    │   │ AddressCache     │   │ AddressCache     │
└────────┬────────┘   └────────┬─────────┘   └────────┬─────────┘
         │ gRPC                │ gRPC                  │ gRPC
         ▼                     ▼                       ▼
┌──────────────────────────────────┐   ┌────────────────────────┐
│    blockchain-wallet-account     │   │ blockchain-wallet-utxo │
│    ChainDispatcher (18 条链)      │   │ ChainDispatcher (BTC)  │
│    ┌─────┬─────┬─────┬────┐     │   │ ┌─────┬─────┬────┐    │
│    │ ETH │ SOL │TRON │ .. │     │   │ │ BTC │ LTC │DASH│    │
│    └──┬──┴──┬──┴──┬──┴──┬─┘     │   │ └──┬──┴──┬──┴──┬─┘    │
└───────┼─────┼─────┼─────┼───────┘   └────┼─────┼─────┼──────┘
        │     │     │     │                 │     │     │
        ▼     ▼     ▼     ▼                 ▼     ▼     ▼
   ┌──────────────────────────────────────────────────────────┐
   │            第三方节点 / 自建节点                            │
   │     Alchemy · Helius · QuickNode · GetBlock · Ankr       │
   └──────────────────────────────────────────────────────────┘
```

## 🧩 核心设计模式

### 1. ChainDispatcher — 多链适配器

所有链通过统一的 `IChainAdaptor` 接口接入，新增链只需实现接口 + 注册工厂：

```
gRPC 请求 → ChainDispatcher → registry[chainName] → 对应链的 Adaptor
                                    ├── Ethereum Adaptor
                                    ├── Solana Adaptor
                                    ├── Tron Adaptor
                                    └── ... 18 条链
```

### 2. Worker 模式 — 四类交易处理

```
BaseSynchronizer ──→ ChannelBank ──→ Deposit Worker (充值处理)
   (区块扫描)        (流式排序)   ──→ Withdraw Worker (提现确认)
                                 ──→ Internal Worker (归集/冷热)
                                 ──→ FallBack Worker (失败重试)
```

交易类型自动判定：**充值** / **提现** / **归集** / **热转冷** / **冷转热**。

### 3. ChannelBank — 并发扫块 + 有序消费

```
多个 goroutine 并发拉取区块（乱序）
              │
              ▼
      ┌───────────────┐
      │  ChannelBank  │
      │ MinHeap 排序   │  ← 收到即排序
      │ nextExpected  │  ← 严格保证顺序
      └───────┬───────┘
              │ 有序输出
              ▼
      按区块号顺序消费
```

灵感来源于 Mantle 的 derivation pipeline，用最小堆实现流式排序。

### 4. AddressCache — 双层地址过滤

```
交易地址 → Bloom Filter (O(1), ~12MB/千万地址)
              │
              ├── 一定不存在 → 跳过（大部分交易在这里过滤）
              └── 可能存在 → Map 精确匹配
                               │
                               ├── 不存在 → 误报，跳过
                               └── 存在 → 获取 BusinessId + AddressType
```

布隆过滤器 100M bits / 3 hash，千万级地址仅 ~12MB 内存，避免每笔交易查 DB。

## ⛓️ 支持链

### Account 模型（blockchain-wallet-account）

| 链 | 状态 | 链 | 状态 |
|:---|:---:|:---|:---:|
| Ethereum | ✅ | Arbitrum | ✅ |
| Solana | ✅ | Optimism | ✅ |
| Tron | ✅ | Linea | ✅ |
| BSC | ✅ | Scroll | ✅ |
| Polygon | ✅ | zkSync | ✅ |
| Mantle | ✅ | BTT | ✅ |
| Aptos | ✅ | Sui | ✅ |
| Cosmos | ✅ | TON | ✅ |
| Stellar (XLM) | ✅ | ICP | ✅ |

### UTXO 模型（blockchain-wallet-utxo）

| 链 | 状态 | 链 | 状态 |
|:---|:---:|:---|:---:|
| Bitcoin | ✅ | Bitcoin Cash | ✅ |
| Litecoin | ✅ | Horizen (ZEN) | ✅ |
| Dash | ✅ | | |

## 🚀 快速开始

### 前置依赖

- Go 1.22+
- MySQL 8.0+
- Protocol Buffers 编译器（protoc）

### 启动节点服务

```bash
# Account 模型链
cd blockchain-wallet-account
cp config.yml.example config.yml  # 修改 RPC 节点配置
go run main.go

# UTXO 模型链
cd blockchain-wallet-utxo
cp config.yml.example config.yml
go run main.go
```

### 启动同步服务

```bash
# Solana 同步
cd blockchain-sync-sol
go run cmd/multichain-sync/main.go \
  --chain-name=Solana \
  --rpc-url=http://localhost:8189 \
  --starting-height=0

# Account 模型链同步
cd blockchain-sync-account
go run cmd/multichain-sync/main.go \
  --chain-name=Ethereum \
  --rpc-url=http://localhost:8189
```

## 🔧 核心功能

| 功能 | 说明 |
|:-----|:-----|
| **统一 gRPC 接口** | 所有链通过相同的 protobuf 接口调用，业务层无需关心链差异 |
| **实时扫块同步** | 支持 finalized 确认、自动重组检测（reorg） |
| **交易构建与签名** | 支持未签名交易构建、签名验证、广播 |
| **地址管理** | 地址生成、验证、格式转换 |
| **费用估算** | 支持慢/正常/快三档费用估算，Solana 支持优先费 |
| **浏览器 API 集成** | Etherscan、SolScan、OKLink 统一封装 |
| **Prometheus 监控** | 同步延迟、RPC 调用、地址缓存命中率 |
| **优雅关闭** | 并行关闭 worker + 超时控制 |

## 📁 各子项目说明

<details>
<summary><b>blockchain-proto</b> — gRPC 接口定义</summary>

定义了两套 gRPC 服务接口：
- `WalletAccountService` — Account 模型链接口（转账/余额/区块/交易查询等）
- `WalletUtxoService` — UTXO 模型链接口（UTXO 查询/交易构建等）
</details>

<details>
<summary><b>blockchain-wallet-account</b> — Account 链节点服务</summary>

通过 `ChainDispatcher` 将 gRPC 请求路由到 18 条链的适配器。每条链实现 `IChainAdaptor` 接口：
- `GetBlockByNumber` / `GetBlockByHash` — 区块查询
- `GetTxByHash` / `GetTxByAddress` — 交易查询
- `GetAccount` / `GetBalance` — 账户余额
- `BuildUnSignTransaction` / `SendTx` — 交易构建发送
</details>

<details>
<summary><b>blockchain-sync-sol</b> — Solana 同步服务</summary>

核心组件：
- `BaseSynchronizer` — 区块扫描引擎（并发拉取 + ChannelBank 有序消费）
- `AddressCache` — 布隆过滤器 + Map 双层地址过滤
- `Deposit/Withdraw/Internal/FallBack` — 四类交易 Worker
- `Notifier` — 交易通知回调（支持 fallback 重试）
</details>

<details>
<summary><b>chain-explorer-api</b> — 链浏览器 API</summary>

统一封装多个区块链浏览器 API：
- Etherscan — ETH/EVM 链交易、账户、Gas 查询
- SolScan — Solana 交易、代币、账户查询
- OKLink — 多链交易、UTXO 查询
</details>

## 📄 License

[MIT](LICENSE)
