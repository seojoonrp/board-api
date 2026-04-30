package apperror

type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string { return e.Message }
func (e *AppError) Unwrap() error { return e.Err }

func NewBadRequest(msg string) *AppError {
	return &AppError{Code: 400, Message: msg}
}

func NewUnauthorized(msg string) *AppError {
	return &AppError{Code: 401, Message: msg}
}

func NewNotFound(msg string) *AppError {
	return &AppError{Code: 404, Message: msg}
}

func NewInternal(err error) *AppError {
	return &AppError{Code: 500, Message: "internal server error", Err: err}
}
