package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestExpiredTokens_StructValidation(t *testing.T) {
	// Arrange
	expireTime := time.Now().Add(24 * time.Hour).Unix()
	token := ExpiredTokens{
		Token:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test",
		ExpireTime: expireTime,
	}

	// Act & Assert
	assert.NotEmpty(t, token.Token, "Token should not be empty")
	assert.Greater(t, token.ExpireTime, int64(0), "ExpireTime should be greater than 0")
	assert.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test", token.Token)
	assert.Equal(t, expireTime, token.ExpireTime)
}

func TestExpiredTokens_EmptyToken(t *testing.T) {
	// Arrange
	token := ExpiredTokens{
		Token:      "",
		ExpireTime: time.Now().Unix(),
	}

	// Act & Assert
	assert.Empty(t, token.Token, "Token should be empty")
	assert.Greater(t, token.ExpireTime, int64(0), "ExpireTime should still be valid")
}

func TestExpiredTokens_ZeroExpireTime(t *testing.T) {
	// Arrange
	token := ExpiredTokens{
		Token:      "valid-token-string",
		ExpireTime: 0,
	}

	// Act & Assert
	assert.NotEmpty(t, token.Token, "Token should not be empty")
	assert.Equal(t, int64(0), token.ExpireTime, "ExpireTime should be zero")
}

func TestExpiredTokens_NegativeExpireTime(t *testing.T) {
	// Arrange
	token := ExpiredTokens{
		Token:      "valid-token-string",
		ExpireTime: -1,
	}

	// Act & Assert
	assert.NotEmpty(t, token.Token, "Token should not be empty")
	assert.Equal(t, int64(-1), token.ExpireTime, "ExpireTime should be negative")
	assert.Less(t, token.ExpireTime, int64(0), "ExpireTime should be less than 0")
}

func TestExpiredTokens_WithGormModel(t *testing.T) {
	// Arrange
	token := ExpiredTokens{
		Model: gorm.Model{
			ID: 1,
		},
		Token:      "test-jwt-token",
		ExpireTime: time.Now().Unix(),
	}

	// Act & Assert
	assert.Equal(t, uint(1), token.ID, "ID should be 1")
	assert.NotEmpty(t, token.Token, "Token should not be empty")
	assert.Greater(t, token.ExpireTime, int64(0), "ExpireTime should be positive")
}

func TestExpiredTokens_IsExpired(t *testing.T) {
	tests := []struct {
		name          string
		expireTime    int64
		currentTime   int64
		expectedResult bool
	}{
		{
			name:          "Token is expired",
			expireTime:    time.Now().Add(-1 * time.Hour).Unix(), // 1 saat önce
			currentTime:   time.Now().Unix(),
			expectedResult: true,
		},
		{
			name:          "Token is not expired",
			expireTime:    time.Now().Add(1 * time.Hour).Unix(), // 1 saat sonra
			currentTime:   time.Now().Unix(),
			expectedResult: false,
		},
		{
			name:          "Token expires exactly now",
			expireTime:    time.Now().Unix(),
			currentTime:   time.Now().Unix(),
			expectedResult: false, // Eşit zamanı expired kabul etmiyoruz
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			token := ExpiredTokens{
				Token:      "test-token",
				ExpireTime: tt.expireTime,
			}

			// Act
			isExpired := token.ExpireTime < tt.currentTime

			// Assert
			assert.Equal(t, tt.expectedResult, isExpired, tt.name)
		})
	}
}

func TestExpiredTokens_LongToken(t *testing.T) {
	// Arrange - Çok uzun token string'i test et
	longToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjE2MTYyMzkwMjJ9.invalid-signature-but-very-long-token-string-for-testing-purposes"
	
	token := ExpiredTokens{
		Token:      longToken,
		ExpireTime: time.Now().Unix(),
	}

	// Act & Assert
	assert.Equal(t, longToken, token.Token, "Long token should be preserved")
	assert.Greater(t, len(token.Token), 100, "Token should be longer than 100 characters")
}

func TestExpiredTokens_MultipleTokensComparison(t *testing.T) {
	// Arrange
	token1 := ExpiredTokens{
		Token:      "token-1",
		ExpireTime: time.Now().Unix(),
	}
	
	token2 := ExpiredTokens{
		Token:      "token-2", 
		ExpireTime: time.Now().Add(1 * time.Hour).Unix(),
	}

	// Act & Assert
	assert.NotEqual(t, token1.Token, token2.Token, "Tokens should be different")
	assert.NotEqual(t, token1.ExpireTime, token2.ExpireTime, "Expire times should be different")
	assert.Less(t, token1.ExpireTime, token2.ExpireTime, "Token1 should expire before token2")
}

func BenchmarkExpiredTokens_Creation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ExpiredTokens{
			Token:      "benchmark-token-string",
			ExpireTime: time.Now().Unix(),
		}
	}
}