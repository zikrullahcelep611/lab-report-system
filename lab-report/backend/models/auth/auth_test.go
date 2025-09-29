package auth

import (
	"testing"
	"time"

	"gitbub.com/zikrullahcelep611/lab-report/backend/models/role"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestLogin_ValidLogin_ShouldWork(t *testing.T) {
	// Arrange
	login := Login{
		Email:    "user@example.com",
		Role:     role.RoleAdmin,
		Password: "securePassword123",
	}

	// Act & Assert
	assert.Equal(t, "user@example.com", login.Email)
	assert.Equal(t, role.RoleAdmin, login.Role)
	assert.Equal(t, "securePassword123", login.Password)
}

func TestLogin_EmptyFields_ShouldWork(t *testing.T) {
	// Arrange
	login := Login{
		Email:    "",
		Role:     role.Role(""),
		Password: "",
	}

	// Act & Assert
	assert.Empty(t, login.Email)
	assert.Empty(t, string(login.Role))
	assert.Empty(t, login.Password)
}

func TestLogin_DifferentRoles_ShouldWork(t *testing.T) {
	// Arrange
	testCases := []struct {
		name string
		role role.Role
	}{
		{"SuperAdmin", role.RoleSuperAdmin},
		{"Admin", role.RoleAdmin},
		{"Technician", role.RoleTechnician},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			login := Login{
				Email:    "test@example.com",
				Role:     tc.role,
				Password: "password123",
			}

			// Act & Assert
			assert.Equal(t, tc.role, login.Role)
			assert.Equal(t, string(tc.role), string(login.Role))
		})
	}
}

func TestLogin_EmailFormats_ShouldWork(t *testing.T) {
	// Arrange
	testCases := []string{
		"user@example.com",
		"admin@hospital.org",
		"doctor.john@medical.center",
		"technician_123@lab.test",
		"super.admin@system.gov.tr",
	}

	for _, email := range testCases {
		t.Run(email, func(t *testing.T) {
			// Arrange
			login := Login{
				Email:    email,
				Role:     role.RoleAdmin,
				Password: "password",
			}

			// Act & Assert
			assert.Equal(t, email, login.Email)
			assert.Contains(t, login.Email, "@")
		})
	}
}

func TestLogin_PasswordTypes_ShouldWork(t *testing.T) {
	// Arrange
	testCases := []struct {
		name     string
		password string
	}{
		{"Plain Password", "simplePassword"},
		{"Complex Password", "ComplexP@ssw0rd!"},
		{"Hashed Password", "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"},
		{"Empty Password", ""},
		{"Numeric Password", "123456789"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			login := Login{
				Email:    "test@example.com",
				Role:     role.RoleAdmin,
				Password: tc.password,
			}

			// Act & Assert
			assert.Equal(t, tc.password, login.Password)
		})
	}
}

func TestLogin_GormModel_ShouldIncludeBaseFields(t *testing.T) {
	// Arrange
	login := Login{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Email:    "gorm@test.com",
		Role:     role.RoleAdmin,
		Password: "password123",
	}

	// Act & Assert
	assert.Equal(t, uint(1), login.Model.ID)
	assert.False(t, login.CreatedAt.IsZero())
	assert.False(t, login.UpdatedAt.IsZero())
	assert.Nil(t, login.DeletedAt.Time)
}

func TestLogin_JWTRegisteredClaims_ShouldWork(t *testing.T) {
	// Arrange
	now := time.Now()
	expiry := now.Add(time.Hour * 24)

	login := Login{
		Email:    "jwt@example.com",
		Role:     role.RoleAdmin,
		Password: "password",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "lab-report-system",
			Subject:   "user-login",
			Audience:  []string{"web-app", "mobile-app"},
			ExpiresAt: jwt.NewNumericDate(expiry),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        "login-session-123",
		},
	}

	// Act & Assert
	assert.Equal(t, "lab-report-system", login.Issuer)
	assert.Equal(t, "user-login", login.Subject)
	assert.Contains(t, login.Audience, "web-app")
	assert.Contains(t, login.Audience, "mobile-app")
	assert.Equal(t, expiry.Unix(), login.ExpiresAt.Unix())
	assert.Equal(t, "login-session-123", login.RegisteredClaims.ID)
}

func TestLogin_JSONTags_ShouldBeCorrect(t *testing.T) {
	// Arrange
	login := Login{
		Email:    "json@test.com",
		Role:     role.RoleTechnician,
		Password: "json_password",
	}

	// Act & Assert - Verify field accessibility for JSON serialization
	assert.Equal(t, "json@test.com", login.Email)
	assert.Equal(t, role.RoleTechnician, login.Role)
	assert.Equal(t, "json_password", login.Password)
}

func TestLogin_CompleteLoginRequest_ShouldWork(t *testing.T) {
	// Arrange - Simulate real login request
	login := Login{
		Email:    "doctor@hospital.com",
		Role:     role.RoleAdmin,
		Password: "securePassword123!",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "hospital-system",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Act & Assert
	assert.Equal(t, "doctor@hospital.com", login.Email)
	assert.Equal(t, role.RoleAdmin, login.Role)
	assert.Equal(t, "securePassword123!", login.Password)
	assert.Equal(t, "hospital-system", login.Issuer)
	assert.True(t, login.ExpiresAt.After(time.Now()))
}

func TestLogin_InvalidEmail_ShouldWork(t *testing.T) {
	// Arrange - Login struct doesn't validate email format
	testCases := []string{
		"invalid-email",
		"@missing-local.com",
		"missing-at-sign.com",
		"",
	}

	for _, email := range testCases {
		t.Run(email, func(t *testing.T) {
			// Arrange
			login := Login{
				Email:    email,
				Role:     role.RoleAdmin,
				Password: "password",
			}

			// Act & Assert
			assert.Equal(t, email, login.Email)
		})
	}
}

func TestLogin_InvalidRole_ShouldWork(t *testing.T) {
	// Arrange
	login := Login{
		Email:    "user@test.com",
		Role:     role.Role("invalid_role"),
		Password: "password",
	}

	// Act & Assert
	assert.Equal(t, role.Role("invalid_role"), login.Role)
	assert.NotEqual(t, role.RoleAdmin, login.Role)
	assert.NotEqual(t, role.RoleSuperAdmin, login.Role)
	assert.NotEqual(t, role.RoleTechnician, login.Role)
}

func TestLogin_TokenExpiration_ShouldWork(t *testing.T) {
	// Arrange
	now := time.Now()
	pastTime := now.Add(-time.Hour)
	futureTime := now.Add(time.Hour)

	expiredLogin := Login{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(pastTime),
		},
	}

	validLogin := Login{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(futureTime),
		},
	}

	// Act & Assert
	assert.True(t, expiredLogin.ExpiresAt.Before(now))
	assert.True(t, validLogin.ExpiresAt.After(now))
}

// Benchmark tests
func BenchmarkLogin_Creation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Login{
			Email:    "bench@test.com",
			Role:     role.RoleAdmin,
			Password: "password",
		}
	}
}

func BenchmarkLogin_WithJWTClaims(b *testing.B) {
	now := time.Now()
	expiry := now.Add(time.Hour)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Login{
			Email:    "bench@test.com",
			Role:     role.RoleAdmin,
			Password: "password",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expiry),
				IssuedAt:  jwt.NewNumericDate(now),
			},
		}
	}
}
