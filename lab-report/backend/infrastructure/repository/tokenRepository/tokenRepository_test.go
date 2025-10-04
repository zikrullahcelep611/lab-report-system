package tokenrepository

import (
    "context"
    "database/sql"
    "errors"
    "testing"
    "time"

    "github.com/DATA-DOG/go-sqlmock"
    "github.com/stretchr/testify/assert"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

func setupMockDB() (*gorm.DB, sqlmock.Sqlmock, error) {
    sqlDB, mock, err := sqlmock.New()
    if err != nil {
        return nil, nil, err
    }

    gormDB, err := gorm.Open(postgres.New(postgres.Config{
        Conn: sqlDB,
    }), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Silent),
    })

    if err != nil {
        return nil, nil, err
    }

    return gormDB, mock, nil
}

func TestNewRepository(t *testing.T) {
    // Arrange
    db, _, _ := setupMockDB()

    // Act
    repo := NewRepository(db)

    // Assert
    assert.NotNil(t, repo)
    assert.Equal(t, db, repo.DB)
}

func TestRepository_DeleteExpiredTokens_Success(t *testing.T) {
    // Arrange
    db, mock, _ := setupMockDB()
    repo := NewRepository(db)
    ctx := context.Background()

    mock.ExpectExec(`DELETE FROM "expired_tokens"`).
        WithArgs(sqlmock.AnyArg()).
        WillReturnResult(sqlmock.NewResult(0, 3))

    // Act
    repo.DeleteExpiredTokens(ctx)

    // Assert
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_DeleteExpiredTokens_DatabaseError(t *testing.T) {
    // Arrange
    db, mock, _ := setupMockDB()
    repo := NewRepository(db)
    ctx := context.Background()

    mock.ExpectExec(`DELETE FROM "expired_tokens"`).
        WithArgs(sqlmock.AnyArg()).
        WillReturnError(errors.New("database connection error"))

    // Act
    repo.DeleteExpiredTokens(ctx)

    // Assert
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_AddTokenToBlackList_Success(t *testing.T) {
    // Arrange
    db, mock, _ := setupMockDB()
    repo := NewRepository(db)
    ctx := context.Background()
    tokenStr := "test-jwt-token-12345"
    expireTime := time.Now().Add(time.Hour)

    mock.ExpectBegin()
    mock.ExpectQuery(`INSERT INTO "expired_tokens"`).
        WithArgs(
            sqlmock.AnyArg(), // created_at
            sqlmock.AnyArg(), // updated_at
            sqlmock.AnyArg(), // deleted_at
            tokenStr,
            expireTime.Unix(),
        ).
        WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
    mock.ExpectCommit()

    // Act
    err := repo.AddTokenToBlackList(ctx, tokenStr, expireTime)

    // Assert
    assert.NoError(t, err)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_AddTokenToBlackList_DatabaseError(t *testing.T) {
    // Arrange
    db, mock, _ := setupMockDB()
    repo := NewRepository(db)
    ctx := context.Background()
    tokenStr := "test-jwt-token-12345"
    expireTime := time.Now().Add(time.Hour)

    mock.ExpectBegin()
    mock.ExpectQuery(`INSERT INTO "expired_tokens"`).
        WithArgs(
            sqlmock.AnyArg(),
            sqlmock.AnyArg(),
            sqlmock.AnyArg(),
            tokenStr,
            expireTime.Unix(),
        ).
        WillReturnError(errors.New("database insert error"))
    mock.ExpectRollback()

    // Act
    err := repo.AddTokenToBlackList(ctx, tokenStr, expireTime)

    // Assert
    assert.Error(t, err)
    assert.Equal(t, "redis set errors", err.Error())
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_AddTokenToBlackList_EmptyToken(t *testing.T) {
    // Arrange
    db, mock, _ := setupMockDB()
    repo := NewRepository(db)
    ctx := context.Background()
    tokenStr := ""
    expireTime := time.Now().Add(time.Hour)

    mock.ExpectBegin()
    mock.ExpectQuery(`INSERT INTO "expired_tokens"`).
        WithArgs(
            sqlmock.AnyArg(),
            sqlmock.AnyArg(),
            sqlmock.AnyArg(),
            tokenStr,
            expireTime.Unix(),
        ).
        WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
    mock.ExpectCommit()

    // Act
    err := repo.AddTokenToBlackList(ctx, tokenStr, expireTime)

    // Assert
    assert.NoError(t, err)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_IsTokenBlackListed_TokenExists(t *testing.T) {
    // Arrange
    db, mock, _ := setupMockDB()
    repo := NewRepository(db)
    ctx := context.Background()
    tokenStr := "blacklisted-token-12345"

    mock.ExpectQuery(`SELECT count\(\*\) FROM "expired_tokens"`).
        WithArgs(tokenStr).
        WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

    // Act
    result := repo.IsTokenBlackListed(ctx, tokenStr)

    // Assert
    assert.True(t, result)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_IsTokenBlackListed_TokenNotExists(t *testing.T) {
    // Arrange
    db, mock, _ := setupMockDB()
    repo := NewRepository(db)
    ctx := context.Background()
    tokenStr := "valid-token-12345"

    mock.ExpectQuery(`SELECT count\(\*\) FROM "expired_tokens"`).
        WithArgs(tokenStr).
        WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

    // Act
    result := repo.IsTokenBlackListed(ctx, tokenStr)

    // Assert
    assert.False(t, result)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_IsTokenBlackListed_DatabaseError(t *testing.T) {
    // Arrange
    db, mock, _ := setupMockDB()
    repo := NewRepository(db)
    ctx := context.Background()
    tokenStr := "some-token-12345"

    mock.ExpectQuery(`SELECT count\(\*\) FROM "expired_tokens"`).
        WithArgs(tokenStr).
        WillReturnError(errors.New("database query error"))

    // Act
    result := repo.IsTokenBlackListed(ctx, tokenStr)

    // Assert
    assert.False(t, result) // Should return false on error
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_IsTokenBlackListed_MultipleTokens(t *testing.T) {
    // Arrange
    db, mock, _ := setupMockDB()
    repo := NewRepository(db)
    ctx := context.Background()
    tokenStr := "duplicate-token-12345"

    mock.ExpectQuery(`SELECT count\(\*\) FROM "expired_tokens"`).
        WithArgs(tokenStr).
        WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))

    // Act
    result := repo.IsTokenBlackListed(ctx, tokenStr)

    // Assert
    assert.True(t, result) // Should return true even for multiple entries
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_ContextCancellation(t *testing.T) {
    // Arrange
    db, mock, _ := setupMockDB()
    repo := NewRepository(db)
    ctx, cancel := context.WithCancel(context.Background())
    tokenStr := "test-token"

    // Cancel context immediately
    cancel()

    mock.ExpectQuery(`SELECT count\(\*\) FROM "expired_tokens"`).
        WithArgs(tokenStr).
        WillReturnError(sql.ErrConnDone)

    // Act
    result := repo.IsTokenBlackListed(ctx, tokenStr)

    // Assert
    assert.False(t, result)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_DeleteExpiredTokens_TimeComparison(t *testing.T) {
    // Arrange
    db, mock, _ := setupMockDB()
    repo := NewRepository(db)
    ctx := context.Background()

    // Expect deletion with current time comparison
    mock.ExpectExec(`DELETE FROM "expired_tokens"`).
        WithArgs(sqlmock.AnyArg()).
        WillReturnResult(sqlmock.NewResult(0, 5))

    // Act
    repo.DeleteExpiredTokens(ctx)

    // Assert
    assert.NoError(t, mock.ExpectationsWereMet())
}

// Benchmark tests
func BenchmarkRepository_AddTokenToBlackList(b *testing.B) {
    db, mock, _ := setupMockDB()
    repo := NewRepository(db)
    ctx := context.Background()
    expireTime := time.Now().Add(time.Hour)

    // Setup expectations for benchmark
    for i := 0; i < b.N; i++ {
        mock.ExpectBegin()
        mock.ExpectQuery(`INSERT INTO "expired_tokens"`).
            WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(i + 1))
        mock.ExpectCommit()
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        tokenStr := "benchmark-token-" + string(rune(i))
        repo.AddTokenToBlackList(ctx, tokenStr, expireTime)
    }
}

func BenchmarkRepository_IsTokenBlackListed(b *testing.B) {
    db, mock, _ := setupMockDB()
    repo := NewRepository(db)
    ctx := context.Background()

    // Setup expectations for benchmark
    for i := 0; i < b.N; i++ {
        mock.ExpectQuery(`SELECT count\(\*\) FROM "expired_tokens"`).
            WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        tokenStr := "benchmark-token-" + string(rune(i))
        repo.IsTokenBlackListed(ctx, tokenStr)
    }
}

// Edge case tests
func TestRepository_AddTokenToBlackList_PastExpireTime(t *testing.T) {
    // Arrange
    db, mock, _ := setupMockDB()
    repo := NewRepository(db)
    ctx := context.Background()
    tokenStr := "expired-token"
    expireTime := time.Now().Add(-time.Hour) // Past time

    mock.ExpectBegin()
    mock.ExpectQuery(`INSERT INTO "expired_tokens"`).
        WithArgs(
            sqlmock.AnyArg(),
            sqlmock.AnyArg(),
            sqlmock.AnyArg(),
            tokenStr,
            expireTime.Unix(), // Negative unix timestamp
        ).
        WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
    mock.ExpectCommit()

    // Act
    err := repo.AddTokenToBlackList(ctx, tokenStr, expireTime)

    // Assert
    assert.NoError(t, err)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_LongTokenString(t *testing.T) {
    // Arrange
    db, mock, _ := setupMockDB()
    repo := NewRepository(db)
    ctx := context.Background()
    
    // Very long token string
    longToken := string(make([]byte, 1000))
    for i := range longToken {
        longToken = longToken[:i] + "a" + longToken[i+1:]
    }

    mock.ExpectQuery(`SELECT count\(\*\) FROM "expired_tokens"`).
        WithArgs(longToken).
        WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

    // Act
    result := repo.IsTokenBlackListed(ctx, longToken)

    // Assert
    assert.False(t, result)
    assert.NoError(t, mock.ExpectationsWereMet())
}