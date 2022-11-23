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
		return Infura{}, errors.NewError(errors.ConnectErr, err)
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
		return nil, errors.NewError(errors.BlockErr, err)
	}

	transactions := lastBlock.Transactions()
	var trxs []blockchain.Transaction

	if nTrx := uint16(len(transactions)); n > nTrx { //todo improve searching the previous the last block to obtain the rest od transactions
		n = nTrx
		transactions = transactions[:n]
	}

	for _, transaction := range transactions[:n] {
		trxMessage, err := transaction.AsMessage(types.LatestSignerForChainID(transaction.ChainId()), nil)
		if err != nil {
			ic.logger.Error("problem converting transaction as a message type: ", err.Error())
			return nil, errors.NewError(errors.TransactionAsMessageErr, err)
		}

		trx := blockchain.NewTransaction().BuildOverview(*transaction, trxMessage, false)
		trxs = append(trxs, trx)
	}
	return trxs, nil
}

func (ic *Infura) FindTransactionProcessed(c context.Context, hash common.Hash) (blockchain.Transaction, errors.Error) {
	transac, pending, err := ic.connection.TransactionByHash(c, hash)
	if err != nil {
		ic.logger.Error("problem getting transaction: ", err.Error())
		return blockchain.Transaction{}, errors.NewError(errors.QueryTransactionErr, err)
	}
	trxAsMessage, err := transac.AsMessage(types.LatestSignerForChainID(transac.ChainId()), nil)

	if pending {
		return blockchain.NewTransaction().BuildOverview(*transac, trxAsMessage, pending), nil
	}

	receipt, err := ic.connection.TransactionReceipt(c, hash)
	if err != nil {
		ic.logger.Error("problem getting receipt for an processed transaction: ", err.Error())
		return blockchain.Transaction{}, errors.NewError(errors.ReceiptErr, err)
	}

	block, err := ic.connection.HeaderByHash(c, receipt.BlockHash)
	if err != nil {
		ic.logger.Error("problem getting block by hash: ", err.Error())
		return blockchain.Transaction{}, errors.NewError(errors.BlockErr, err)
	}

	return blockchain.NewTransaction().BuildOverview(*transac, trxAsMessage, pending).
		BuildBlock(block).
		BuildFee(*transac, block).
		BuildReceipt(*transac, receipt, block), nil
}

func (ic *Infura) Ping() {

}
