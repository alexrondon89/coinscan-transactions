package blockchain

import (
	"context"
	"fmt"
	"github.com/alexrondon89/coinscan-transactions/internal/platform"
	"github.com/alexrondon89/coinscan-transactions/internal/platform/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"math"
	"math/big"
	"strconv"
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
	Receipt  *receipt `json:"receipt,omitempty"`
	Fee      *Fee     `json:"Fee,omitempty"`
}

type overview struct {
	Pending bool   `json:"pending"`
	From    string `json:"from"`
	To      string `json:"to,omitempty"`
	Hash    string `json:"transactionHash"`
}

type receipt struct {
	TransactionFee denomination `json:"transactionFee"`
	Value          denomination `json:"value"`
	Status         uint64       `json:"status"`
	GasUsed        uint64       `json:"gasUsed"`
	GasLimit       uint64       `json:"gasLimit"`
}

type Fee struct {
	MaxFeePerGasOffered denomination `json:"maxFeePerGasOffered"`
	GasPrice            denomination `json:"gasPrice"`
	GasTipCap           denomination `json:"gasTipCap"`
	BaseFee             denomination `json:"baseFee"`
}

type block struct {
	Hash      string    `json:"hash"`
	Number    uint64    `json:"number"`
	TimeStamp time.Time `json:"timeStamp"`
}

type denomination struct {
	Wei  string
	Gwei string
	Eth  string
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

func (t Transaction) BuildFee(trx types.Transaction, blockInfo *types.Header) Transaction {
	t.Fee = &Fee{}
	maxFeePerGasOffered := t.buildDenomination(trx.GasFeeCap()) // https://ethereum.org/es/developers/docs/gas/#maxfee maximo pago ofrecido
	t.Fee.MaxFeePerGasOffered = maxFeePerGasOffered

	gasTipCap := t.buildDenomination(trx.GasTipCap()) // maxima propina para el minero
	t.Fee.GasTipCap = gasTipCap

	baseFee := t.buildDenomination(blockInfo.BaseFee) // comision base por el bloque
	t.Fee.BaseFee = baseFee

	gasPrice := t.buildDenomination(t.calculateGasPrice(trx, blockInfo)) // precio por cada unidad de gas
	t.Fee.GasPrice = gasPrice

	return t
}

func (t Transaction) BuildReceipt(trx types.Transaction, receiptOfTrx *types.Receipt, blockInfo *types.Header) Transaction {
	t.Receipt = &receipt{}
	value := t.buildDenomination(trx.Value())
	t.Receipt.Value = value

	gasPrice := t.calculateGasPrice(trx, blockInfo)
	trxFee := gasPrice * receiptOfTrx.GasUsed
	trxFeeAmount := t.buildDenomination(trxFee)
	t.Receipt.TransactionFee = trxFeeAmount

	t.Receipt.Status = receiptOfTrx.Status
	t.Receipt.GasUsed = receiptOfTrx.GasUsed
	t.Receipt.GasLimit = trx.Gas()

	return t
}

func (t Transaction) buildDenomination(unit interface{}) denomination {
	den := denomination{}

	weiAmount, err := convertUnit(unit, params.Wei)
	if err == nil {
		weiAmountFloat, _ := weiAmount.Float64()
		weiAmountString := setUnitFormat(weiAmountFloat, "0")
		den.Wei = weiAmountString
	}

	gweiAmount, err := convertUnit(unit, params.GWei)
	if err == nil {
		gweiAmountFloat, _ := gweiAmount.Float64()
		gweiAmountString := setUnitFormat(gweiAmountFloat, "9")
		den.Gwei = gweiAmountString
	}

	ethAmount, err := convertUnit(unit, params.Ether)
	if err == nil {
		ethAmountFloat, _ := ethAmount.Float64()
		ethAmountString := setUnitFormat(ethAmountFloat, "18")
		den.Eth = ethAmountString
	}

	return den
}

func convertUnit(unit interface{}, exponent float64) (*big.Float, errors.Error) {
	unitConverted, err := platform.ConvertToUnitDesired(unit, exponent)
	if err != nil {
		return nil, err
	}

	return unitConverted, nil
}

func setUnitFormat(unit float64, numberOfDecimals string) string {
	format := "%." + numberOfDecimals + "f"
	return fmt.Sprintf(format, unit)
}

func (t Transaction) calculateGasPrice(trx types.Transaction, blockInfo *types.Header) uint64 {
	if t.Fee.GasPrice.Gwei != "" {
		gasPriceInFloat, err := strconv.ParseFloat(t.Fee.GasPrice.Gwei, 64)
		if err == nil {
			priceInWei := gasPriceInFloat * math.Pow(10, 9)
			return uint64(priceInWei)
		}
	}

	baseFee := blockInfo.BaseFee.Uint64()
	fmt.Println(baseFee)
	gasTipCap := trx.GasTipCap().Uint64()
	gasFeeCap := trx.GasFeeCap().Uint64()
	price := baseFee + gasTipCap
	if price < gasFeeCap {
		return price
	}

	return gasFeeCap
}
