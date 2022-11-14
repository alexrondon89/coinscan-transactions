package server

import (
	"github.com/alexrondon89/coinscan-transactions/cmd/server/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"log"
)

type WebServer struct {
	logger          *logrus.Logger
	instance        *fiber.App
	ethereumHandler *handler.EthereumHandler
}

func NewInstance(logger *logrus.Logger, ethereumHandler *handler.EthereumHandler) *WebServer {
	return &WebServer{
		logger:          logger,
		instance:        fiber.New(),
		ethereumHandler: ethereumHandler,
	}
}

func (ws *WebServer) AddEthereumRoutes() {
	ws.instance.Get("/lastTransactions", ws.ethereumHandler.HandlerLastTransactions)
	ws.instance.Get("/transaction/:hash", ws.ethereumHandler.HandlerTransaction)
}

func (ws *WebServer) Start() {
	defer ws.instance.Shutdown()

	err := ws.instance.Listen(":3000")
	if err != nil {
		log.Fatal("coinScan transactions service could not start: ", err.Error())
	}
}
