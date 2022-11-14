package client

var (
	ConnectErr              = NewClientErr("error creating client connection to blockchain", "0001CLIENT", 500)
	QueryTransactionErr     = NewClientErr("error getting transacion", "0002CLIENT", 500)
	BlockErr                = NewClientErr("error getting block from blockchain", "0003CLIENT", 500)
	ReceiptErr              = NewClientErr("error getting receipt for processed transaction", "0004CLIENT", 500)
	TransactionAsMessageErr = NewClientErr("error getting transaction as a message", "0005CLIENT", 500)
)

type ClientErr struct {
	Message      string `json:"message"`
	InternalCode string `json:"codeError"`
	HttpCode     int    `json:"codeHttp"`
}

func NewClientErr(message string, internalCode string, httpCode int) ClientErr {
	return ClientErr{
		Message:      message,
		InternalCode: internalCode,
		HttpCode:     httpCode,
	}
}

func (c ClientErr) Error() string {
	return c.Message
}

func (c ClientErr) InternalStatusCode() string {
	return c.InternalCode
}

func (c ClientErr) StatusCode() int {
	return c.HttpCode
}
