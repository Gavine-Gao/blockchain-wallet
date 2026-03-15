package blockchain_transaction_syncs

import (
	"context"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Gavine-Gao/blockchain-sync-account/config"
	"github.com/Gavine-Gao/blockchain-sync-account/database"
	"github.com/Gavine-Gao/blockchain-sync-account/rpcclient"
	"github.com/Gavine-Gao/blockchain-sync-account/rpcclient/chain-account/account"
	"github.com/Gavine-Gao/blockchain-sync-account/worker"
)

type BlockChainSync struct {
	Synchronizer *worker.BaseSynchronizer
	Deposit      *worker.Deposit
	Withdraw     *worker.Withdraw
	Internal     *worker.Internal
	FallBack     *worker.FallBack

	shutdown context.CancelCauseFunc
	stopped  atomic.Bool
}

func NewBlockChainSync(ctx context.Context, cfg *config.Config, shutdown context.CancelCauseFunc) (*BlockChainSync, error) {
	db, err := database.NewDB(ctx, cfg.MasterDB)
	if err != nil {
		log.Error("init database fail", err)
		return nil, err
	}

	log.Info("New deposit", "ChainAccountRpc", cfg.ChainAccountRpc)
	conn, err := grpc.NewClient(cfg.ChainAccountRpc, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("Connect to da retriever fail", "err", err)
		return nil, err
	}
	client := account.NewWalletAccountServiceClient(conn)
	accountClient, err := rpcclient.NewWalletChainAccountClient(context.Background(), client, cfg.ChainNode.ChainName)
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
	fallback, err := worker.NewFallBack(cfg, db, accountClient, deposit, shutdown)
	if err != nil {
		log.Error("new fallback worker fail", "err", err)
		return nil, err
	}

	out := &BlockChainSync{
		Deposit:  deposit,
		Withdraw: withdraw,
		Internal: internal,
		FallBack: fallback,
		shutdown: shutdown,
	}
	return out, nil
}

func (mcs *BlockChainSync) Start(ctx context.Context) error {
	err := mcs.Deposit.Start()
	if err != nil {
		return err
	}
	//err = mcs.Withdraw.Start()
	//if err != nil {
	//	return err
	//}
	//err = mcs.Internal.Start()
	//if err != nil {
	//	return err
	//}
	err = mcs.FallBack.Start()
	if err != nil {
		return err
	}
	return nil
}

func (mcs *BlockChainSync) Stop(ctx context.Context) error {
	err := mcs.Deposit.Close()
	if err != nil {
		return err
	}
	//err = mcs.Withdraw.Close()
	//if err != nil {
	//	return err
	//}
	//err = mcs.Internal.Close()
	//if err != nil {
	//	return err
	//}
	err = mcs.FallBack.Close()
	if err != nil {
		return err
	}
	return nil
}

func (mcs *BlockChainSync) Stopped() bool {
	return mcs.stopped.Load()
}
