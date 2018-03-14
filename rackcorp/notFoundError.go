package rackcorp

type notFoundError struct {
	Message string
}

func (e *notFoundError) Error() string {
	return e.Message
}

func newNotFoundError(message string) *notFoundError {
	return &notFoundError{
		Message: message,
	}
}
