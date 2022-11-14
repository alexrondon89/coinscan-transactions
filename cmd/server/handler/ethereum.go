package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	"github.com/alexrondon89/coinscan-transactions/cmd/config"
	"github.com/alexrondon89/coinscan-transactions/internal"
)

type EthereumHandler struct {
	logger  *logrus.Logger
	config  *config.Config
	Service internal.Service
}

func NewEth(logger *logrus.Logger, config *config.Config, ethereumService internal.Service) EthereumHandler {
	return EthereumHandler{
		config:  config,
		logger:  logger,
		Service: ethereumService,
	}
}

func (eh *EthereumHandler) HandlerLastTransactions(c *fiber.Ctx) error {
	trxList, err := eh.Service.GetLastTransactions(c, eh.config.Ethereum.Cache.NumberOfElements)
	if err != nil {
		return c.Status(err.StatusCode()).JSON(err.Type())
	}
	return c.Status(http.StatusOK).JSON(trxList)
}

func (eh *EthereumHandler) HandlerTransaction(c *fiber.Ctx) error {
	hash := c.Params("hash")
	trx, err := eh.Service.GetTransaction(c, hash)
	if err != nil {
		return c.Status(err.StatusCode()).JSON(err.Type())
	}
	return c.Status(http.StatusOK).JSON(trx)
}
