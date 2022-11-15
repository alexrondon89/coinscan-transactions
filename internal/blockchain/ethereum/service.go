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
	logger               *logrus.Logger
	config               *config.Config
	client               blockchain.Client
	cacheTransactions    []blockchain.Transaction
	numberOfTransactions uint16
}

func NewService(logger *logrus.Logger, config *config.Config, client blockchain.Client) (*Service, errors.Error) {
	s := &Service{
		logger:               logger,
		client:               client,
		config:               config,
		numberOfTransactions: config.Ethereum.Cache.NumberOfElements,
	}

	s.updateCacheLastTransactions()
	return s, nil
}

func (s *Service) updateCacheLastTransactions() {
	go func() {
		ticker := time.NewTicker(time.Duration(s.config.Ethereum.Cache.TimeToUpdate) * time.Second)
		for {
			select {
			case <-ticker.C:
				lastTransactions, err := s.client.LastTransactions(context.Background(), s.numberOfTransactions)
				s.logger.Info("amount: ", len(lastTransactions))

				if err == nil {
					s.logger.Info("new execution to update cache with the last transactions ")
					s.cacheTransactions = lastTransactions
					continue
				}
				s.logger.Error("problem refreshing cacheTransactions: " + err.Error())
			}
		}
	}()
}

func (s *Service) GetLastTransactions(c *fiber.Ctx, n uint16) ([]blockchain.Transaction, errors.Error) {
	if n > s.config.Ethereum.Cache.MaxNumberOfElements {
		s.logger.Error("error getting last transactions, max amount set exceeded")
		return nil, errors.NewError(errors.MaxNumberOfTransactionsExceeded, nil)
	}

	if len(s.cacheTransactions) == 0 || s.numberOfTransactions != n {
		cacheTransactions, err := s.client.LastTransactions(context.Background(), n)
		if err != nil {
			s.logger.Error("error recovering transactions using client")
			return nil, errors.NewError(errors.GetTransactionsError, err)
		}

		s.logger.Info("cache updated with new amount of transactions: ", n, " transactions")
		s.numberOfTransactions = n
		s.cacheTransactions = cacheTransactions
		return s.cacheTransactions, nil
	}

	s.logger.Info("last transactions recovered from cache")
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
