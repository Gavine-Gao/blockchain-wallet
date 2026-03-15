package multichain_transaction_syncs

import (
	"context"
	"sync/atomic"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/ethereum/go-ethereum/log"

	"github.com/Gavine-Gao/blockchain-sync-btc/config"
	"github.com/Gavine-Gao/blockchain-sync-btc/database"
	"github.com/Gavine-Gao/blockchain-sync-btc/rpcclient/syncclient"
	"github.com/Gavine-Gao/blockchain-sync-btc/rpcclient/syncclient/utxo"
	"github.com/Gavine-Gao/blockchain-sync-btc/worker"
)

type BlockChainSync struct {
	Synchronizer *worker.BaseSynchronizer
	Deposit      *worker.Deposit
	Withdraw     *worker.Withdraw
	Internal     *worker.Internal

	shutdown context.CancelCauseFunc
	stopped  atomic.Bool
}

func NewBlockChainSync(ctx context.Context, cfg *config.Config, shutdown context.CancelCauseFunc) (*BlockChainSync, error) {
	db, err := database.NewDB(ctx, cfg.MasterDB)
	if err != nil {
		log.Error("init database fail", err)
		return nil, err
	}

	log.Info("New deposit", "ChainBtcRpc", cfg.ChainBtcRpc)
	conn, err := grpc.NewClient(cfg.ChainBtcRpc, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("Connect to da retriever fail", "err", err)
		return nil, err
	}
	client := utxo.NewWalletUtxoServiceClient(conn)
	accountClient, err := syncclient.NewWalletBtcAccountClient(context.Background(), client, cfg.ChainNode.ChainName)
	if err != nil {
		log.Error("new wallet account client fail", "err", err)
		return nil, err
	}

	deposit, err := worker.NewDeposit(cfg, db, accountClient, shutdown)
	if err != nil {
		log.Error("new deposit worker fail", "err", err)
		return nil, err
	}
	withdraw, err := worker.NewWithdraw(cfg, db, accountClient, shutdown)
	if err != nil {
		log.Error("new withdraw worker fail", "err", err)
		return nil, err
	}
	internal, err := worker.NewInternal(cfg, db, accountClient, shutdown)
	if err != nil {
		log.Error("new internal worker fail", "err", err)
		return nil, err
	}

	out := &BlockChainSync{
		Deposit:  deposit,
		Withdraw: withdraw,
		Internal: internal,
		shutdown: shutdown,
	}
	return out, nil
}

func (bcs *BlockChainSync) Start(ctx context.Context) error {
	err := bcs.Deposit.Start()
	if err != nil {
		return err
	}
	err = bcs.Withdraw.Start()
	if err != nil {
		return err
	}
	err = bcs.Internal.Start()
	if err != nil {
		return err
	}
	return nil
}

func (bcs *BlockChainSync) Stop(ctx context.Context) error {
	err := bcs.Deposit.Close()
	if err != nil {
		return err
	}
	err = bcs.Withdraw.Close()
	if err != nil {
		return err
	}
	err = bcs.Internal.Close()
	if err != nil {
		return err
	}
	return nil
}

func (bcs *BlockChainSync) Stopped() bool {
	return bcs.stopped.Load()
}
