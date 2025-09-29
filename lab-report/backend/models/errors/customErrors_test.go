package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserNotFoundError_Error_ShouldReturnMessage(t *testing.T) {
	// Arrange
	errorMessage := "User not found in database"
	err := &UserNotFoundError{Message: errorMessage}

	// Act
	result := err.Error()

	// Assert
	assert.Equal(t, errorMessage, result)
}

func TestUserNotFoundError_EmptyMessage_ShouldReturnEmpty(t *testing.T) {
	// Arrange
	err := &UserNotFoundError{Message: ""}

	// Act
	result := err.Error()

	// Assert
	assert.Empty(t, result)
}

func TestUserNotFoundError_IsError_ShouldImplementErrorInterface(t *testing.T) {
	// Arrange
	err := &UserNotFoundError{Message: "test error"}

	// Act & Assert
	var errorInterface error = err
	assert.NotNil(t, errorInterface)
	assert.Equal(t, "test error", errorInterface.Error())
}

func TestTokenIsNullError_Error_ShouldReturnMessage(t *testing.T) {
	// Arrange
	errorMessage := "Token is null or empty"
	err := &TokenIsNullError{Message: errorMessage}

	// Act
	result := err.Error()

	// Assert
	assert.Equal(t, errorMessage, result)
}

func TestTokenIsNullError_EmptyMessage_ShouldReturnEmpty(t *testing.T) {
	// Arrange
	err := &TokenIsNullError{Message: ""}

	// Act
	result := err.Error()

	// Assert
	assert.Empty(t, result)
}

func TestTokenIsNullError_IsError_ShouldImplementErrorInterface(t *testing.T) {
	// Arrange
	err := &TokenIsNullError{Message: "token missing"}

	// Act & Assert
	var errorInterface error = err
	assert.NotNil(t, errorInterface)
	assert.Equal(t, "token missing", errorInterface.Error())
}

func TestInvalidTokenError_Error_ShouldReturnMessage(t *testing.T) {
	// Arrange
	errorMessage := "Token is invalid or expired"
	err := &InvalidTokenError{Message: errorMessage}

	// Act
	result := err.Error()

	// Assert
	assert.Equal(t, errorMessage, result)
}

func TestInvalidTokenError_EmptyMessage_ShouldReturnEmpty(t *testing.T) {
	// Arrange
	err := &InvalidTokenError{Message: ""}

	// Act
	result := err.Error()

	// Assert
	assert.Empty(t, result)
}

func TestInvalidTokenError_IsError_ShouldImplementErrorInterface(t *testing.T) {
	// Arrange
	err := &InvalidTokenError{Message: "invalid jwt token"}

	// Act & Assert
	var errorInterface error = err
	assert.NotNil(t, errorInterface)
	assert.Equal(t, "invalid jwt token", errorInterface.Error())
}

func TestCustomErrors_ErrorComparison_ShouldWork(t *testing.T) {
	// Arrange
	userNotFoundErr := &UserNotFoundError{Message: "user error"}
	tokenNullErr := &TokenIsNullError{Message: "token error"}
	invalidTokenErr := &InvalidTokenError{Message: "invalid error"}

	// Act & Assert
	assert.NotEqual(t, userNotFoundErr.Error(), tokenNullErr.Error())
	assert.NotEqual(t, tokenNullErr.Error(), invalidTokenErr.Error())
	assert.NotEqual(t, userNotFoundErr.Error(), invalidTokenErr.Error())
}

func TestCustomErrors_ErrorsIs_ShouldWork(t *testing.T) {
	// Arrange
	userErr1 := &UserNotFoundError{Message: "user not found"}
	userErr2 := &UserNotFoundError{Message: "user not found"}
	tokenErr := &TokenIsNullError{Message: "token null"}

	// Act & Assert
	// Note: errors.Is() checks for error identity/wrapping, not type equality
	assert.False(t, errors.Is(userErr1, userErr2)) // Different instances
	assert.False(t, errors.Is(userErr1, tokenErr)) // Different types
	assert.True(t, errors.Is(userErr1, userErr1))  // Same instance
}

func TestCustomErrors_TypeAssertion_ShouldWork(t *testing.T) {
	// Arrange
	var err error

	// Test UserNotFoundError
	err = &UserNotFoundError{Message: "user error"}
	userErr, ok := err.(*UserNotFoundError)
	assert.True(t, ok)
	assert.Equal(t, "user error", userErr.Message)

	// Test TokenIsNullError
	err = &TokenIsNullError{Message: "token error"}
	tokenErr, ok := err.(*TokenIsNullError)
	assert.True(t, ok)
	assert.Equal(t, "token error", tokenErr.Message)

	// Test InvalidTokenError
	err = &InvalidTokenError{Message: "invalid error"}
	invalidErr, ok := err.(*InvalidTokenError)
	assert.True(t, ok)
	assert.Equal(t, "invalid error", invalidErr.Message)
}

func TestCustomErrors_ErrorWrapping_ShouldWork(t *testing.T) {
	// Arrange
	baseErr := &UserNotFoundError{Message: "original user error"}
	wrappedErr := errors.New("wrapped: " + baseErr.Error())

	// Act & Assert
	assert.Contains(t, wrappedErr.Error(), "original user error")
	assert.Contains(t, wrappedErr.Error(), "wrapped:")
}

func TestCustomErrors_CommonErrorMessages_ShouldWork(t *testing.T) {
	// Arrange
	testCases := []struct {
		name     string
		errType  string
		message  string
		expected string
	}{
		{
			name:     "UserNotFound_StandardMessage",
			errType:  "UserNotFoundError",
			message:  "User with ID 123 not found",
			expected: "User with ID 123 not found",
		},
		{
			name:     "TokenNull_StandardMessage",
			errType:  "TokenIsNullError",
			message:  "Authorization token is required",
			expected: "Authorization token is required",
		},
		{
			name:     "InvalidToken_StandardMessage",
			errType:  "InvalidTokenError",
			message:  "JWT token is invalid or expired",
			expected: "JWT token is invalid or expired",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var err error

			switch tc.errType {
			case "UserNotFoundError":
				err = &UserNotFoundError{Message: tc.message}
			case "TokenIsNullError":
				err = &TokenIsNullError{Message: tc.message}
			case "InvalidTokenError":
				err = &InvalidTokenError{Message: tc.message}
			}

			// Act & Assert
			assert.Equal(t, tc.expected, err.Error())
		})
	}
}

func TestCustomErrors_NilPointer_ShouldPanic(t *testing.T) {
	// Arrange
	var userErr *UserNotFoundError
	var tokenErr *TokenIsNullError
	var invalidErr *InvalidTokenError

	// Act & Assert - These should panic when Error() is called on nil
	assert.Panics(t, func() { _ = userErr.Error() })
	assert.Panics(t, func() { _ = tokenErr.Error() })
	assert.Panics(t, func() { _ = invalidErr.Error() })
}

// Benchmark tests
func BenchmarkUserNotFoundError_Error(b *testing.B) {
	err := &UserNotFoundError{Message: "User not found"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func BenchmarkTokenIsNullError_Error(b *testing.B) {
	err := &TokenIsNullError{Message: "Token is null"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func BenchmarkInvalidTokenError_Error(b *testing.B) {
	err := &InvalidTokenError{Message: "Invalid token"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}
