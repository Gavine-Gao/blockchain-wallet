# Gavine - 区块链钱包基础设施

多链区块链钱包基础设施项目集，包含以下子项目：

## 项目结构

| 项目 | 说明 |
|------|------|
| `blockchain-proto` | Protocol Buffer 定义（gRPC 接口） |
| `blockchain-wallet-account` | Account 模型链节点服务（ETH/SOL/TRON 等 18 条链） |
| `blockchain-wallet-utxo` | UTXO 模型链节点服务（BTC） |
| `blockchain-sync-sol` | Solana 链上数据同步服务 |
| `blockchain-sync-account` | Account 模型链同步服务 |
| `blockchain-sync-btc` | BTC UTXO 同步服务 |

## 架构

```
┌──────────────────────────┐
│     业务层（钱包 API）      │
└─────────┬────────────────┘
          │ gRPC
┌─────────▼────────────────┐
│   blockchain-sync-*      │  链上数据同步
│   (sol / account / btc)  │
└─────────┬────────────────┘
          │ gRPC
┌─────────▼────────────────┐
│  blockchain-wallet-*     │  链节点抽象层
│  (account / utxo)        │
└─────────┬────────────────┘
          │ JSON-RPC / REST
┌─────────▼────────────────┐
│  第三方节点 / 自建节点      │
│  (Alchemy/Helius/GetBlock)│
└──────────────────────────┘
```
