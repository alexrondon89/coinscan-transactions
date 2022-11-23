package errors

import "net/http"

var (
	ConnectErr              = newClientErr("error creating client connection to blockchain", "0001CLIENT", http.StatusInternalServerError)
	QueryTransactionErr     = newClientErr("error getting transacion", "0002CLIENT", http.StatusInternalServerError)
	BlockErr                = newClientErr("error getting block from blockchain", "0003CLIENT", http.StatusInternalServerError)
	ReceiptErr              = newClientErr("receipt for transaction not exist or it is not processed yet", "0004CLIENT", http.StatusBadRequest)
	TransactionAsMessageErr = newClientErr("error getting transaction as a message", "0005CLIENT", http.StatusInternalServerError)
	ConversionUnitError     = newClientErr("error converting value to ethereum scale unit", "0006CLIENT", http.StatusInternalServerError)
)

func newClientErr(message, internalCode string, httpCode int) ErrorType {
	return ErrorType{
		Message:      message,
		InternalCode: internalCode,
		HttpCode:     httpCode,
	}
}
