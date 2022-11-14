package errors

var (
	InitializationError = newServiceErr("error initialization client with transactions", "0001SERVICE", 500)
)

func newServiceErr(message, internalCode string, httpCode int) ErrorType {
	return ErrorType{
		Message:      message,
		InternalCode: internalCode,
		HttpCode:     httpCode,
	}
}
