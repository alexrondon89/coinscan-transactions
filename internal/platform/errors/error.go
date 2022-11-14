package errors

type Error interface {
	Error() string
	StatusCode() int
	InternalCode() string
	Type() ErrorType
}

func NewError(errorType ErrorType, originalError error) Error {
	return Err{
		ErrorType:     errorType,
		OriginalError: originalError,
	}
}

type Err struct {
	OriginalError error
	ErrorType     ErrorType
}

type ErrorType struct {
	Message      string `json:"message"`
	InternalCode string `json:"codeError"`
	HttpCode     int    `json:"codeHttp"`
}

func (err Err) Error() string {
	return err.ErrorType.Message
}

func (err Err) StatusCode() int {
	return err.ErrorType.HttpCode
}

func (err Err) InternalCode() string {
	return err.ErrorType.InternalCode
}

func (err Err) Type() ErrorType {
	return err.ErrorType
}
