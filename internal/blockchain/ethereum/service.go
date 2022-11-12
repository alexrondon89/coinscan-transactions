package ethereum

import (
	"coinScan/cmd/config"
	"coinScan/internal/blockchain"
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"time"
)

type Service struct {
	logger            *logrus.Logger
	config            *config.Config
	client            blockchain.Client
	cacheTransactions []blockchain.Transaction
}

func NewService(logger *logrus.Logger, config *config.Config, client blockchain.Client) Service {
	s := Service{
		logger: logger,
		client: client,
		config: config,
	}

	s.cacheTransactions, _ = s.client.LastTransactions(context.Background(), s.config.Ethereum.Cache.NumberOfElements)
	s.updateCacheLastTransactions()
	return s
}

func (s *Service) updateCacheLastTransactions() {
	go func() {
		ticker := time.NewTicker(time.Duration(s.config.Ethereum.Cache.TimeToUpdate) * time.Second)
		for {
			select {
			case <-ticker.C:
				lastTransactions, err := s.client.LastTransactions(context.Background(), s.config.Ethereum.Cache.NumberOfElements)
				if err == nil {
					s.logger.Info("cache updated with the last transactions")
					s.cacheTransactions = lastTransactions
					continue
				}
				s.logger.Error("error during cache updated")
			}
		}
	}()
}

func (s *Service) GetLastTransactions(c *fiber.Ctx, n uint16) ([]blockchain.Transaction, error) {
	trxs, err := s.client.LastTransactions(c.Context(), n)
	if err != nil {
		return nil, err
	}
	s.cacheTransactions = trxs
	return s.cacheTransactions, nil
}

func (s *Service) GetTransaction(c *fiber.Ctx, hash string) (blockchain.Transaction, error) {
	hashType := common.HexToHash(hash)
	trx, _ := s.client.FindTransaction(c.Context(), hashType)
	return trx, nil
}
