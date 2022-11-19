package blockchain

import (
	"context"
	"github.com/alexrondon89/coinscan-transactions/internal/platform"
	"github.com/alexrondon89/coinscan-transactions/internal/platform/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"math/big"
	"time"
)

type Client interface {
	Ping()
	LastTransactions(context.Context, uint16) ([]Transaction, errors.Error)
	FindTransactionProcessed(c context.Context, hash common.Hash) (Transaction, errors.Error)
}

// Transaction todo add more details to transaction
type Transaction struct {
	Overview overview `json:"trx"`
	Block    *block   `json:"block,omitempty"`
	Values   *values  `json:"values,omitempty"`
	Fee      *Fee     `json:"Fee,omitempty"`
}

type overview struct {
	Pending bool   `json:"pending"`
	From    string `json:"from"`
	To      string `json:"to,omitempty"`
	Hash    string `json:"transactionHash"`
}

type values struct {
	TransactionFee uint64     `json:"transactionFee"`
	Value          *big.Float `json:"value"`
}

type Fee struct {
	MaxFeePerGasOffered uint64 `json:"maxFeePerGasOffered"`
	GasLimit            uint64 `json:"gasLimit"`
	GasUsed             uint64 `json:"gasUsed"`
	GasPrice            uint64 `json:"gasPrice"`
	GasTipCap           uint64 `json:"gasTipCap"`
	BaseFee             uint64 `json:"baseFee"`
}

type block struct {
	Hash      string    `json:"hash"`
	Number    uint64    `json:"number"`
	TimeStamp time.Time `json:"timeStamp"`
}

func NewTransaction() Transaction {
	return Transaction{}
}

func (t Transaction) BuildOverview(trx types.Transaction, trxAsMessage types.Message, pending bool) Transaction {
	t.Overview = overview{
		Pending: pending,
		From:    trxAsMessage.From().Hex(),
		Hash:    trx.Hash().Hex(),
	}

	if to := trx.To(); to != nil {
		t.Overview.To = to.Hex()
	}

	return t
}

func (t Transaction) BuildBlock(blockInfo *types.Header) Transaction {
	t.Block = &block{
		Hash:      blockInfo.Hash().Hex(),
		Number:    blockInfo.Number.Uint64(),
		TimeStamp: time.Unix(int64(blockInfo.Time), 0).UTC(),
	}
	return t
}

func (t Transaction) BuildFee(trx types.Transaction, receipt *types.Receipt, blockInfo *types.Header) Transaction {
	t.Fee = &Fee{
		MaxFeePerGasOffered: trx.GasFeeCap().Uint64(),          // https://ethereum.org/es/developers/docs/gas/#maxfee maximo pago ofrecido
		GasLimit:            trx.Gas(),                         //https://ethereum.org/es/developers/docs/gas/#maxfee maximo gas a consumirse por operaciones
		GasUsed:             receipt.GasUsed,                   // gas consumido por operaciones
		GasTipCap:           trx.GasTipCap().Uint64(),          // maxima propina para el minero
		BaseFee:             blockInfo.BaseFee.Uint64(),        // comision base por el bloque
		GasPrice:            calculateGasPrice(trx, blockInfo), // precio por cada unidad de gas
	}

	return t
}

func (t Transaction) BuildValues(trx types.Transaction, receipt *types.Receipt, blockInfo *types.Header) Transaction {
	t.Values = &values{
		Value: platform.ConvertToUnitDesired(trx.Value(), params.Ether),
	}

	if t.Fee != nil {
		t.Values.TransactionFee = t.Fee.GasPrice * receipt.GasUsed
		return t
	}

	t.Values.TransactionFee = calculateGasPrice(trx, blockInfo) * receipt.GasUsed
	return t
}

func calculateGasPrice(trx types.Transaction, blockInfo *types.Header) uint64 {
	baseFee := blockInfo.BaseFee.Uint64()
	gasTipCap := trx.GasTipCap().Uint64()
	gasFeeCap := trx.GasFeeCap().Uint64()
	price := baseFee + gasTipCap
	if price < gasFeeCap {
		return price
	}

	return gasFeeCap
}
