package errors

type UserNotFoundError struct {
	Message string
}

func (e *UserNotFoundError) Error() string {
	return e.Message
}

type TokenIsNullError struct {
	Message string
}

func (e *TokenIsNullError) Error() string {
	return e.Message
}

type InvalidTokenError struct {
	Message string
}

func (e *InvalidTokenError) Error() string {
	return e.Message
}
