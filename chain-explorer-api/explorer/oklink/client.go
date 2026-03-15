package oklink

import (
	"time"

	"github.com/Gavine-Gao/chain-explorer-api/explorer"
	"github.com/Gavine-Gao/chain-explorer-api/explorer/base"
)

const ChainExplorerName = "oklink"

type ChainExplorerAdaptor struct {
	explorer.ChainExplorerAdaptor
	baseClient *base.BaseClient
}

func NewChainExplorerAdaptor(key string, baseURL string, verbose bool, timeout time.Duration) (*ChainExplorerAdaptor, error) {
	baseClient := base.NewBaseClient(key, baseURL, verbose, timeout)
	return &ChainExplorerAdaptor{
		baseClient: baseClient,
	}, nil
}
