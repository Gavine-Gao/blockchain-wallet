package bitcoincash

import (
	"math"
	"math/big"
	"strconv"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/shopspring/decimal"

	"github.com/Gavine-Gao/blockchain-wallet-utxo/rpc/utxo"
)

// BCH 主网参数（BCH 地址格式与 BTC legacy 地址相同：1开头P2PKH，3开头P2SH）
var bchMainNetParams = chaincfg.MainNetParams

const (
	bchDecimals = 8
)

type BlockData struct {
	Hash   string   `json:"hash"`
	Height uint64   `json:"height"`
	Tx     []string `json:"tx"`
}

type DecodeTxRes struct {
	Hash       string
	SignHashes [][]byte
	Vins       []*utxo.Vin
	Vouts      []*utxo.Vout
	CostFee    *big.Int
}

func bchToSatoshi(bchCount float64) *big.Int {
	amount := strconv.FormatFloat(bchCount, 'f', -1, 64)
	amountDm, _ := decimal.NewFromString(amount)
	tenDm := decimal.NewFromFloat(math.Pow(10, float64(bchDecimals)))
	satoshiDm, _ := big.NewInt(0).SetString(amountDm.Mul(tenDm).String(), 10)
	return satoshiDm
}
