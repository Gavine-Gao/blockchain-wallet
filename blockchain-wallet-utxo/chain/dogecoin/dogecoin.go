package dogecoin

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/log"
	"github.com/pkg/errors"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"

	"github.com/Gavine-Gao/blockchain-wallet-utxo/chain"
	"github.com/Gavine-Gao/blockchain-wallet-utxo/chain/base"
	btctypes "github.com/Gavine-Gao/blockchain-wallet-utxo/chain/bitcoin/types"
	"github.com/Gavine-Gao/blockchain-wallet-utxo/config"
	common2 "github.com/Gavine-Gao/blockchain-wallet-utxo/rpc/common"
	"github.com/Gavine-Gao/blockchain-wallet-utxo/rpc/utxo"
)

const ChainName = "Dogecoin"

// Dogecoin 主网参数
var dogeMainNetParams = chaincfg.Params{
	Net:              0xc0c0c0c0,
	PubKeyHashAddrID: 0x1e, // D 开头
	ScriptHashAddrID: 0x16, // 9 或 A 开头
	Bech32HRPSegwit:  "doge",
}

type ChainAdaptor struct {
	dogeClient     *base.BaseClient
	dogeDataClient *base.BaseDataClient
}

func NewChainAdaptor(conf *config.Config) (chain.IChainAdaptor, error) {
	baseClient, err := base.NewBaseClient(conf.WalletNode.Doge.RpcUrl, conf.WalletNode.Doge.RpcUser, conf.WalletNode.Doge.RpcPass)
	if err != nil {
		log.Error("new dogecoin rpc client fail", "err", err)
		return nil, err
	}
	baseDataClient, err := base.NewBaseDataClient(conf.WalletNode.Doge.DataApiUrl, conf.WalletNode.Doge.DataApiKey, "DOGE", "Dogecoin")
	if err != nil {
		log.Error("new dogecoin data client fail", "err", err)
		return nil, err
	}
	return &ChainAdaptor{
		dogeClient:     baseClient,
		dogeDataClient: baseDataClient,
	}, nil
}

func (c *ChainAdaptor) GetSupportChains(req *utxo.SupportChainsRequest) (*utxo.SupportChainsResponse, error) {
	return &utxo.SupportChainsResponse{
		Code:    common2.ReturnCode_SUCCESS,
		Msg:     "Support this chain",
		Support: true,
	}, nil
}

func (c *ChainAdaptor) ConvertAddress(req *utxo.ConvertAddressRequest) (*utxo.ConvertAddressResponse, error) {
	var address string
	compressedPubKeyBytes, err := hex.DecodeString(req.PublicKey)
	if err != nil {
		return &utxo.ConvertAddressResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "decode public key fail: " + err.Error(),
		}, nil
	}
	pubKeyHash := btcutil.Hash160(compressedPubKeyBytes)
	switch req.Format {
	case "p2pkh":
		p2pkhAddr, err := btcutil.NewAddressPubKeyHash(pubKeyHash, &dogeMainNetParams)
		if err != nil {
			log.Error("create p2pkh address fail", "err", err)
			return &utxo.ConvertAddressResponse{
				Code: common2.ReturnCode_ERROR,
				Msg:  "create p2pkh address fail: " + err.Error(),
			}, nil
		}
		address = p2pkhAddr.EncodeAddress()
	case "p2sh":
		witnessAddr, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, &dogeMainNetParams)
		if err != nil {
			log.Error("create witness address fail", "err", err)
			return &utxo.ConvertAddressResponse{
				Code: common2.ReturnCode_ERROR,
				Msg:  "create witness address fail: " + err.Error(),
			}, nil
		}
		script, err := txscript.PayToAddrScript(witnessAddr)
		if err != nil {
			log.Error("create p2sh address script fail", "err", err)
			return &utxo.ConvertAddressResponse{
				Code: common2.ReturnCode_ERROR,
				Msg:  "create p2sh address script fail: " + err.Error(),
			}, nil
		}
		p2shAddr, err := btcutil.NewAddressScriptHash(script, &dogeMainNetParams)
		if err != nil {
			log.Error("create p2sh address fail", "err", err)
			return &utxo.ConvertAddressResponse{
				Code: common2.ReturnCode_ERROR,
				Msg:  "create p2sh address fail: " + err.Error(),
			}, nil
		}
		address = p2shAddr.EncodeAddress()
	default:
		return &utxo.ConvertAddressResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "Do not support address type: " + req.Format,
		}, errors.New("do not support address type")
	}
	return &utxo.ConvertAddressResponse{
		Code:    common2.ReturnCode_SUCCESS,
		Msg:     "create address success",
		Address: address,
	}, nil
}

func (c *ChainAdaptor) ValidAddress(req *utxo.ValidAddressRequest) (*utxo.ValidAddressResponse, error) {
	address, err := btcutil.DecodeAddress(req.Address, &dogeMainNetParams)
	if err != nil {
		return &utxo.ValidAddressResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, nil
	}
	if !address.IsForNet(&dogeMainNetParams) {
		return &utxo.ValidAddressResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "address is not valid for dogecoin network",
		}, nil
	}
	return &utxo.ValidAddressResponse{
		Code:  common2.ReturnCode_SUCCESS,
		Msg:   "verify address success",
		Valid: true,
	}, nil
}

func (c *ChainAdaptor) GetFee(req *utxo.FeeRequest) (*utxo.FeeResponse, error) {
	gasFeeResp, err := c.dogeDataClient.GetFee()
	if err != nil {
		return &utxo.FeeResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	return &utxo.FeeResponse{
		Code:       common2.ReturnCode_SUCCESS,
		Msg:        "get dogecoin fee success",
		BestFee:    gasFeeResp.BestTransactionFee,
		BestFeeSat: gasFeeResp.BestTransactionFeeSat,
		SlowFee:    gasFeeResp.SlowGasPrice,
		NormalFee:  gasFeeResp.StandardGasPrice,
		FastFee:    gasFeeResp.RapidGasPrice,
	}, nil
}

func (c *ChainAdaptor) GetAccount(req *utxo.AccountRequest) (*utxo.AccountResponse, error) {
	balance, err := c.dogeDataClient.GetAccountBalance(req.Address)
	if err != nil {
		return &utxo.AccountResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "Get dogecoin account info fail",
		}, err
	}
	return &utxo.AccountResponse{
		Code:    common2.ReturnCode_SUCCESS,
		Msg:     "Get dogecoin account info success",
		Balance: balance.BalanceStr,
	}, nil
}

func (c *ChainAdaptor) GetUnspentOutputs(req *utxo.UnspentOutputsRequest) (*utxo.UnspentOutputsResponse, error) {
	utxoList, err := c.dogeDataClient.GetAccountUtxoList(req.Address)
	if err != nil {
		log.Error("get dogecoin utxo fail", "err", err)
		return nil, err
	}
	var utxoRetList []*utxo.UnspentOutput
	for _, utxoItem := range utxoList {
		txOutN, _ := strconv.Atoi(utxoItem.Index)
		unspentOutput := &utxo.UnspentOutput{
			TxId:          utxoItem.TxId,
			TxOutputN:     uint64(txOutN),
			Height:        utxoItem.Height,
			BlockTime:     utxoItem.BlockTime,
			Address:       utxoItem.Address,
			UnspentAmount: utxoItem.UnspentAmount,
			Confirmations: 0,
			Index:         uint64(txOutN),
		}
		utxoRetList = append(utxoRetList, unspentOutput)
	}
	return &utxo.UnspentOutputsResponse{
		Code:           common2.ReturnCode_SUCCESS,
		Msg:            "get dogecoin utxo success",
		UnspentOutputs: utxoRetList,
	}, nil
}

func (c *ChainAdaptor) GetBlockByNumber(req *utxo.BlockNumberRequest) (*utxo.BlockResponse, error) {
	blockHash, err := c.dogeClient.Client.GetBlockHash(req.Height)
	if err != nil {
		log.Error("get block hash by number fail", "err", err)
		return &utxo.BlockResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "get block hash fail",
		}, err
	}
	var params []json.RawMessage
	numBlocksJSON, _ := json.Marshal(blockHash)
	params = []json.RawMessage{numBlocksJSON}
	block, err := c.dogeClient.Client.RawRequest("getblock", params)
	if err != nil {
		return &utxo.BlockResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "get block fail: " + err.Error(),
		}, err
	}
	var resultBlock btctypes.BlockData
	err = json.Unmarshal(block, &resultBlock)
	if err != nil {
		log.Error("Unmarshal json fail", "err", err)
		return &utxo.BlockResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "unmarshal block fail: " + err.Error(),
		}, err
	}
	var txList []*utxo.TransactionList
	for _, txHash := range resultBlock.Tx {
		txList = append(txList, &utxo.TransactionList{Hash: txHash})
	}
	return &utxo.BlockResponse{
		Code:   common2.ReturnCode_SUCCESS,
		Msg:    "get block by number success",
		Height: uint64(req.Height),
		Hash:   blockHash.String(),
		TxList: txList,
	}, nil
}

func (c *ChainAdaptor) GetBlockByHash(req *utxo.BlockHashRequest) (*utxo.BlockResponse, error) {
	var params []json.RawMessage
	numBlocksJSON, _ := json.Marshal(req.Hash)
	params = []json.RawMessage{numBlocksJSON}
	block, err := c.dogeClient.Client.RawRequest("getblock", params)
	if err != nil {
		return &utxo.BlockResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "get block fail: " + err.Error(),
		}, err
	}
	var resultBlock btctypes.BlockData
	err = json.Unmarshal(block, &resultBlock)
	if err != nil {
		log.Error("Unmarshal json fail", "err", err)
		return &utxo.BlockResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "unmarshal block fail: " + err.Error(),
		}, err
	}
	var txList []*utxo.TransactionList
	for _, txHash := range resultBlock.Tx {
		txList = append(txList, &utxo.TransactionList{Hash: txHash})
	}
	return &utxo.BlockResponse{
		Code:   common2.ReturnCode_SUCCESS,
		Msg:    "get block by hash success",
		Height: resultBlock.Height,
		Hash:   req.Hash,
		TxList: txList,
	}, nil
}

func (c *ChainAdaptor) GetBlockHeaderByHash(req *utxo.BlockHeaderHashRequest) (*utxo.BlockHeaderResponse, error) {
	var params []json.RawMessage
	hashJSON, _ := json.Marshal(req.Hash)
	params = []json.RawMessage{hashJSON}
	block, err := c.dogeClient.Client.RawRequest("getblock", params)
	if err != nil {
		return &utxo.BlockHeaderResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "get block fail: " + err.Error(),
		}, err
	}
	var resultBlock btctypes.BlockData
	err = json.Unmarshal(block, &resultBlock)
	if err != nil {
		return &utxo.BlockHeaderResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "unmarshal block fail: " + err.Error(),
		}, err
	}
	return &utxo.BlockHeaderResponse{
		Code:       common2.ReturnCode_SUCCESS,
		Msg:        "get block header by hash success",
		ParentHash: resultBlock.PreviousBlockHash,
		Number:     strconv.FormatUint(resultBlock.Height, 10),
		BlockHash:  resultBlock.Hash,
		MerkleRoot: resultBlock.Merkleroot,
	}, nil
}

func (c *ChainAdaptor) GetBlockHeaderByNumber(req *utxo.BlockHeaderNumberRequest) (*utxo.BlockHeaderResponse, error) {
	blockNumber := req.Height
	if req.Height == 0 {
		latestBlock, err := c.dogeClient.Client.GetBlockCount()
		if err != nil {
			return &utxo.BlockHeaderResponse{
				Code: common2.ReturnCode_ERROR,
				Msg:  "get latest block fail",
			}, err
		}
		blockNumber = latestBlock
	}
	blockHash, err := c.dogeClient.Client.GetBlockHash(blockNumber)
	if err != nil {
		log.Error("get block hash by number fail", "err", err)
		return &utxo.BlockHeaderResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "get block hash fail",
		}, err
	}
	blockHeader, err := c.dogeClient.Client.GetBlockHeader(blockHash)
	if err != nil {
		return &utxo.BlockHeaderResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "get block header fail",
		}, err
	}
	return &utxo.BlockHeaderResponse{
		Code:       common2.ReturnCode_SUCCESS,
		Msg:        "get block header success",
		ParentHash: blockHeader.PrevBlock.String(),
		Number:     strconv.FormatInt(blockNumber, 10),
		BlockHash:  blockHash.String(),
		MerkleRoot: blockHeader.MerkleRoot.String(),
	}, nil
}

func (c *ChainAdaptor) SendTx(req *utxo.SendTxRequest) (*utxo.SendTxResponse, error) {
	r := bytes.NewReader([]byte(req.RawTx))
	var msgTx wire.MsgTx
	err := msgTx.Deserialize(r)
	if err != nil {
		return &utxo.SendTxResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	txHash, err := c.dogeClient.SendRawTransaction(&msgTx, true)
	if err != nil {
		return &utxo.SendTxResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	if strings.Compare(msgTx.TxHash().String(), txHash.String()) != 0 {
		log.Error("broadcast transaction, tx hash mismatch", "local hash", msgTx.TxHash().String(), "hash from net", txHash.String(), "signedTx", req.RawTx)
	}
	return &utxo.SendTxResponse{
		Code:   common2.ReturnCode_SUCCESS,
		Msg:    "send tx success",
		TxHash: txHash.String(),
	}, nil
}

func (c *ChainAdaptor) GetTxByAddress(req *utxo.TxAddressRequest) (*utxo.TxAddressResponse, error) {
	txListByAddress, err := c.dogeDataClient.GetTxListByAddress(req.Address, uint64(req.Page), uint64(req.Pagesize))
	if err != nil {
		return &utxo.TxAddressResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, err
	}
	var tx_list []*utxo.TxMessage
	for _, txItem := range txListByAddress.TransactionList {
		var from_addrs []*utxo.Address
		var to_addrs []*utxo.Address
		var value_list []*utxo.Value
		var direction int32
		from_addrs = append(from_addrs, &utxo.Address{Address: txItem.From})
		tx_fee := txItem.TxFee
		to_addrs = append(to_addrs, &utxo.Address{Address: txItem.To})
		value_list = append(value_list, &utxo.Value{Value: txItem.Amount})
		datetime := txItem.TransactionTime
		if strings.EqualFold(req.Address, from_addrs[0].Address) {
			direction = 0
		} else {
			direction = 1
		}
		tx := &utxo.TxMessage{
			Hash:     txItem.TxId,
			Froms:    from_addrs,
			Tos:      to_addrs,
			Values:   value_list,
			Fee:      tx_fee,
			Status:   utxo.TxStatus_Success,
			Type:     direction,
			Height:   txItem.Height,
			Datetime: datetime,
		}
		tx_list = append(tx_list, tx)
	}
	return &utxo.TxAddressResponse{
		Code: common2.ReturnCode_SUCCESS,
		Msg:  "get transaction list success",
		Tx:   tx_list,
	}, nil
}

func (c *ChainAdaptor) GetTxByHash(req *utxo.TxHashRequest) (*utxo.TxHashResponse, error) {
	tx, err := c.dogeDataClient.GetTxByHash(req.Hash)
	if err != nil {
		return nil, err
	}
	var fromAddrs []*utxo.Address
	var toAddrs []*utxo.Address
	var valueList []*utxo.Value
	for _, input := range tx.InputDetails {
		fromAddrs = append(fromAddrs, &utxo.Address{Address: input.InputHash})
	}
	for _, out := range tx.OutputDetails {
		toAddrs = append(toAddrs, &utxo.Address{Address: out.OutputHash})
		valueList = append(valueList, &utxo.Value{Value: out.Amount})
	}
	datetime := tx.TransactionTime
	txMsg := &utxo.TxMessage{
		Hash:     tx.Txid,
		Froms:    fromAddrs,
		Tos:      toAddrs,
		Values:   valueList,
		Fee:      tx.Txfee,
		Status:   utxo.TxStatus_Success,
		Type:     0,
		Height:   tx.Height,
		Datetime: datetime,
	}
	return &utxo.TxHashResponse{
		Code: common2.ReturnCode_SUCCESS,
		Msg:  "get transaction detail success",
		Tx:   txMsg,
	}, nil
}

func (c *ChainAdaptor) CreateUnSignTransaction(req *utxo.UnSignTransactionRequest) (*utxo.UnSignTransactionResponse, error) {
	txHash, buf, err := c.CalcSignHashes(req.Vin, req.Vout)
	if err != nil {
		log.Error("calc sign hashes fail", "err", err)
		return nil, err
	}
	return &utxo.UnSignTransactionResponse{
		Code:       common2.ReturnCode_SUCCESS,
		Msg:        "create un sign transaction success",
		TxData:     buf,
		SignHashes: txHash,
	}, nil
}

func (c *ChainAdaptor) BuildSignedTransaction(req *utxo.SignedTransactionRequest) (*utxo.SignedTransactionResponse, error) {
	r := bytes.NewReader(req.TxData)
	var msgTx wire.MsgTx
	err := msgTx.Deserialize(r)
	if err != nil {
		return &utxo.SignedTransactionResponse{Code: common2.ReturnCode_ERROR, Msg: err.Error()}, err
	}
	if len(req.Signatures) != len(msgTx.TxIn) {
		err = errors.New("Signature number != Txin number")
		return &utxo.SignedTransactionResponse{Code: common2.ReturnCode_ERROR, Msg: err.Error()}, err
	}
	if len(req.PublicKeys) != len(msgTx.TxIn) {
		err = errors.New("Pubkey number != Txin number")
		return &utxo.SignedTransactionResponse{Code: common2.ReturnCode_ERROR, Msg: err.Error()}, err
	}
	for i, in := range msgTx.TxIn {
		btcecPub, err2 := btcec.ParsePubKey(req.PublicKeys[i])
		if err2 != nil {
			return &utxo.SignedTransactionResponse{Code: common2.ReturnCode_ERROR, Msg: err2.Error()}, err2
		}
		var pkData []byte
		if btcec.IsCompressedPubKey(req.PublicKeys[i]) {
			pkData = btcecPub.SerializeCompressed()
		} else {
			pkData = btcecPub.SerializeUncompressed()
		}
		preTx, err2 := c.dogeClient.GetRawTransactionVerbose(&in.PreviousOutPoint.Hash)
		if err2 != nil {
			return &utxo.SignedTransactionResponse{Code: common2.ReturnCode_ERROR, Msg: err2.Error()}, err2
		}
		fromAddress, err2 := btcutil.DecodeAddress(preTx.Vout[in.PreviousOutPoint.Index].ScriptPubKey.Address, &dogeMainNetParams)
		if err2 != nil {
			return &utxo.SignedTransactionResponse{Code: common2.ReturnCode_ERROR, Msg: err2.Error()}, err2
		}
		fromPkScript, err2 := txscript.PayToAddrScript(fromAddress)
		if err2 != nil {
			return &utxo.SignedTransactionResponse{Code: common2.ReturnCode_ERROR, Msg: err2.Error()}, err2
		}
		if len(req.Signatures[i]) < 64 {
			return &utxo.SignedTransactionResponse{Code: common2.ReturnCode_ERROR, Msg: "Invalid signature length"}, errors.New("Invalid signature length")
		}
		var rScalar btcec.ModNScalar
		rScalar.SetByteSlice(req.Signatures[i][0:32])
		var sScalar btcec.ModNScalar
		sScalar.SetByteSlice(req.Signatures[i][32:64])
		btcecSig := ecdsa.NewSignature(&rScalar, &sScalar)
		sig := append(btcecSig.Serialize(), byte(txscript.SigHashAll))
		sigScript, err2 := txscript.NewScriptBuilder().AddData(sig).AddData(pkData).Script()
		if err2 != nil {
			return &utxo.SignedTransactionResponse{Code: common2.ReturnCode_ERROR, Msg: err2.Error()}, err2
		}
		msgTx.TxIn[i].SignatureScript = sigScript
		amount := dogeToSatoshi(preTx.Vout[in.PreviousOutPoint.Index].Value).Int64()
		vm, err2 := txscript.NewEngine(fromPkScript, &msgTx, i, txscript.StandardVerifyFlags, nil, nil, amount, nil)
		if err2 != nil {
			return &utxo.SignedTransactionResponse{Code: common2.ReturnCode_ERROR, Msg: err2.Error()}, err2
		}
		if err3 := vm.Execute(); err3 != nil {
			return &utxo.SignedTransactionResponse{Code: common2.ReturnCode_ERROR, Msg: err3.Error()}, err3
		}
	}
	buf := bytes.NewBuffer(make([]byte, 0, msgTx.SerializeSize()))
	err = msgTx.Serialize(buf)
	if err != nil {
		return &utxo.SignedTransactionResponse{Code: common2.ReturnCode_ERROR, Msg: err.Error()}, err
	}
	hash := msgTx.TxHash()
	return &utxo.SignedTransactionResponse{
		Code: common2.ReturnCode_SUCCESS, SignedTxData: buf.Bytes(), Hash: (&hash).CloneBytes(),
	}, nil
}

func (c *ChainAdaptor) DecodeTransaction(req *utxo.DecodeTransactionRequest) (*utxo.DecodeTransactionResponse, error) {
	res, err := c.DecodeTx(req.RawData, req.Vins, false)
	if err != nil {
		return &utxo.DecodeTransactionResponse{Code: common2.ReturnCode_ERROR, Msg: err.Error()}, err
	}
	return &utxo.DecodeTransactionResponse{
		Code: common2.ReturnCode_SUCCESS, Msg: "decode transaction success",
		SignHashes: res.SignHashes, Status: utxo.TxStatus_Other,
		Vins: res.Vins, Vouts: res.Vouts, CostFee: res.CostFee.String(),
	}, nil
}

func (c *ChainAdaptor) VerifySignedTransaction(req *utxo.VerifyTransactionRequest) (*utxo.VerifyTransactionResponse, error) {
	_, err := c.DecodeTx([]byte(""), nil, true)
	if err != nil {
		return &utxo.VerifyTransactionResponse{Code: common2.ReturnCode_ERROR, Msg: err.Error()}, err
	}
	return &utxo.VerifyTransactionResponse{Code: common2.ReturnCode_SUCCESS, Msg: "verify transaction success", Verify: true}, nil
}

func (c *ChainAdaptor) CalcSignHashes(Vins []*utxo.Vin, Vouts []*utxo.Vout) ([][]byte, []byte, error) {
	if len(Vins) == 0 || len(Vouts) == 0 {
		return nil, nil, errors.New("invalid len in or out")
	}
	rawTx := wire.NewMsgTx(wire.TxVersion)
	for _, in := range Vins {
		utxoHash, err := chainhash.NewHashFromStr(in.Hash)
		if err != nil {
			return nil, nil, err
		}
		rawTx.AddTxIn(wire.NewTxIn(wire.NewOutPoint(utxoHash, in.Index), nil, nil))
	}
	for _, out := range Vouts {
		toAddr, err := btcutil.DecodeAddress(out.Address, &dogeMainNetParams)
		if err != nil {
			return nil, nil, err
		}
		toPkScript, err := txscript.PayToAddrScript(toAddr)
		if err != nil {
			return nil, nil, err
		}
		rawTx.AddTxOut(wire.NewTxOut(out.Amount, toPkScript))
	}
	signHashes := make([][]byte, len(Vins))
	for i, in := range Vins {
		fromAddr, err := btcutil.DecodeAddress(in.Address, &dogeMainNetParams)
		if err != nil {
			return nil, nil, err
		}
		fromPkScript, err := txscript.PayToAddrScript(fromAddr)
		if err != nil {
			return nil, nil, err
		}
		sh, err := txscript.CalcSignatureHash(fromPkScript, txscript.SigHashAll, rawTx, i)
		if err != nil {
			return nil, nil, err
		}
		signHashes[i] = sh
	}
	buf := bytes.NewBuffer(make([]byte, 0, rawTx.SerializeSize()))
	_ = rawTx.Serialize(buf)
	return signHashes, buf.Bytes(), nil
}

func (c *ChainAdaptor) DecodeTx(txData []byte, vins []*utxo.Vin, sign bool) (*DecodeTxRes, error) {
	var msgTx wire.MsgTx
	if err := msgTx.Deserialize(bytes.NewReader(txData)); err != nil {
		return nil, err
	}
	offline := len(vins) > 0
	if offline && len(vins) != len(msgTx.TxIn) {
		return nil, errors.New("the length of deserialized tx's in differs from vin")
	}
	ins, totalIn, err := c.DecodeVins(msgTx, offline, vins, sign)
	if err != nil {
		return nil, err
	}
	outs, totalOut, err := c.DecodeVouts(msgTx)
	if err != nil {
		return nil, err
	}
	sh, _, err := c.CalcSignHashes(ins, outs)
	if err != nil {
		return nil, err
	}
	res := DecodeTxRes{SignHashes: sh, Vins: ins, Vouts: outs, CostFee: totalIn.Sub(totalIn, totalOut)}
	if sign {
		res.Hash = msgTx.TxHash().String()
	}
	return &res, nil
}

func (c *ChainAdaptor) DecodeVins(msgTx wire.MsgTx, offline bool, vins []*utxo.Vin, sign bool) ([]*utxo.Vin, *big.Int, error) {
	ins := make([]*utxo.Vin, 0, len(msgTx.TxIn))
	total := big.NewInt(0)
	for idx, in := range msgTx.TxIn {
		vin, err := c.GetVin(offline, vins, idx, in)
		if err != nil {
			return nil, nil, err
		}
		if sign {
			if err = c.VerifySign(vin, msgTx, idx); err != nil {
				return nil, nil, err
			}
		}
		total.Add(total, big.NewInt(vin.Amount))
		ins = append(ins, vin)
	}
	return ins, total, nil
}

func (c *ChainAdaptor) DecodeVouts(msgTx wire.MsgTx) ([]*utxo.Vout, *big.Int, error) {
	outs := make([]*utxo.Vout, 0, len(msgTx.TxOut))
	total := big.NewInt(0)
	for _, out := range msgTx.TxOut {
		_, addrs, _, err := txscript.ExtractPkScriptAddrs(out.PkScript, &dogeMainNetParams)
		if err != nil {
			return nil, nil, err
		}
		outs = append(outs, &utxo.Vout{Address: addrs[0].EncodeAddress(), Amount: out.Value})
		total.Add(total, big.NewInt(out.Value))
	}
	return outs, total, nil
}

func (c *ChainAdaptor) GetVin(offline bool, vins []*utxo.Vin, index int, in *wire.TxIn) (*utxo.Vin, error) {
	var vin *utxo.Vin
	if offline {
		vin = vins[index]
	} else {
		preTx, err := c.dogeClient.GetRawTransactionVerbose(&in.PreviousOutPoint.Hash)
		if err != nil {
			return nil, err
		}
		out := preTx.Vout[in.PreviousOutPoint.Index]
		vin = &utxo.Vin{Amount: dogeToSatoshi(out.Value).Int64(), Address: out.ScriptPubKey.Address}
	}
	vin.Hash = in.PreviousOutPoint.Hash.String()
	vin.Index = in.PreviousOutPoint.Index
	return vin, nil
}

func (c *ChainAdaptor) VerifySign(vin *utxo.Vin, msgTx wire.MsgTx, index int) error {
	fromAddr, err := btcutil.DecodeAddress(vin.Address, &dogeMainNetParams)
	if err != nil {
		return err
	}
	pkScript, err := txscript.PayToAddrScript(fromAddr)
	if err != nil {
		return err
	}
	vm, err := txscript.NewEngine(pkScript, &msgTx, index, txscript.StandardVerifyFlags, nil, nil, vin.Amount, nil)
	if err != nil {
		return err
	}
	return vm.Execute()
}
