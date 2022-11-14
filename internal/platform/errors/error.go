package errors

type Error interface {
	Error() string
	StatusCode() int
	InternalStatusCode() string
}

func NewError(errorType Error, originalError error) Error {
	return Err{
		ErrorType:     errorType,
		OriginalError: originalError,
	}
}

type Err struct {
	OriginalError error
	ErrorType     Error
}

func (err Err) Error() string {
	return err.ErrorType.Error()
}

func (err Err) StatusCode() int {
	return err.ErrorType.StatusCode()
}

func (err Err) InternalStatusCode() string {
	return err.ErrorType.InternalStatusCode()
}
