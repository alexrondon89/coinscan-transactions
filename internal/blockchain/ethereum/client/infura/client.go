package infura

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"

	"github.com/alexrondon89/coinscan-transactions/cmd/config"
	"github.com/alexrondon89/coinscan-transactions/internal/blockchain"
	"github.com/alexrondon89/coinscan-transactions/internal/platform/errors"
	"github.com/alexrondon89/coinscan-transactions/internal/platform/errors/client"
)

type Infura struct {
	logger     *logrus.Logger
	config     *config.Config
	connection *ethclient.Client
}

func New(logger *logrus.Logger, config *config.Config) (Infura, errors.Error) {
	conn, err := ethclient.Dial(config.Ethereum.Node.Url)
	if err != nil {
		logger.Error("problem getting connection: ", err.Error())
		return Infura{}, errors.NewError(client.ConnectErr, err)
	}

	return Infura{
		logger:     logger,
		config:     config,
		connection: conn,
	}, nil
}

func (ic *Infura) LastTransactions(c context.Context, n uint16) ([]blockchain.Transaction, errors.Error) {
	lastBlock, err := ic.connection.BlockByNumber(c, nil) // todo add possibility to search by block number
	if err != nil {
		ic.logger.Error("problem getting last block by number: ", err.Error())
		return nil, errors.NewError(client.BlockErr, err)
	}

	transactions := lastBlock.Transactions()
	var trxs []blockchain.Transaction

	if nTrx := uint16(len(transactions)); n > nTrx { //todo improve searching the previous the last block to obtain the rest od transactions
		transactions = transactions[:nTrx]
	}

	for _, transaction := range transactions {
		trxMessage, err := transaction.AsMessage(types.LatestSignerForChainID(transaction.ChainId()), nil)
		if err != nil {
			ic.logger.Error("problem converting transaction as a message type: ", err.Error())
			return nil, errors.NewError(client.TransactionAsMessageErr, err)
		}

		trx := blockchain.NewTransaction().BuildTrxSection(trxMessage, *lastBlock)
		trxs = append(trxs, trx)
	}
	return trxs, nil
}

func (ic *Infura) FindTransactionProcessed(c context.Context, hash common.Hash) (blockchain.Transaction, errors.Error) {
	receipt, err := ic.connection.TransactionReceipt(c, hash)
	if err != nil {
		ic.logger.Error("problem getting receipt for an processed transaction: ", err.Error())
		return blockchain.Transaction{}, errors.NewError(client.ReceiptErr, err)
	}

	block, err := ic.connection.BlockByHash(c, receipt.BlockHash)
	if err != nil {
		ic.logger.Error("problem getting block by hash: ", err.Error())
		return blockchain.Transaction{}, errors.NewError(client.BlockErr, err)
	}

	transaction := block.Transaction(hash)
	trxMessage, err := transaction.AsMessage(types.LatestSignerForChainID(transaction.ChainId()), nil)
	if err != nil {
		ic.logger.Error("problem converting transaction as a message type: ", err.Error())
		return blockchain.Transaction{}, errors.NewError(client.TransactionAsMessageErr, err)
	}

	return blockchain.NewTransaction().BuildTrxSection(trxMessage, *block).BuildGasSection(trxMessage, receipt).BuildBlockSection(receipt), nil
}

func (ic *Infura) Ping() {

}
