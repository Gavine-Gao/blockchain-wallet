package dash

import (
	"math"
	"math/big"
	"strconv"

	"github.com/shopspring/decimal"

	"github.com/Gavine-Gao/blockchain-wallet-utxo/rpc/utxo"
)

const (
	dashDecimals = 8
)

type DecodeTxRes struct {
	Hash       string
	SignHashes [][]byte
	Vins       []*utxo.Vin
	Vouts      []*utxo.Vout
	CostFee    *big.Int
}

func dashToSatoshi(dashCount float64) *big.Int {
	amount := strconv.FormatFloat(dashCount, 'f', -1, 64)
	amountDm, _ := decimal.NewFromString(amount)
	tenDm := decimal.NewFromFloat(math.Pow(10, float64(dashDecimals)))
	satoshiDm, _ := big.NewInt(0).SetString(amountDm.Mul(tenDm).String(), 10)
	return satoshiDm
}
