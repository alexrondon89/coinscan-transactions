package internal

import (
	"coinScan/internal/blockchain"
	"github.com/gofiber/fiber/v2"
)

type Service interface {
	GetLastTransactions(c *fiber.Ctx, uint162 uint16) ([]blockchain.Transaction, error)
	GetTransaction(c *fiber.Ctx, hash string) (blockchain.Transaction, error)
}
