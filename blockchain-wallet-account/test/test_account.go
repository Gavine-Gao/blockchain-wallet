package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	"github.com/Gavine-Gao/blockchain-wallet-account/rpc/account"
)

func main() {
	// 连接到gRPC服务器（优化配置）
	conn, err := grpc.Dial("127.0.0.1:8189", 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                30 * time.Second,
			Timeout:             5 * time.Second,
			PermitWithoutStream: true,
		}),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(4*1024*1024), // 4MB
			grpc.MaxCallSendMsgSize(4*1024*1024), // 4MB
		),
	)
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()

	// 创建客户端
	client := account.NewWalletAccountServiceClient(conn)

	// 测试getAccount接口
	fmt.Println("=== 测试 getAccount 接口 ===")

	start := time.Now()
	resp, err := client.GetAccount(context.Background(), &account.AccountRequest{
		ConsumerToken:   "test-token",
		Chain:           "Ethereum",
		Coin:            "ETH",
		Network:         "mainnet",
		Address:         "0x8ba1f109551bD432803012645Hac136c", // Holesky测试网地址
		ContractAddress: "0x00",
	})
	duration := time.Since(start)

	if err != nil {
		fmt.Printf("❌ 请求失败: %v\n", err)
	} else {
		fmt.Printf("✅ 请求成功\n")
		fmt.Printf("响应: %+v\n", resp)
	}
	fmt.Printf("⏱️  请求耗时: %v\n", duration)

	// 如果耗时超过5秒，标记为慢
	if duration > 5*time.Second {
		fmt.Printf("⚠️  警告: 请求耗时过长 (>5秒)\n")
	}
}
