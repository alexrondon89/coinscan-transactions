package infura

import (
	"context"
	"fmt"
	"math/big"

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

func (ic *Infura) FindTransactions(c context.Context, numberOfTransactionsRequested uint16, blockNumber *big.Int) ([]blockchain.Transaction, errors.Error) {
	block, err := ic.FindBlockByNumber(c, blockNumber)
	if err != nil {
		return nil, err
	}

	transactions := block.Transactions()
	lengthOfBlock := uint16(transactions.Len())
	var trxs []blockchain.Transaction

	if isNecessaryToFindInPreviousBlock(numberOfTransactionsRequested, lengthOfBlock) {
		missingTransactions := numberOfTransactionsRequested - lengthOfBlock
		previousBlock := big.NewInt(block.Number().Int64() - 1)
		ic.logger.Info(fmt.Sprintf("recursive call: is necessary to get %d missing Transactions in previous block %d", missingTransactions, previousBlock))

		missingTransactionsFound, err := ic.FindTransactions(c, missingTransactions, previousBlock)
		if err != nil {
			ic.logger.Error("problem getting transaction in recursive execution: ", err.Error())
			return nil, errors.NewError(errors.GetTransactionsError, err)
		}

		trxs = append(trxs, missingTransactionsFound...)
	}

	amountOfTrxToIterate := NumberOfTransactionsToObtainInCurrentBlock(numberOfTransactionsRequested, lengthOfBlock)
	for index := 0; uint16(index) < amountOfTrxToIterate; index++ {
		transaction := transactions[index]
		trxMessage, err := transaction.AsMessage(types.LatestSignerForChainID(transaction.ChainId()), nil)
		if err != nil {
			ic.logger.Error("problem converting transaction as a message type: ", err.Error())
			return nil, errors.NewError(errors.TransactionAsMessageErr, err)
		}

		trx := NewTransaction().BuildOverview(*transaction, trxMessage, false).BuildBlock(block.Header()).trx
		trxs = append([]blockchain.Transaction{trx}, trxs...) // to get the newest transaction in first position
	}

	ic.logger.Info("transactions recovered: ", len(trxs))
	return trxs, nil
}

func (ic *Infura) FindTransaction(c context.Context, hash string) (blockchain.Transaction, errors.Error) {
	hashType := common.HexToHash(hash)
	transac, pending, err := ic.connection.TransactionByHash(c, hashType)
	if err != nil {
		ic.logger.Error("problem getting transaction: ", err.Error())
		return blockchain.Transaction{}, errors.NewError(errors.QueryTransactionErr, err)
	}
	trxAsMessage, err := transac.AsMessage(types.LatestSignerForChainID(transac.ChainId()), nil)

	if pending {
		return NewTransaction().BuildOverview(*transac, trxAsMessage, pending).trx, nil
	}

	receipt, err := ic.connection.TransactionReceipt(c, hashType)
	if err != nil {
		ic.logger.Error("problem getting receipt for an processed transaction: ", err.Error())
		return blockchain.Transaction{}, errors.NewError(errors.ReceiptErr, err)
	}

	block, err := ic.connection.HeaderByHash(c, receipt.BlockHash)
	if err != nil {
		ic.logger.Error("problem getting block by hash: ", err.Error())
		return blockchain.Transaction{}, errors.NewError(errors.BlockErr, err)
	}

	return NewTransaction().BuildOverview(*transac, trxAsMessage, pending).
		BuildBlock(block).
		BuildFee(*transac, block).
		BuildReceipt(*transac, receipt, block).trx, nil
}

func (ic *Infura) FindBlockByNumber(c context.Context, blockNumber *big.Int) (*types.Block, errors.Error) {
	if blockNumber != nil {
		block, err := ic.connection.BlockByNumber(c, big.NewInt(blockNumber.Int64()))
		if err != nil {
			ic.logger.Error("problem getting last block: ", err.Error())
			return block, errors.NewError(errors.BlockErr, err)
		}

		return block, nil
	}

	number, err := ic.connection.BlockNumber(c)
	if err != nil {
		ic.logger.Error("problem getting last block number: ", err.Error())
		return nil, errors.NewError(errors.BlockErr, err)
	}

	block, err := ic.connection.BlockByNumber(c, big.NewInt(int64(number)))
	if err != nil {
		ic.logger.Error("problem getting last block: ", err.Error())
		return nil, errors.NewError(errors.BlockErr, err)
	}

	return block, nil
}

func isNecessaryToFindInPreviousBlock(numberOfTransactionsRequested, lengthOfBlock uint16) bool {
	return numberOfTransactionsRequested > lengthOfBlock
}

func NumberOfTransactionsToObtainInCurrentBlock(numberOfTransactionsRequested, lengthOfBlock uint16) uint16 {
	if numberOfTransactionsRequested > lengthOfBlock {
		return lengthOfBlock
	}

	return numberOfTransactionsRequested
}
