package api

type ApiError struct {
	Message string
	Err     error
}

func (e *ApiError) Error() string {
	if e.Message == "" && e.Err == nil {
		return "Unknown Rackcorp API error."
	}
	if e.Message == "" {
		return e.Err.Error()
	}
	if e.Err == nil {
		return e.Message
	}
	return e.Message + ": " + e.Err.Error()
}

func newApiError(resp response, err error) *ApiError {
	result := &ApiError{
		Err: err,
	}
	if resp.Debug != "" {
		result.Message = resp.Debug
		return result
	}
	if resp.Message != "" {
		result.Message = resp.Message
		return result
	}
	result.Message = resp.Code
	return result
}
