package redisrepository

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
)

// Test data structures
type TestData struct {
	Name  string
	Age   int
	Email string
}

type ComplexTestData struct {
	ID       int
	UserData TestData
	Tags     []string
	Metadata map[string]interface{}
}

func setupMockRedis() (*redis.Client, redismock.ClientMock) {
	db, mock := redismock.NewClientMock()
	return db, mock
}

func TestNewRepository_ShouldCreateRepository(t *testing.T) {
	// Arrange
	rdb, _ := setupMockRedis()

	// Act
	repo := NewRepository(rdb)

	// Assert
	assert.NotNil(t, repo)
	assert.Equal(t, rdb, repo.rdb)
}

func TestRepository_SetData_Success(t *testing.T) {
	// Arrange
	rdb, mock := setupMockRedis()
	repo := NewRepository(rdb)
	ctx := context.Background()

	testData := TestData{
		Name:  "John Doe",
		Age:   30,
		Email: "john@example.com",
	}

	// Mock Redis SET operation
	mock.ExpectSet("", nil, 5*time.Minute).SetVal("OK")

	// Act
	cacheKey, err := repo.SetData(ctx, testData)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, cacheKey)
	assert.Len(t, cacheKey, 36) // UUID length
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_SetData_NilData_ShouldReturnError(t *testing.T) {
	// Arrange
	rdb, mock := setupMockRedis()
	repo := NewRepository(rdb)
	ctx := context.Background()

	// Act
	cacheKey, err := repo.SetData(ctx, nil)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "data is nil", err.Error())
	assert.Empty(t, cacheKey)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_SetData_RedisError_ShouldReturnError(t *testing.T) {
	// Arrange
	rdb, mock := setupMockRedis()
	repo := NewRepository(rdb)
	ctx := context.Background()

	testData := TestData{
		Name:  "John Doe",
		Age:   30,
		Email: "john@example.com",
	}

	// Mock Redis SET operation to return error
	mock.ExpectSet("", nil, 5*time.Minute).SetErr(redis.TxFailedErr)

	// Act
	cacheKey, err := repo.SetData(ctx, testData)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "redis set errors", err.Error())
	assert.Empty(t, cacheKey)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetData_KeyNotFound_ShouldReturnError(t *testing.T) {
	// Arrange
	rdb, mock := setupMockRedis()
	repo := NewRepository(rdb)
	ctx := context.Background()
	cacheKey := "nonexistent-key"

	// Mock Redis GET operation to return redis.Nil
	mock.ExpectGet(cacheKey).RedisNil()

	// Act
	var retrievedData TestData
	err := repo.GetData(ctx, cacheKey, &retrievedData)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "redis get error: key does not exist", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetData_RedisError_ShouldReturnError(t *testing.T) {
	// Arrange
	rdb, mock := setupMockRedis()
	repo := NewRepository(rdb)
	ctx := context.Background()
	cacheKey := "test-key"

	// Mock Redis GET operation to return error
	mock.ExpectGet(cacheKey).SetErr(redis.TxFailedErr)

	// Act
	var retrievedData TestData
	err := repo.GetData(ctx, cacheKey, &retrievedData)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "redis get errors", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_DeleteData_Success(t *testing.T) {
	// Arrange
	rdb, mock := setupMockRedis()
	repo := NewRepository(rdb)
	ctx := context.Background()
	cacheKey := "test-cache-key"

	// Mock Redis DEL operation
	mock.ExpectDel(cacheKey).SetVal(1) // 1 key deleted

	// Act
	err := repo.DeleteData(ctx, cacheKey)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_DeleteData_RedisError_ShouldReturnError(t *testing.T) {
	// Arrange
	rdb, mock := setupMockRedis()
	repo := NewRepository(rdb)
	ctx := context.Background()
	cacheKey := "test-key"

	// Mock Redis DEL operation to return error
	mock.ExpectDel(cacheKey).SetErr(redis.TxFailedErr)

	// Act
	err := repo.DeleteData(ctx, cacheKey)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "redis delete errors", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_DeleteData_KeyNotExists_ShouldSucceed(t *testing.T) {
	// Arrange
	rdb, mock := setupMockRedis()
	repo := NewRepository(rdb)
	ctx := context.Background()
	cacheKey := "nonexistent-key"

	// Mock Redis DEL operation - returns 0 when key doesn't exist
	mock.ExpectDel(cacheKey).SetVal(0) // 0 keys deleted

	// Act
	err := repo.DeleteData(ctx, cacheKey)

	// Assert
	assert.NoError(t, err) // Should still succeed even if key doesn't exist
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_SetData_DifferentDataTypes(t *testing.T) {
	// Arrange
	rdb, mock := setupMockRedis()
	repo := NewRepository(rdb)
	ctx := context.Background()

	testCases := []struct {
		name string
		data interface{}
	}{
		{"String", "test string"},
		{"Integer", 42},
		{"Float", 3.14159},
		{"Boolean", true},
		{"Slice", []string{"a", "b", "c"}},
		{"Map", map[string]int{"one": 1, "two": 2}},
		{"Struct", TestData{Name: "Test", Age: 30, Email: "test@example.com"}},
		{"Complex Struct", ComplexTestData{
			ID:       1,
			UserData: TestData{Name: "Complex", Age: 25, Email: "complex@example.com"},
			Tags:     []string{"tag1", "tag2"},
			Metadata: map[string]interface{}{"key": "value"},
		}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Mock Redis SET operation
			mock.ExpectSet("", nil, 5*time.Minute).SetVal("OK")

			// Act
			cacheKey, err := repo.SetData(ctx, tc.data)

			// Assert
			assert.NoError(t, err)
			assert.NotEmpty(t, cacheKey)
		})
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Context_ShouldPassContext(t *testing.T) {
	// Arrange
	rdb, mock := setupMockRedis()
	repo := NewRepository(rdb)

	// Create a context with cancel to test context propagation
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	testData := TestData{Name: "Test", Age: 30, Email: "test@example.com"}

	// Mock Redis operation that should respect context
	mock.ExpectSet("", nil, 5*time.Minute).SetErr(context.Canceled)

	// Act
	_, err := repo.SetData(ctx, testData)

	// Assert
	assert.Error(t, err)
	// Note: The actual error might be different depending on how Redis handles context cancellation
}

func TestRepository_EmptyStringData_ShouldWork(t *testing.T) {
	// Arrange
	rdb, mock := setupMockRedis()
	repo := NewRepository(rdb)
	ctx := context.Background()

	emptyString := ""

	// Mock Redis SET operation
	mock.ExpectSet("", nil, 5*time.Minute).SetVal("OK")

	// Act
	cacheKey, err := repo.SetData(ctx, emptyString)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, cacheKey)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_ZeroValues_ShouldWork(t *testing.T) {
	// Arrange
	rdb, mock := setupMockRedis()
	repo := NewRepository(rdb)
	ctx := context.Background()

	testCases := []struct {
		name string
		data interface{}
	}{
		{"Zero Int", 0},
		{"Zero Float", 0.0},
		{"False Boolean", false},
		{"Empty Slice", []string{}},
		{"Empty Map", map[string]int{}},
		{"Empty Struct", TestData{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Mock Redis SET operation
			mock.ExpectSet("", nil, 5*time.Minute).SetVal("OK")

			// Act
			cacheKey, err := repo.SetData(ctx, tc.data)

			// Assert
			assert.NoError(t, err)
			assert.NotEmpty(t, cacheKey)
		})
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

// Benchmark tests
func BenchmarkRepository_SetData(b *testing.B) {
	rdb, mock := setupMockRedis()
	repo := NewRepository(rdb)
	ctx := context.Background()

	testData := TestData{
		Name:  "Benchmark User",
		Age:   30,
		Email: "benchmark@example.com",
	}

	// Setup mock for all iterations
	for i := 0; i < b.N; i++ {
		mock.ExpectSet("", nil, 5*time.Minute).SetVal("OK")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.SetData(ctx, testData)
	}
}

func BenchmarkRepository_DeleteData(b *testing.B) {
	rdb, mock := setupMockRedis()
	repo := NewRepository(rdb)
	ctx := context.Background()
	cacheKey := "benchmark-key"

	// Setup mock for all iterations
	for i := 0; i < b.N; i++ {
		mock.ExpectDel(cacheKey).SetVal(1)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = repo.DeleteData(ctx, cacheKey)
	}
}
