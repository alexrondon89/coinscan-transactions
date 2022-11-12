package server

import (
	"coinScan/cmd/server/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type WebServer struct {
	logger          *logrus.Logger
	instance        *fiber.App
	ethereumHandler *handler.EthereumHandler
	coinGecoService interface{}
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

	//ws.instance.Get("/blocksAmount", ws.ethereumHandler.HandlerBlockAmounts)
}

func (ws *WebServer) AddCoinGecoRoutes() {
	ws.instance.Get("/coinPrices", func(c *fiber.Ctx) error {
		return c.SendString("here will go coingeco service")
	})
}

func (ws *WebServer) Start() {
	defer ws.instance.Shutdown()

	err := ws.instance.Listen(":3000")
	if err != nil {
		panic("server could not start...")
	}

}
