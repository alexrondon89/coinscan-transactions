package blockchain

import (
	"context"
	"github.com/alexrondon89/coinscan-transactions/internal/platform/errors"
	"time"
)

type Client interface {
	Ping()
	LastTransactions(c context.Context, amountOfTransactions uint16) ([]Transaction, errors.Error)
	FindTransactionProcessed(c context.Context, hash string) (Transaction, errors.Error)
}

type Transaction struct {
	Overview Overview `json:"trx"`
	Block    *Block   `json:"Block,omitempty"`
	Receipt  *Receipt `json:"Receipt,omitempty"`
	Fee      *Fee     `json:"Fee,omitempty"`
}

type Overview struct {
	Pending bool   `json:"pending"`
	From    string `json:"from"`
	To      string `json:"to,omitempty"`
	Hash    string `json:"transactionHash"`
}

type Receipt struct {
	TransactionFee Denomination `json:"transactionFee"`
	Value          Denomination `json:"value"`
	Status         uint64       `json:"status"`
	GasUsed        uint64       `json:"gasUsed"`
	GasLimit       uint64       `json:"gasLimit"`
}

type Fee struct {
	MaxFeePerGasOffered Denomination `json:"maxFeePerGasOffered"`
	GasPrice            Denomination `json:"gasPrice"`
	GasTipCap           Denomination `json:"gasTipCap"`
	BaseFee             Denomination `json:"baseFee"`
}

type Block struct {
	Hash      string    `json:"hash"`
	Number    uint64    `json:"number"`
	TimeStamp time.Time `json:"timeStamp"`
}

type Denomination struct {
	Exp18 string `json:"exp18"` // wei
	Exp9  string `json:"exp9"`  // gwei
	Exp0  string `json:"exp0"`  // eth
}
