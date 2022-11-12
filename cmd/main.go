package main

import (
	"coinScan/cmd/config"
	"coinScan/cmd/server"
	"coinScan/cmd/server/handler"
	"coinScan/internal/blockchain/ethereum"
	"coinScan/internal/blockchain/ethereum/client/infura"
	"coinScan/internal/platform"
	"fmt"
)

func main() {
	//configuration
	log := platform.NewLogrus()
	configApp := config.Load()

	//clients
	infuraCli, _ := infura.New(log, configApp)

	//service
	ethSrv := ethereum.NewService(log, configApp, &infuraCli)

	//handler
	ethHandler := handler.NewEth(log, configApp, &ethSrv)

	//server
	fiberServer := server.NewInstance(log, &ethHandler)
	fiberServer.AddEthereumRoutes()
	fiberServer.Start()
	fmt.Println(configApp)
}
