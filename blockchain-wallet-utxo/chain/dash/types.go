package dash

import (
	"github.com/btcsuite/btcd/chaincfg"
)

// Dash 主网参数
var dashMainNetParams = chaincfg.Params{
	Net:              0xbf0c6bbd,
	PubKeyHashAddrID: 0x4c, // X 开头
	ScriptHashAddrID: 0x10, // 7 开头
	Bech32HRPSegwit:  "dash",
}

// BlockData Dash 区块数据结构
type BlockData struct {
	Hash              string   `json:"hash"`
	Confirmations     uint64   `json:"confirmations"`
	Size              uint64   `json:"size"`
	Height            uint64   `json:"height"`
	Version           uint64   `json:"version"`
	Merkleroot        string   `json:"merkleroot"`
	Tx                []string `json:"tx"`
	Time              uint64   `json:"time"`
	Nonce             uint64   `json:"nonce"`
	Bits              string   `json:"bits"`
	Difficulty        float64  `json:"difficulty"`
	PreviousBlockHash string   `json:"previousblockhash"`
	NextBlockHash     string   `json:"nextblockhash"`
}
