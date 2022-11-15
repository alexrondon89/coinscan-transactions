package errors

import "net/http"

var (
	SizeError = newHandlerErr("a valid number must be passed in size query param", "0001HANDLER", http.StatusBadRequest)
)

func newHandlerErr(message, internalCode string, httpCode int) ErrorType {
	return ErrorType{
		Message:      message,
		InternalCode: internalCode,
		HttpCode:     httpCode,
	}
}
