package infura

import (
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"

	"github.com/alexrondon89/coinscan-transactions/internal/blockchain"
	"github.com/alexrondon89/coinscan-transactions/internal/platform"
	"github.com/alexrondon89/coinscan-transactions/internal/platform/errors"
)

type EthTransaction struct {
	trx blockchain.Transaction
}

func NewTransaction() EthTransaction {
	return EthTransaction{}
}

func (t EthTransaction) BuildOverview(trx types.Transaction, trxAsMessage types.Message, pending bool) EthTransaction {
	t.trx.Overview = blockchain.Overview{
		Pending: pending,
		From:    trxAsMessage.From().Hex(),
		Hash:    trx.Hash().Hex(),
	}

	if to := trx.To(); to != nil {
		t.trx.Overview.To = to.Hex()
	}

	return t
}

func (t EthTransaction) BuildBlock(blockInfo *types.Header) EthTransaction {
	t.trx.Block = &blockchain.Block{
		Hash:      blockInfo.Hash().Hex(),
		Number:    blockInfo.Number.Uint64(),
		TimeStamp: time.Unix(int64(blockInfo.Time), 0).UTC(),
	}
	return t
}

func (t EthTransaction) BuildFee(trx types.Transaction, blockInfo *types.Header) EthTransaction {
	t.trx.Fee = &blockchain.Fee{}
	maxFeePerGasOffered := t.buildDenomination(trx.GasFeeCap()) // https://ethereum.org/es/developers/docs/gas/#maxfee maximo pago ofrecido
	t.trx.Fee.MaxFeePerGasOffered = maxFeePerGasOffered

	gasTipCap := t.buildDenomination(trx.GasTipCap()) // maxima propina para el minero
	t.trx.Fee.GasTipCap = gasTipCap

	baseFee := t.buildDenomination(blockInfo.BaseFee) // comision base por el bloque
	t.trx.Fee.BaseFee = baseFee

	gasPrice := t.buildDenomination(t.calculateGasPrice(trx, blockInfo)) // precio por cada unidad de gas
	t.trx.Fee.GasPrice = gasPrice

	return t
}

func (t EthTransaction) BuildReceipt(trx types.Transaction, receiptOfTrx *types.Receipt, blockInfo *types.Header) EthTransaction {
	t.trx.Receipt = &blockchain.Receipt{}
	value := t.buildDenomination(trx.Value())
	t.trx.Receipt.Value = value

	gasPrice := t.calculateGasPrice(trx, blockInfo)
	trxFee := gasPrice * receiptOfTrx.GasUsed
	trxFeeAmount := t.buildDenomination(trxFee)
	t.trx.Receipt.TransactionFee = trxFeeAmount

	t.trx.Receipt.Status = receiptOfTrx.Status
	t.trx.Receipt.GasUsed = receiptOfTrx.GasUsed
	t.trx.Receipt.GasLimit = trx.Gas()

	return t
}

func (t EthTransaction) buildDenomination(unit interface{}) blockchain.Denomination {
	den := blockchain.Denomination{}

	weiAmount, err := convertUnit(unit, params.Wei)
	if err == nil {
		weiAmountFloat, _ := weiAmount.Float64()
		weiAmountString := setUnitFormat(weiAmountFloat, "0")
		den.Exp18 = weiAmountString
	}

	gweiAmount, err := convertUnit(unit, params.GWei)
	if err == nil {
		gweiAmountFloat, _ := gweiAmount.Float64()
		gweiAmountString := setUnitFormat(gweiAmountFloat, "9")
		den.Exp9 = gweiAmountString
	}

	ethAmount, err := convertUnit(unit, params.Ether)
	if err == nil {
		ethAmountFloat, _ := ethAmount.Float64()
		ethAmountString := setUnitFormat(ethAmountFloat, "18")
		den.Exp0 = ethAmountString
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

func (t EthTransaction) calculateGasPrice(trx types.Transaction, blockInfo *types.Header) uint64 {
	if gwei := t.trx.Fee.GasPrice.Exp9; gwei != "" {
		weiAmount, err := convertUnit(gwei, math.Pow(10, -9)) // due to new(big.Float).Quo() => x/y
		if err == nil {
			priceInWei, _ := weiAmount.Uint64()
			return priceInWei
		}
	}

	baseFee := blockInfo.BaseFee.Uint64()
	gasTipCap := trx.GasTipCap().Uint64()
	gasFeeCap := trx.GasFeeCap().Uint64()
	price := baseFee + gasTipCap
	if price < gasFeeCap {
		return price
	}

	return gasFeeCap
}
