package etherscan

import "github.com/Gavine-Gao/chain-explorer-api/common/transaction"

func (cea *ChainExplorerAdaptor) GetTxByHash(request *transaction.TxRequest) (*transaction.TxResponse, error) {
	return &transaction.TxResponse{}, nil
}
