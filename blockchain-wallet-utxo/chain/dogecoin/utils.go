package dogecoin

import (
	"math"
	"math/big"
	"strconv"

	"github.com/shopspring/decimal"

	"github.com/Gavine-Gao/blockchain-wallet-utxo/rpc/utxo"
)

const (
	dogeDecimals = 8
)

type DecodeTxRes struct {
	Hash       string
	SignHashes [][]byte
	Vins       []*utxo.Vin
	Vouts      []*utxo.Vout
	CostFee    *big.Int
}

func dogeToSatoshi(dogeCount float64) *big.Int {
	amount := strconv.FormatFloat(dogeCount, 'f', -1, 64)
	amountDm, _ := decimal.NewFromString(amount)
	tenDm := decimal.NewFromFloat(math.Pow(10, float64(dogeDecimals)))
	satoshiDm, _ := big.NewInt(0).SetString(amountDm.Mul(tenDm).String(), 10)
	return satoshiDm
}
