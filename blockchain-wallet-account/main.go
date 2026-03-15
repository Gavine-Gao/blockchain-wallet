package main

import (
	"crypto/tls"
	"flag"
	"github.com/Gavine-Gao/blockchain-wallet-account/rpc/account"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/Gavine-Gao/blockchain-wallet-account/chaindispatcher"
	"github.com/Gavine-Gao/blockchain-wallet-account/config"
)

func init() {
	// 配置全局HTTP代理
	setupHTTPProxy()
}

func setupHTTPProxy() {
	// 检查环境变量中的代理设置
	if proxyURL := os.Getenv("https_proxy"); proxyURL != "" {
		proxy, err := url.Parse(proxyURL)
		if err == nil {
			// 配置优化的HTTP传输层使用代理
			http.DefaultTransport = &http.Transport{
				Proxy:                 http.ProxyURL(proxy),
				TLSClientConfig:       &tls.Config{InsecureSkipVerify: false},
				MaxIdleConns:          100,
				MaxIdleConnsPerHost:   10,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				DisableKeepAlives:     false,
			}
			log.Info("HTTP代理已配置（优化版）", "proxy", proxyURL)
		} else {
			log.Error("代理URL解析失败", "err", err)
		}
	} else {
		// 即使没有代理也优化HTTP传输层
		http.DefaultTransport = &http.Transport{
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: false},
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   10,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			DisableKeepAlives:     false,
		}
		log.Info("HTTP传输层已优化（无代理）")
	}
}

func main() {
	var f = flag.String("c", "config.yml", "config path")
	flag.Parse()
	conf, err := config.New(*f)
	if err != nil {
		panic(err)
	}
	dispatcher, err := chaindispatcher.New(conf)
	if err != nil {
		log.Error("Setup dispatcher failed", "err", err)
		panic(err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(dispatcher.Interceptor))
	defer grpcServer.GracefulStop()

	account.RegisterWalletAccountServiceServer(grpcServer, dispatcher)

	listen, err := net.Listen("tcp", ":"+conf.Server.Port)
	if err != nil {
		log.Error("net listen failed", "err", err)
		panic(err)
	}
	reflection.Register(grpcServer)

	log.Info("blockchian wallet rpc services start success", "port", conf.Server.Port)

	if err := grpcServer.Serve(listen); err != nil {
		log.Error("grpc server serve failed", "err", err)
		panic(err)
	}
}
