package internal

import (
	"github.com/alexrondon89/coinscan-transactions/internal/blockchain"
	"github.com/alexrondon89/coinscan-transactions/internal/platform/errors"
	"github.com/gofiber/fiber/v2"
)

type Service interface {
	GetLastTransactions(c *fiber.Ctx, uint162 uint16) ([]blockchain.Transaction, errors.Error)
	GetTransaction(c *fiber.Ctx, hash string) (blockchain.Transaction, errors.Error)
}
