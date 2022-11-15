package errors

import "net/http"

var (
	InitializationError             = newServiceErr("error initialization client with transactions", "0001SERVICE", http.StatusInternalServerError)
	EmptyCacheTransactionErr        = newServiceErr("cache transactions is empty", "0002SERVICE", http.StatusInternalServerError)
	GetTransactionsError            = newServiceErr("error getting transactions from ethereum", "0003SERVICE", http.StatusInternalServerError)
	MaxNumberOfTransactionsExceeded = newServiceErr("error getting last transactions, max amount set exceeded", "0004SERVICE", http.StatusBadRequest)
)

func newServiceErr(message, internalCode string, httpCode int) ErrorType {
	return ErrorType{
		Message:      message,
		InternalCode: internalCode,
		HttpCode:     httpCode,
	}
}
