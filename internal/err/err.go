package err

type ApiError struct {
	Type string
	Err  error
}

func (e *ApiError) Error() string {
	return e.Err.Error()
}

func (e *ApiError) Unwrap() error {
	return e.Err
}

const (
	ErrorTypePasswordHash = "password_hash_failure"
	ErrorTypeEmailExists  = "email_already_exists"
)

func NewPasswordHashError(err error) *ApiError {
	return &ApiError{
		Type: ErrorTypePasswordHash,
		Err:  err,
	}
}

func NewEmailExistsError(err error) *ApiError {
	return &ApiError{
		Type: ErrorTypeEmailExists,
		Err:  err,
	}
}
