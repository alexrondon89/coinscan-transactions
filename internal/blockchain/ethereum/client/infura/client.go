package infura

import (
	"coinScan/cmd/config"
	"coinScan/internal/blockchain"
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

type Infura struct {
	logger     *logrus.Logger
	config     *config.Config
	connection *ethclient.Client
}

func New(logger *logrus.Logger, config *config.Config) (Infura, error) {
	conn, err := ethclient.Dial(config.Ethereum.Node.Url)
	if err != nil {
		return Infura{}, err
	}

	return Infura{
		logger:     logger,
		config:     config,
		connection: conn,
	}, nil
}

func (ic *Infura) LastTransactions(c context.Context, n uint16) ([]blockchain.Transaction, error) {
	lastBlock, _ := ic.connection.BlockByNumber(c, nil) // todo add possibility to search by block number
	transactions := lastBlock.Transactions()
	var trxs []blockchain.Transaction

	if nTrx := uint16(len(transactions)); n > nTrx { //todo improve searching the previou of the last block to obtain the rest od transactions
		transactions = transactions[:nTrx]
	}

	for _, transaction := range transactions {
		trxMessage, err := transaction.AsMessage(types.LatestSignerForChainID(transaction.ChainId()), nil)
		if err != nil {
			ic.logger.Error("errors getting las transactions from infura: %s", err)
			return nil, err
		}
		trx := blockchain.NewTransaction().BuildTrxSection(trxMessage, *lastBlock)
		trxs = append(trxs, trx)
	}
	return trxs, nil
}

func (ic *Infura) FindTransaction(c context.Context, hash common.Hash) (blockchain.Transaction, error) {
	receipt, _ := ic.connection.TransactionReceipt(c, hash)
	block, _ := ic.connection.BlockByHash(c, receipt.BlockHash)
	transaction := block.Transaction(hash)
	trxMessage, _ := transaction.AsMessage(types.LatestSignerForChainID(transaction.ChainId()), nil)

	return blockchain.NewTransaction().BuildTrxSection(trxMessage, *block).BuildGasSection(trxMessage, receipt).BuildBlockSection(receipt), nil
}

func (ic *Infura) Ping() {

}
