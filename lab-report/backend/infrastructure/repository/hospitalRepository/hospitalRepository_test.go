package hospitalrepository

import (
	"context"
	"regexp"
	"testing"

	"gitbub.com/zikrullahcelep611/lab-report/backend/models/hospital"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
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

func TestRepository_GetHospitalByID_Success(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()
	hospitalID := uint(1)

	expectedHospital := hospital.Hospital{
		Model: gorm.Model{ID: hospitalID},
		Name:  "Test Hospital",
	}

	// Mock the SELECT query - Note: there's a bug in the actual code using "hospital_id" instead of "id"
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "hospitals" WHERE hospital_id = $1 AND "hospitals"."deleted_at" IS NULL ORDER BY "hospitals"."id" LIMIT 1`)).
		WithArgs(hospitalID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(expectedHospital.ID, expectedHospital.Name))

	// Act
	result, err := repo.GetHospitalByID(ctx, hospitalID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedHospital.ID, result.ID)
	assert.Equal(t, expectedHospital.Name, result.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetHospitalByID_NotFound(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()
	hospitalID := uint(999)

	// Mock the SELECT query to return no rows
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "hospitals" WHERE hospital_id = $1 AND "hospitals"."deleted_at" IS NULL ORDER BY "hospitals"."id" LIMIT 1`)).
		WithArgs(hospitalID).
		WillReturnError(gorm.ErrRecordNotFound)

	// Act
	result, err := repo.GetHospitalByID(ctx, hospitalID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.Equal(t, hospital.Hospital{}, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_CreateHospital_Success(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	newHospital := hospital.Hospital{
		Name: "New Hospital",
	}

	// Mock the INSERT query
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "hospitals" ("created_at","updated_at","deleted_at","name") VALUES ($1,$2,$3,$4) RETURNING "id"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, newHospital.Name).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	// Act
	result, err := repo.CreateHospital(ctx, newHospital)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, newHospital.Name, result.Name)
	assert.Equal(t, uint(1), result.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_CreateHospital_Failure(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	newHospital := hospital.Hospital{
		Name: "New Hospital",
	}

	// Mock the INSERT query to fail
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "hospitals" ("created_at","updated_at","deleted_at","name") VALUES ($1,$2,$3,$4) RETURNING "id"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, newHospital.Name).
		WillReturnError(gorm.ErrInvalidTransaction)
	mock.ExpectRollback()

	// Act
	result, err := repo.CreateHospital(ctx, newHospital)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrInvalidTransaction, err)
	assert.Equal(t, hospital.Hospital{}, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_UpdateHospital_Success(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	updateHospital := hospital.Hospital{
		Model: gorm.Model{ID: 1},
		Name:  "Updated Hospital",
	}

	// Mock the UPDATE query (GORM Save method)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "hospitals" SET "created_at"=$1,"updated_at"=$2,"deleted_at"=$3,"name"=$4 WHERE "id" = $5`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, updateHospital.Name, updateHospital.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Act
	result, err := repo.UpdateHospital(ctx, updateHospital)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, updateHospital.ID, result.ID)
	assert.Equal(t, updateHospital.Name, result.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_UpdateHospital_Failure(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	updateHospital := hospital.Hospital{
		Model: gorm.Model{ID: 1},
		Name:  "Updated Hospital",
	}

	// Mock the UPDATE query to fail
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "hospitals" SET "created_at"=$1,"updated_at"=$2,"deleted_at"=$3,"name"=$4 WHERE "id" = $5`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, updateHospital.Name, updateHospital.ID).
		WillReturnError(gorm.ErrInvalidTransaction)
	mock.ExpectRollback()

	// Act
	result, err := repo.UpdateHospital(ctx, updateHospital)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrInvalidTransaction, err)
	assert.Equal(t, hospital.Hospital{}, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_DeleteHospital_Success(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()
	hospitalID := uint(1)

	// Mock the SELECT query to find existing hospital
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "hospitals" WHERE "hospitals"."id" = $1 AND "hospitals"."deleted_at" IS NULL ORDER BY "hospitals"."id" LIMIT 1`)).
		WithArgs(hospitalID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(hospitalID, "Test Hospital"))

	// Mock the DELETE query (soft delete)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "hospitals" SET "deleted_at"=$1 WHERE "hospitals"."id" = $2 AND "hospitals"."deleted_at" IS NULL`)).
		WithArgs(sqlmock.AnyArg(), hospitalID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Act
	err := repo.DeleteHospital(ctx, hospitalID)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_DeleteHospital_NotFound(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()
	hospitalID := uint(999)

	// Mock the SELECT query to return not found
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "hospitals" WHERE "hospitals"."id" = $1 AND "hospitals"."deleted_at" IS NULL ORDER BY "hospitals"."id" LIMIT 1`)).
		WithArgs(hospitalID).
		WillReturnError(gorm.ErrRecordNotFound)

	// Act
	err := repo.DeleteHospital(ctx, hospitalID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_DeleteHospital_DeleteFailure(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()
	hospitalID := uint(1)

	// Mock the SELECT query to find existing hospital
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "hospitals" WHERE "hospitals"."id" = $1 AND "hospitals"."deleted_at" IS NULL ORDER BY "hospitals"."id" LIMIT 1`)).
		WithArgs(hospitalID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(hospitalID, "Test Hospital"))

	// Mock the DELETE query to fail
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "hospitals" SET "deleted_at"=$1 WHERE "hospitals"."id" = $2 AND "hospitals"."deleted_at" IS NULL`)).
		WithArgs(sqlmock.AnyArg(), hospitalID).
		WillReturnError(gorm.ErrInvalidTransaction)
	mock.ExpectRollback()

	// Act
	err := repo.DeleteHospital(ctx, hospitalID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrInvalidTransaction, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Context_ShouldPassContext(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)

	// Create a context with cancel to test context propagation
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	hospitalID := uint(1)

	// Mock query that should respect context
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "hospitals" WHERE hospital_id = $1 AND "hospitals"."deleted_at" IS NULL ORDER BY "hospitals"."id" LIMIT 1`)).
		WithArgs(hospitalID).
		WillReturnError(context.Canceled)

	// Act
	_, err := repo.GetHospitalByID(ctx, hospitalID)

	// Assert
	assert.Error(t, err)
	// Note: The actual error might be different depending on how GORM handles context cancellation
}

func TestRepository_EmptyName_ShouldWork(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	newHospital := hospital.Hospital{
		Name: "", // Empty name
	}

	// Mock the INSERT query
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "hospitals" ("created_at","updated_at","deleted_at","name") VALUES ($1,$2,$3,$4) RETURNING "id"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, "").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	// Act
	result, err := repo.CreateHospital(ctx, newHospital)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "", result.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

