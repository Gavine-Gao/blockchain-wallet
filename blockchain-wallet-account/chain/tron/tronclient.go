package tron

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	defaultRequestTimeout = 10 * time.Second
	defaultRetryCount     = 3
)

// TronClient Define a Tron RPC client
type TronClient struct {
	rpc *resty.Client
}

// DialTronClient Initialize and return a TronClient instance
func DialTronClient(rpcURL, rpcUser, rpcPass string) *TronClient {
	client := resty.New()
	client.SetHeader(rpcUser, rpcPass)
	client.SetBaseURL(rpcURL)
	client.SetTimeout(defaultRequestTimeout)
	client.SetRetryCount(defaultRetryCount)

	return &TronClient{
		rpc: client,
	}
}

// JsonRpc Call JSON-RPC
func (client *TronClient) JsonRpcBlock(params interface{}, result interface{}) error {
	// 构造请求体，确保 id_or_num 是字符串类型
	var idOrNum string
	switch v := params.(type) {
	case int64:
		idOrNum = fmt.Sprintf("\"%d\"", v) // 添加引号包裹数字
	case string:
		idOrNum = fmt.Sprintf("\"%s\"", v) // 添加引号包裹字符串
	default:
		return fmt.Errorf("unsupported params type: %T", params)
	}

	requestBody := map[string]interface{}{
		"id_or_num": json.RawMessage(idOrNum), // 使用 RawMessage 保持引号
		"detail":    true,
	}

	// 打印请求信息
	requestJSON, _ := json.MarshalIndent(requestBody, "", "  ")
	fmt.Printf("Request URL: %s\n", client.rpc.BaseURL+"/wallet/getblock")
	fmt.Printf("Request Headers:\n%v\n", client.rpc.Header)
	fmt.Printf("Request Body:\n%s\n", string(requestJSON))

	resp, err := client.rpc.R().
		SetBody(requestBody).
		SetResult(result).
		Post("/wallet/getblock")

	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}

	// 打印响应信息
	fmt.Printf("Response Status Code: %d\n", resp.StatusCode())
	fmt.Printf("Response Body:\n%s\n", string(resp.Body()))

	// 检查是否包含错误信息
	if strings.Contains(string(resp.Body()), "Error") {
		return fmt.Errorf("API error: %s", string(resp.Body()))
	}

	if resp.IsError() {
		return fmt.Errorf("API request failed with status code: %d", resp.StatusCode())
	}

	return nil
}

// JsonRpc Call JSON-RPC
func (client *TronClient) JsonRpcBlockHeader(params interface{}, result interface{}) error {
	// 构造请求体，确保 id_or_num 是字符串类型
	var idOrNum string
	switch v := params.(type) {
	case int64:
		idOrNum = fmt.Sprintf("\"%d\"", v) // 添加引号包裹数字
	case string:
		idOrNum = fmt.Sprintf("\"%s\"", v) // 添加引号包裹字符串
	default:
		return fmt.Errorf("unsupported params type: %T", params)
	}

	requestBody := map[string]interface{}{
		"id_or_num": json.RawMessage(idOrNum), // 使用 RawMessage 保持引号
		"detail":    false,
	}

	// 打印请求信息
	requestJSON, _ := json.MarshalIndent(requestBody, "", "  ")
	fmt.Printf("Request URL: %s\n", client.rpc.BaseURL+"/wallet/getblock")
	fmt.Printf("Request Headers:\n%v\n", client.rpc.Header)
	fmt.Printf("Request Body:\n%s\n", string(requestJSON))

	resp, err := client.rpc.R().
		SetBody(requestBody).
		SetResult(result).
		Post("/wallet/getblock")

	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}

	// 打印响应信息
	fmt.Printf("Response Status Code: %d\n", resp.StatusCode())
	fmt.Printf("Response Body:\n%s\n", string(resp.Body()))

	// 检查是否包含错误信息
	if strings.Contains(string(resp.Body()), "Error") {
		return fmt.Errorf("API error: %s", string(resp.Body()))
	}

	if resp.IsError() {
		return fmt.Errorf("API request failed with status code: %d", resp.StatusCode())
	}

	return nil
}

// Solidity Call Solidity
func (client *TronClient) Solidity(method string, params interface{}, result interface{}) error {
	_, err := client.rpc.R().SetBody(params).SetResult(result).Post("/walletsolidity/" + method)
	return err
}

// Wallet Call Wallet
func (client *TronClient) Wallet(method string, params interface{}, result interface{}) error {
	_, err := client.rpc.R().SetBody(params).SetResult(result).Post("/wallet/" + method)
	return err
}

// GetBlockByNumber Obtain block information based on block number
func (client *TronClient) GetBlockByNumber(blockNumber interface{}) (*BlockResponse, error) {

	var response BlockResponse
	err := client.JsonRpcBlock(blockNumber, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get block by number: %v", err)
	}
	return &response, nil
}

// GetBlockByNumber Obtain block information based on block number
func (client *TronClient) GetBlockByHush(hush string) (*Block, error) {
	params := []interface{}{hush, true}
	var response Response[Block]
	err := client.JsonRpcBlock(params, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get block by number: %v", err)
	}
	return &response.Result, nil
}

// GetBlockHeaderByNumber 获取区块头信息
func (client *TronClient) GetBlockHeaderByNumber(blockNumber int64) (*BlockResponse, error) {
	var response BlockResponse
	err := client.JsonRpcBlockHeader(blockNumber, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get block header: %v", err)
	}

	// 检查响应是否为空
	if response.BlockID == "" {
		return nil, fmt.Errorf("empty response received")
	}

	return &response, nil
}

// GetBlockByNumber Obtain block information based on block number
func (client *TronClient) GetBlockHeaderByHash(blockHush string) (*BlockResponse, error) {
	var response BlockResponse
	err := client.JsonRpcBlockHeader(blockHush, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get block header: %v", err)
	}

	// 检查响应是否为空
	if response.BlockID == "" {
		return nil, fmt.Errorf("empty response received")
	}

	return &response, nil
}

// GetBlockByHash Obtain block information based on block hash
func (client *TronClient) GetBlockByHash(blockHash string) (*Block, error) {
	params := []interface{}{blockHash, false}
	var response Response[Block]
	err := client.JsonRpcGetBlockByHash(params, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get block by hash: %v", err)

	}
	return &response.Result, nil
}
func (client *TronClient) JsonRpcGetBlockByHash(params interface{}, result interface{}) error {
	return nil
}
func (client *TronClient) JsonRpcGetBalance(params interface{}, result interface{}) error {

	requestBody := map[string]interface{}{
		"address": params, // 使用 RawMessage 保持引号
		"visible": true,
	}

	// 打印请求信息
	requestJSON, _ := json.MarshalIndent(requestBody, "", "  ")
	fmt.Printf("Request URL: %s\n", client.rpc.BaseURL+"/wallet/getblock")
	fmt.Printf("Request Headers:\n%v\n", client.rpc.Header)
	fmt.Printf("Request Body:\n%s\n", string(requestJSON))

	resp, err := client.rpc.R().
		SetBody(requestBody).
		SetResult(result).
		Post("/wallet/getaccount")

	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}

	// 打印响应信息
	fmt.Printf("Response Status Code: %d\n", resp.StatusCode())
	fmt.Printf("Response Body:\n%s\n", string(resp.Body()))

	// 检查是否包含错误信息
	if strings.Contains(string(resp.Body()), "Error") {
		return fmt.Errorf("API error: %s", string(resp.Body()))
	}

	if resp.IsError() {
		return fmt.Errorf("API request failed with status code: %d", resp.StatusCode())
	}

	return nil
}

// GetAccount Get account information
func (client *TronClient) GetBalance(address string) (*Account, error) {
	params := []interface{}{address}
	var response Account
	err := client.JsonRpcGetBalance(params, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get block by hash: %v", err)

	}
	return &response, nil

}

// GetAccount Get account information
func (client *TronClient) GetTransactionByHash(hush string) (*Transaction, error) {
	params := []interface{}{hush}
	var response Transaction
	err := client.JsonRpcGetTransactionByHash(params, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get block by hash: %v", err)

	}
	return &response, nil

}

func (client *TronClient) JsonRpcGetTransactionByHash(params interface{}, result interface{}) error {
	requestBody := map[string]interface{}{
		"value": params, // 使用 RawMessage 保持引号
	}

	// 打印请求信息
	requestJSON, _ := json.MarshalIndent(requestBody, "", "  ")
	fmt.Printf("Request URL: %s\n", client.rpc.BaseURL+"/wallet/getblock")
	fmt.Printf("Request Headers:\n%v\n", client.rpc.Header)
	fmt.Printf("Request Body:\n%s\n", string(requestJSON))

	resp, err := client.rpc.R().
		SetBody(requestBody).
		SetResult(result).
		Post("/walletsolidity/gettransactionbyid")

	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}

	// 打印响应信息
	fmt.Printf("Response Status Code: %d\n", resp.StatusCode())
	fmt.Printf("Response Body:\n%s\n", string(resp.Body()))

	// 检查是否包含错误信息
	if strings.Contains(string(resp.Body()), "Error") {
		return fmt.Errorf("API error: %s", string(resp.Body()))
	}

	if resp.IsError() {
		return fmt.Errorf("API request failed with status code: %d", resp.StatusCode())
	}

	return nil
}

// BroadcastTransaction 广播已签名交易到 TRON 网络
func (client *TronClient) BroadcastTransaction(signedTx interface{}) (*BroadcastReturns, error) {
	var result BroadcastReturns
	err := client.Wallet("broadcasttransaction", signedTx, &result)
	if err != nil {
		return nil, fmt.Errorf("broadcast transaction failed: %v", err)
	}
	return &result, nil
}

// CreateTRXTransaction 创建 TRX 转账未签名交易
func (client *TronClient) CreateTRXTransaction(from, to string, amount int64) (*UnSignTransaction, error) {
	req := &CreateTransactionRequest{
		OwnerAddress: from,
		ToAddress:    to,
		Amount:       amount,
	}
	var result UnSignTransaction
	err := client.Wallet("createtransaction", req, &result)
	if err != nil {
		return nil, fmt.Errorf("create TRX transaction failed: %v", err)
	}
	if result.TxID == "" {
		return nil, fmt.Errorf("create TRX transaction failed: empty txID")
	}
	return &result, nil
}

// CreateTRC20Transaction 创建 TRC20 代币转账未签名交易
func (client *TronClient) CreateTRC20Transaction(from, to, contractAddress string, amount int64) (*UnSignTransaction, error) {
	// 构造 transfer(address,uint256) 的 ABI 编码参数
	toHex := strings.TrimPrefix(to, "0x")
	if strings.HasPrefix(to, "T") {
		// Base58 地址需要转换
		addr, err := Base58ToHex(to)
		if err != nil {
			return nil, fmt.Errorf("convert to address failed: %v", err)
		}
		toHex = strings.TrimPrefix(addr, "0x")
	}
	// 去掉 41 前缀
	toHex = strings.TrimPrefix(toHex, "41")
	// 填充地址到 32 字节
	paddedTo := PadLeftZero(toHex, 64)
	// 填充金额到 32 字节
	amountHex := fmt.Sprintf("%x", amount)
	paddedAmount := PadLeftZero(amountHex, 64)

	parameter := paddedTo + paddedAmount

	req := &TriggerSmartContractRequest{
		OwnerAddress:     from,
		ContractAddress:  contractAddress,
		FunctionSelector: "transfer(address,uint256)",
		Parameter:        parameter,
		FeeLimit:         100000000, // 100 TRX fee limit
		CallValue:        0,
	}

	var result UnSignTrc20Transaction
	err := client.Wallet("triggersmartcontract", req, &result)
	if err != nil {
		return nil, fmt.Errorf("create TRC20 transaction failed: %v", err)
	}
	if !result.Result.Result {
		return nil, fmt.Errorf("create TRC20 transaction failed: result is false")
	}
	return &result.Transaction, nil
}
