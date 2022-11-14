package service

var (
	InitializationError = NewServiceErr("error initialization client with transactions", "0001SERVICE", 500)
)

type ServiceErr struct {
	Message      string
	InternalCode string
	HttpCode     int
}

func NewServiceErr(message string, internalCode string, httpCode int) ServiceErr {
	return ServiceErr{
		Message:      message,
		InternalCode: internalCode,
		HttpCode:     httpCode,
	}
}

func (c ServiceErr) Error() string {
	return c.Message
}

func (c ServiceErr) InternalStatusCode() string {
	return c.InternalCode
}

func (c ServiceErr) StatusCode() int {
	return c.HttpCode
}
