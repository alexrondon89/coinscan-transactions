package handler

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	"github.com/alexrondon89/coinscan-transactions/cmd/config"
	"github.com/alexrondon89/coinscan-transactions/internal"
	"github.com/alexrondon89/coinscan-transactions/internal/platform/errors"
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
	size := c.Query("size", strconv.Itoa(int(eh.config.Ethereum.Cache.NumberOfElements)))
	sizeNumber, err := strconv.Atoi(size)
	if err != nil {
		errToReturn := errors.NewError(errors.SizeError, err)
		return c.Status(errToReturn.StatusCode()).JSON(errToReturn.Type())
	}

	trxList, errTrx := eh.Service.GetLastTransactions(c, uint16(sizeNumber))
	if errTrx != nil {
		return c.Status(errTrx.StatusCode()).JSON(errTrx.Type())
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
