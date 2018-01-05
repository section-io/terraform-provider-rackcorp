package rackcorp

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

func newNotFoundError(message string) *NotFoundError {
	return &NotFoundError{
		Message: message,
	}
}
