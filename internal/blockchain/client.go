package blockchain

import (
	"coinScan/internal/platform"
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"math/big"
	"time"
)

type Client interface {
	Ping()
	//Connect() *ethclient.Client
	LastTransactions(context.Context, uint16) ([]Transaction, error)
	FindTransaction(c context.Context, hash common.Hash) (Transaction, error)
}

// Transaction todo add more details to transaction
type Transaction struct {
	Block block `json:"block"`
	Trx   trx   `json:"trx"`
	Gas   gas   `json:"gas"`
}

type trx struct {
	From      string     `json:"from"`
	To        string     `json:"to"`
	Value     *big.Float `json:"value"`
	TimeStamp time.Time  `json:"timeStamp"`

	TransactionFee  uint64     `json:"transactionFee"`
	TransactionHash string     `json:"transactionHash"`
	Cost            *big.Float `json:"cost"`
}

type gas struct {
	GasLimit  uint64   `json:"gasLimit"`
	GasUsed   uint64   `json:"gasUsed"`
	GasPrice  *big.Int `json:"gasPrice"`
	GasFeeCap uint64   `json:"gasFeeCap"`
	GasTipCap uint64   `json:"gasTipCap"`
}

type block struct {
	Hash   string `json:"hash"`
	Number uint64 `json:"number"`
}

func NewTransaction() Transaction {
	return Transaction{}
}

func (t Transaction) BuildBlockSection(receipt *types.Receipt) Transaction {
	t.Block = block{
		Hash:   receipt.BlockHash.Hex(),
		Number: receipt.BlockNumber.Uint64(),
	}
	return t
}

func (t Transaction) BuildGasSection(trxMessage types.Message, receipt *types.Receipt) Transaction {
	t.Gas = gas{
		GasLimit:  trxMessage.Gas(),
		GasUsed:   receipt.GasUsed,
		GasFeeCap: trxMessage.GasFeeCap().Uint64(),
		GasTipCap: trxMessage.GasTipCap().Uint64(),
		GasPrice:  trxMessage.GasPrice(),
	}
	return t

}

func (t Transaction) BuildTrxSection(trxMessage types.Message, block types.Block) Transaction {
	t.Trx = trx{
		From:      trxMessage.From().Hex(),
		To:        trxMessage.To().Hex(),
		Value:     platform.ConvertToUnitDesired(trxMessage.Value(), params.Ether),
		TimeStamp: time.Unix(int64(block.Time()), 0).UTC(),
	}
	return t
}
