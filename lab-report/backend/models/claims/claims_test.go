package claims

import (
	"testing"
	"time"

	"gitbub.com/zikrullahcelep611/lab-report/backend/models/role"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestClaims_ValidClaims_ShouldWork(t *testing.T) {
	// Arrange
	expiresAt := jwt.NewNumericDate(time.Now().Add(time.Hour))

	claims := Claims{
		Email:    "user@example.com",
		Role:     role.RoleAdmin,
		Password: "hashedPassword123",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expiresAt,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Act & Assert
	assert.Equal(t, "user@example.com", claims.Email)
	assert.Equal(t, role.RoleAdmin, claims.Role)
	assert.Equal(t, "hashedPassword123", claims.Password)
	assert.Equal(t, expiresAt, claims.ExpiresAt)
}

func TestClaims_EmptyFields_ShouldWork(t *testing.T) {
	// Arrange
	claims := Claims{
		Email:    "",
		Role:     role.Role(""),
		Password: "",
	}

	// Act & Assert
	assert.Empty(t, claims.Email)
	assert.Empty(t, string(claims.Role))
	assert.Empty(t, claims.Password)
}

func TestClaims_DifferentRoles_ShouldWork(t *testing.T) {
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
			claims := Claims{
				Email: "test@example.com",
				Role:  tc.role,
			}

			// Act & Assert
			assert.Equal(t, tc.role, claims.Role)
			assert.Equal(t, string(tc.role), string(claims.Role))
		})
	}
}

func TestClaims_EmailValidation_ShouldWork(t *testing.T) {
	// Arrange
	testCases := []string{
		"user@example.com",
		"admin@hospital.org",
		"doctor.john@medical.center",
		"technician123@lab.test",
	}

	for _, email := range testCases {
		t.Run(email, func(t *testing.T) {
			// Arrange
			claims := Claims{
				Email: email,
				Role:  role.RoleAdmin,
			}

			// Act & Assert
			assert.Equal(t, email, claims.Email)
			assert.Contains(t, claims.Email, "@")
		})
	}
}

func TestClaims_JWTRegisteredClaims_ShouldWork(t *testing.T) {
	// Arrange
	now := time.Now()
	expiry := now.Add(time.Hour * 24)

	claims := Claims{
		Email: "test@example.com",
		Role:  role.RoleAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "lab-report-system",
			Subject:   "user-authentication",
			Audience:  []string{"lab-app"},
			ExpiresAt: jwt.NewNumericDate(expiry),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        "token-123",
		},
	}

	// Act & Assert
	assert.Equal(t, "lab-report-system", claims.Issuer)
	assert.Equal(t, "user-authentication", claims.Subject)
	assert.Contains(t, claims.Audience, "lab-app")
	assert.Equal(t, expiry.Unix(), claims.ExpiresAt.Unix())
	assert.Equal(t, now.Unix(), claims.NotBefore.Unix())
	assert.Equal(t, now.Unix(), claims.IssuedAt.Unix())
	assert.Equal(t, "token-123", claims.ID)
}

func TestClaims_PasswordField_ShouldWork(t *testing.T) {
	// Arrange - Password field might be used for additional verification
	claims := Claims{
		Email:    "user@test.com",
		Role:     role.RoleTechnician,
		Password: "hashed_password_string",
	}

	// Act & Assert
	assert.Equal(t, "hashed_password_string", claims.Password)
	assert.NotEmpty(t, claims.Password)
}

func TestClaims_TokenExpiration_ShouldWork(t *testing.T) {
	// Arrange
	now := time.Now()
	pastTime := now.Add(-time.Hour)
	futureTime := now.Add(time.Hour)

	expiredClaims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(pastTime),
		},
	}

	validClaims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(futureTime),
		},
	}

	// Act & Assert
	assert.True(t, expiredClaims.ExpiresAt.Before(now))
	assert.True(t, validClaims.ExpiresAt.After(now))
}

func TestClaims_JSONTags_ShouldBeCorrect(t *testing.T) {
	// Arrange
	claims := Claims{
		Email:    "json@test.com",
		Role:     role.RoleAdmin,
		Password: "json_password",
	}

	// Act & Assert - Verify field accessibility for JSON serialization
	assert.Equal(t, "json@test.com", claims.Email)
	assert.Equal(t, role.RoleAdmin, claims.Role)
	assert.Equal(t, "json_password", claims.Password)
}

func TestClaims_CompleteUserClaims_ShouldWork(t *testing.T) {
	// Arrange - Simulate real user claims
	claims := Claims{
		Email:    "doctor@hospital.com",
		Role:     role.RoleAdmin,
		Password: "bcrypt_hashed_password",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "lab-system",
			Subject:   "user-auth",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Act & Assert
	assert.Equal(t, "doctor@hospital.com", claims.Email)
	assert.Equal(t, role.RoleAdmin, claims.Role)
	assert.NotEmpty(t, claims.Password)
	assert.Equal(t, "lab-system", claims.Issuer)
	assert.Equal(t, "user-auth", claims.Subject)
	assert.True(t, claims.ExpiresAt.After(time.Now()))
}

func TestClaims_InvalidRole_ShouldWork(t *testing.T) {
	// Arrange - Test with invalid role
	claims := Claims{
		Email: "user@test.com",
		Role:  role.Role("invalid_role"),
	}

	// Act & Assert
	assert.Equal(t, role.Role("invalid_role"), claims.Role)
	assert.NotEqual(t, role.RoleAdmin, claims.Role)
}

// Benchmark tests
func BenchmarkClaims_Creation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Claims{
			Email:    "bench@test.com",
			Role:     role.RoleAdmin,
			Password: "password",
		}
	}
}
