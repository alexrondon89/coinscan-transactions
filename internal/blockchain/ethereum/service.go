package ethereum

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	"github.com/alexrondon89/coinscan-transactions/cmd/config"
	"github.com/alexrondon89/coinscan-transactions/internal/blockchain"
	"github.com/alexrondon89/coinscan-transactions/internal/platform/errors"
)

type Service struct {
	logger            *logrus.Logger
	config            *config.Config
	client            blockchain.Client
	cacheTransactions []blockchain.Transaction
}

func NewService(logger *logrus.Logger, config *config.Config, client blockchain.Client) (Service, errors.Error) {
	s := Service{
		logger: logger,
		client: client,
		config: config,
	}

	cacheTransactions, err := s.client.LastTransactions(context.Background(), s.config.Ethereum.Cache.NumberOfElements)
	if err != nil {
		return Service{}, errors.NewError(errors.InitializationError, err)
	}

	s.cacheTransactions = cacheTransactions
	s.updateCacheLastTransactions()
	return s, nil
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
				s.logger.Error("problem refreshing cacheTransactions: " + err.Error())
			}
		}
	}()
}

func (s *Service) GetLastTransactions(c *fiber.Ctx, n uint16) ([]blockchain.Transaction, errors.Error) {
	trxs, err := s.client.LastTransactions(c.Context(), n)
	if err != nil {
		return nil, err
	}

	s.cacheTransactions = trxs
	return s.cacheTransactions, nil
}

func (s *Service) GetTransaction(c *fiber.Ctx, hash string) (blockchain.Transaction, errors.Error) {
	hashType := common.HexToHash(hash)
	trx, err := s.client.FindTransactionProcessed(c.Context(), hashType)
	if err != nil {
		return blockchain.Transaction{}, err
	}

	return trx, nil
}
