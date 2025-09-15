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

type TokenIsInvalidError struct{
	Message string
}

func (e *TokenIsInvalidError) Error() string{
	return e.Message
}