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

```
┌──────────────────────────────────┐
│         业务层（钱包 API）         │
└───────────────┬──────────────────┘
                │ gRPC
┌───────────────▼──────────────────┐
│       blockchain-sync-*          │
│   链上数据同步（SOL / ETH / BTC） │
│   • 实时扫块 • 交易解析 • 通知回调 │
└───────────────┬──────────────────┘
                │ gRPC
┌───────────────▼──────────────────┐
│      blockchain-wallet-*         │
│   链节点抽象层（account / utxo）   │
│   • 统一接口 • 多链适配 • 签名构建 │
└───────────────┬──────────────────┘
                │ JSON-RPC / REST
┌───────────────▼──────────────────┐
│     第三方节点 / 自建节点          │
│   Alchemy · Helius · QuickNode   │
└──────────────────────────────────┘
```

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

| 链 | 状态 |
|:---|:---:|
| Bitcoin | ✅ |
| Litecoin | ✅ |
| Dash | ✅ |
| Bitcoin Cash | ✅ |
| Horizen (ZEN) | ✅ |

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

## 📄 License

[MIT](LICENSE)
