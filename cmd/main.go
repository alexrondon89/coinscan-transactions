package main

import (
	"fmt"

	"github.com/alexrondon89/coinscan-transactions/cmd/config"
	"github.com/alexrondon89/coinscan-transactions/cmd/server"
	"github.com/alexrondon89/coinscan-transactions/cmd/server/handler"
	"github.com/alexrondon89/coinscan-transactions/internal/blockchain/ethereum"
	"github.com/alexrondon89/coinscan-transactions/internal/blockchain/ethereum/client/infura"
	"github.com/alexrondon89/coinscan-transactions/internal/platform"
)

func main() {
	//configuration
	log := platform.NewLogrus()
	configApp, err := config.Load()
	if err != nil {
		log.Fatal("coinScan transactions service could not start due to error in configApp initialization: ", err.Error())
	}

	//clients
	infuraCli, err := infura.New(log, configApp)
	if err != nil {
		log.Fatal("coinScan transactions service could not start due to error in infura client initialization: ", err.Error())
	}

	//service
	ethSrv, err := ethereum.NewService(log, configApp, &infuraCli)
	if err != nil {
		log.Fatal("coinScan transactions service could not start due to error in ethereum service initialization: ", err.Error())
	}

	//handler
	ethHandler := handler.NewEth(log, configApp, ethSrv)

	//server
	fiberServer := server.NewInstance(log, &ethHandler)
	fiberServer.AddEthereumRoutes()
	fiberServer.Start()
	fmt.Println(configApp)
}
