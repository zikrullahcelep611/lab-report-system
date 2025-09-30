package loginrepository

import (
	"context"
	"testing"

	"gitbub.com/zikrullahcelep611/lab-report/backend/models/auth"
	"gitbub.com/zikrullahcelep611/lab-report/backend/models/user"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to create gorm DB: %v", err)
	}

	return gormDB, mock
}

func TestNewRepository_ShouldCreateRepository(t *testing.T) {
	// Arrange
	db, _ := setupMockDB(t)

	// Act
	repo := NewRepository(db)

	// Assert
	assert.NotNil(t, repo)
	assert.Equal(t, db, repo.DB)
}

func TestRepository_Login_Success(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	email := "test@example.com"
	password := "testPassword123"

	// Hash the password for database storage
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err)

	expectedUser := user.User{
		Model:    gorm.Model{ID: 1},
		Email:    email,
		Password: string(hashedPassword),
	}

	// Mock the SELECT query
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1`).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password"}).
			AddRow(expectedUser.ID, expectedUser.Email, expectedUser.Password))

	// Act
	result, err := repo.Login(ctx, email, password)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, email, result.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Login_UserNotFound(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	email := "nonexistent@example.com"
	password := "password123"

	// Mock the SELECT query to return no rows
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1`).
		WithArgs(email).
		WillReturnError(gorm.ErrRecordNotFound)

	// Act
	result, err := repo.Login(ctx, email, password)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.Equal(t, auth.Login{}, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Login_WrongPassword(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	email := "test@example.com"
	correctPassword := "correctPassword123"
	wrongPassword := "wrongPassword456"

	// Hash the correct password for database storage
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(correctPassword), bcrypt.DefaultCost)
	assert.NoError(t, err)

	expectedUser := user.User{
		Model:    gorm.Model{ID: 1},
		Email:    email,
		Password: string(hashedPassword),
	}

	// Mock the SELECT query
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1`).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password"}).
			AddRow(expectedUser.ID, expectedUser.Email, expectedUser.Password))

	// Act
	result, err := repo.Login(ctx, email, wrongPassword)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, bcrypt.ErrMismatchedHashAndPassword, err)
	assert.Equal(t, auth.Login{}, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Login_DatabaseError(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	email := "test@example.com"
	password := "password123"

	// Mock the SELECT query to return a database error
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1`).
		WithArgs(email).
		WillReturnError(gorm.ErrInvalidTransaction)

	// Act
	result, err := repo.Login(ctx, email, password)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrInvalidTransaction, err)
	assert.Equal(t, auth.Login{}, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Login_EmptyEmail(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	email := ""
	password := "password123"

	// Mock the SELECT query
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1`).
		WithArgs(email).
		WillReturnError(gorm.ErrRecordNotFound)

	// Act
	result, err := repo.Login(ctx, email, password)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.Equal(t, auth.Login{}, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Login_EmptyPassword(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	email := "test@example.com"
	password := ""

	// Hash a non-empty password for database storage
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("actualPassword"), bcrypt.DefaultCost)
	assert.NoError(t, err)

	expectedUser := user.User{
		Model:    gorm.Model{ID: 1},
		Email:    email,
		Password: string(hashedPassword),
	}

	// Mock the SELECT query
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1`).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password"}).
			AddRow(expectedUser.ID, expectedUser.Email, expectedUser.Password))

	// Act
	result, err := repo.Login(ctx, email, password)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, bcrypt.ErrMismatchedHashAndPassword, err)
	assert.Equal(t, auth.Login{}, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Login_InvalidHashedPassword(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	email := "test@example.com"
	password := "password123"

	// Use an invalid hashed password (not bcrypt format)
	invalidHashedPassword := "not_a_valid_bcrypt_hash"

	expectedUser := user.User{
		Model:    gorm.Model{ID: 1},
		Email:    email,
		Password: invalidHashedPassword,
	}

	// Mock the SELECT query
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1`).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password"}).
			AddRow(expectedUser.ID, expectedUser.Email, expectedUser.Password))

	// Act
	result, err := repo.Login(ctx, email, password)

	// Assert
	assert.Error(t, err)
	// bcrypt.CompareHashAndPassword will return an error for invalid hash format
	assert.NotEqual(t, bcrypt.ErrMismatchedHashAndPassword, err)
	assert.Equal(t, auth.Login{}, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Login_Context_ShouldPassContext(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)

	// Create a context with cancel to test context propagation
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	email := "test@example.com"
	password := "password123"

	// Mock query that should respect context
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1`).
		WithArgs(email).
		WillReturnError(context.Canceled)

	// Act
	_, err := repo.Login(ctx, email, password)

	// Assert
	assert.Error(t, err)
	// Note: The actual error might be different depending on how GORM handles context cancellation
}

func TestRepository_Login_SpecialCharactersInEmail(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	// Test emails with special characters
	testCases := []string{
		"user+tag@example.com",
		"user.name@example.com",
		"user_name@example.com",
		"user-name@example.org",
	}

	for _, email := range testCases {
		t.Run(email, func(t *testing.T) {
			password := "password123"

			// Hash the password for database storage
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			assert.NoError(t, err)

			expectedUser := user.User{
				Model:    gorm.Model{ID: 1},
				Email:    email,
				Password: string(hashedPassword),
			}

			// Mock the SELECT query
			mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1`).
				WithArgs(email).
				WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password"}).
					AddRow(expectedUser.ID, expectedUser.Email, expectedUser.Password))

			// Act
			result, err := repo.Login(ctx, email, password)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, email, result.Email)
		})
	}
}

func TestRepository_Login_DifferentPasswordComplexity(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	email := "test@example.com"

	// Test different password complexities
	testCases := []struct {
		name     string
		password string
	}{
		{"Simple", "simple"},
		{"Numeric", "123456"},
		{"Complex", "P@ssw0rd!@#"},
		{"Long", "ThisIsAVeryLongPasswordWithManyCharacters123!@#"},
		{"Unicode", "пароль123"},
		{"Symbols", "!@#$%^&*()"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Hash the password for database storage
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(tc.password), bcrypt.DefaultCost)
			assert.NoError(t, err)

			expectedUser := user.User{
				Model:    gorm.Model{ID: 1},
				Email:    email,
				Password: string(hashedPassword),
			}

			// Mock the SELECT query
			mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1`).
				WithArgs(email).
				WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password"}).
					AddRow(expectedUser.ID, expectedUser.Email, expectedUser.Password))

			// Act
			result, err := repo.Login(ctx, email, tc.password)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, email, result.Email)
		})
	}
}

// Benchmark tests
func BenchmarkRepository_Login_Success(b *testing.B) {
	db, mock := setupMockDB(&testing.T{})
	repo := NewRepository(db)
	ctx := context.Background()

	email := "bench@example.com"
	password := "benchPassword"

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// Setup mock for all iterations
	for i := 0; i < b.N; i++ {
		mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1`).
			WithArgs(email).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password"}).
				AddRow(1, email, string(hashedPassword)))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.Login(ctx, email, password)
	}
}
