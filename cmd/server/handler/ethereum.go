package handler

import (
	"github.com/alexrondon89/coinscan-transactions/internal"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	"github.com/alexrondon89/coinscan-transactions/cmd/config"
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
		return c.Status(err.StatusCode()).JSON(trxList)
	}
	return c.Status(err.StatusCode()).JSON(trxList)
}

func (eh *EthereumHandler) HandlerTransaction(c *fiber.Ctx) error {
	hash := c.Params("hash")
	trx, err := eh.Service.GetTransaction(c, hash)
	if err != nil {
		return err
	}
	return c.Status(200).JSON(trx)
}
