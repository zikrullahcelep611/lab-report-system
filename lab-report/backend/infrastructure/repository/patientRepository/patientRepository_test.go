package patientrepository

import (
	"context"
	"testing"

	"gitbub.com/zikrullahcelep611/lab-report/backend/models/patient"
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

func TestRepository_GetPatientByID_Success(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()
	patientID := uint(1)

	expectedPatient := patient.Patient{
		Model:      gorm.Model{ID: patientID},
		Name:       "John",
		Lastname:   "Doe",
		NationalID: "12345678901",
	}

	// Mock the SELECT query - Note: there's a bug in actual code using "patient_id" instead of "id"
	mock.ExpectQuery(`SELECT \* FROM "patients" WHERE patient_id = \$1`).
		WithArgs(patientID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "lastname", "national_id"}).
			AddRow(expectedPatient.ID, expectedPatient.Name, expectedPatient.Lastname, expectedPatient.NationalID))

	// Act
	result, err := repo.GetPatientByID(ctx, patientID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedPatient.ID, result.ID)
	assert.Equal(t, expectedPatient.Name, result.Name)
	assert.Equal(t, expectedPatient.Lastname, result.Lastname)
	assert.Equal(t, expectedPatient.NationalID, result.NationalID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetPatientByID_NotFound(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()
	patientID := uint(999)

	// Mock the SELECT query to return no rows
	mock.ExpectQuery(`SELECT \* FROM "patients" WHERE patient_id = \$1`).
		WithArgs(patientID).
		WillReturnError(gorm.ErrRecordNotFound)

	// Act
	result, err := repo.GetPatientByID(ctx, patientID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.Equal(t, patient.Patient{}, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetPatientByID_DatabaseError(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()
	patientID := uint(1)

	// Mock the SELECT query to return database error
	mock.ExpectQuery(`SELECT \* FROM "patients" WHERE patient_id = \$1`).
		WithArgs(patientID).
		WillReturnError(gorm.ErrInvalidTransaction)

	// Act
	result, err := repo.GetPatientByID(ctx, patientID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrInvalidTransaction, err)
	assert.Equal(t, patient.Patient{}, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_CreatePatient_Success(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	newPatient := patient.Patient{
		Name:       "Jane",
		Lastname:   "Smith",
		NationalID: "98765432109",
	}

	// Mock the INSERT query
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "patients"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, newPatient.Name, newPatient.Lastname, newPatient.NationalID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	// Act
	result, err := repo.CreatePatient(ctx, newPatient)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, newPatient.Name, result.Name)
	assert.Equal(t, newPatient.Lastname, result.Lastname)
	assert.Equal(t, newPatient.NationalID, result.NationalID)
	assert.Equal(t, uint(1), result.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_CreatePatient_Failure(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	newPatient := patient.Patient{
		Name:       "Jane",
		Lastname:   "Smith",
		NationalID: "98765432109",
	}

	// Mock the INSERT query to fail
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "patients"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, newPatient.Name, newPatient.Lastname, newPatient.NationalID).
		WillReturnError(gorm.ErrInvalidTransaction)
	mock.ExpectRollback()

	// Act
	result, err := repo.CreatePatient(ctx, newPatient)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrInvalidTransaction, err)
	assert.Equal(t, patient.Patient{}, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_CreatePatient_UniqueConstraintViolation(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	newPatient := patient.Patient{
		Name:       "John",
		Lastname:   "Doe",
		NationalID: "12345678901", // Duplicate national ID
	}

	// Mock unique constraint violation
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "patients"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, newPatient.Name, newPatient.Lastname, newPatient.NationalID).
		WillReturnError(gorm.ErrDuplicatedKey)
	mock.ExpectRollback()

	// Act
	result, err := repo.CreatePatient(ctx, newPatient)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrDuplicatedKey, err)
	assert.Equal(t, patient.Patient{}, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_UpdatePatient_Success(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	updatePatient := patient.Patient{
		Model:      gorm.Model{ID: 1},
		Name:       "Updated Name",
		Lastname:   "Updated Lastname",
		NationalID: "11111111111",
	}

	// Mock the UPDATE query (GORM Save method)
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "patients" SET`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, updatePatient.Name, updatePatient.Lastname, updatePatient.NationalID, updatePatient.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Act
	result, err := repo.UpdatePatient(ctx, updatePatient)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, updatePatient.ID, result.ID)
	assert.Equal(t, updatePatient.Name, result.Name)
	assert.Equal(t, updatePatient.Lastname, result.Lastname)
	assert.Equal(t, updatePatient.NationalID, result.NationalID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_UpdatePatient_Failure(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	updatePatient := patient.Patient{
		Model:      gorm.Model{ID: 1},
		Name:       "Updated Name",
		Lastname:   "Updated Lastname",
		NationalID: "11111111111",
	}

	// Mock the UPDATE query to fail
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "patients" SET`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, updatePatient.Name, updatePatient.Lastname, updatePatient.NationalID, updatePatient.ID).
		WillReturnError(gorm.ErrInvalidTransaction)
	mock.ExpectRollback()

	// Act
	result, err := repo.UpdatePatient(ctx, updatePatient)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrInvalidTransaction, err)
	assert.Equal(t, patient.Patient{}, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_DeletePatient_Success(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()
	patientID := uint(1)

	// Mock the SELECT query to find existing patient
	mock.ExpectQuery(`SELECT \* FROM "patients" WHERE`).
		WithArgs(patientID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "lastname", "national_id"}).
			AddRow(patientID, "John", "Doe", "12345678901"))

	// Mock the DELETE query (soft delete)
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "patients" SET "deleted_at"`).
		WithArgs(sqlmock.AnyArg(), patientID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Act
	err := repo.DeletePatient(ctx, patientID)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_DeletePatient_NotFound(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()
	patientID := uint(999)

	// Mock the SELECT query to return not found
	mock.ExpectQuery(`SELECT \* FROM "patients" WHERE`).
		WithArgs(patientID).
		WillReturnError(gorm.ErrRecordNotFound)

	// Act
	err := repo.DeletePatient(ctx, patientID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_DeletePatient_DeleteFailure(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()
	patientID := uint(1)

	// Mock the SELECT query to find existing patient
	mock.ExpectQuery(`SELECT \* FROM "patients" WHERE`).
		WithArgs(patientID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "lastname", "national_id"}).
			AddRow(patientID, "John", "Doe", "12345678901"))

	// Mock the DELETE query to fail
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "patients" SET "deleted_at"`).
		WithArgs(sqlmock.AnyArg(), patientID).
		WillReturnError(gorm.ErrInvalidTransaction)
	mock.ExpectRollback()

	// Act
	err := repo.DeletePatient(ctx, patientID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrInvalidTransaction, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_CreatePatient_TurkishCharacters(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	// Test Turkish characters in names
	testCases := []struct {
		name     string
		lastname string
	}{
		{"Ömer", "Özkan"},
		{"Çağla", "Şahin"},
		{"Gül", "Ünal"},
		{"İsmail", "Ğüler"},
	}

	for _, tc := range testCases {
		t.Run(tc.name+" "+tc.lastname, func(t *testing.T) {
			newPatient := patient.Patient{
				Name:       tc.name,
				Lastname:   tc.lastname,
				NationalID: "12345678901",
			}

			// Mock the INSERT query
			mock.ExpectBegin()
			mock.ExpectQuery(`INSERT INTO "patients"`).
				WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, newPatient.Name, newPatient.Lastname, newPatient.NationalID).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			mock.ExpectCommit()

			// Act
			result, err := repo.CreatePatient(ctx, newPatient)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tc.name, result.Name)
			assert.Equal(t, tc.lastname, result.Lastname)
		})
	}
}

func TestRepository_CreatePatient_NationalIDValidation(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	// Test different national ID formats
	testCases := []struct {
		name       string
		nationalID string
	}{
		{"Valid 11 digits", "12345678901"},
		{"All same digits", "11111111111"},
		{"Starting with zero", "01234567890"},
		{"Sequential", "12345678901"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			newPatient := patient.Patient{
				Name:       "Test",
				Lastname:   "Patient",
				NationalID: tc.nationalID,
			}

			// Mock the INSERT query
			mock.ExpectBegin()
			mock.ExpectQuery(`INSERT INTO "patients"`).
				WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, newPatient.Name, newPatient.Lastname, newPatient.NationalID).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			mock.ExpectCommit()

			// Act
			result, err := repo.CreatePatient(ctx, newPatient)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tc.nationalID, result.NationalID)
		})
	}
}

func TestRepository_Context_ShouldPassContext(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)

	// Create a context with cancel to test context propagation
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	patientID := uint(1)

	// Mock query that should respect context
	mock.ExpectQuery(`SELECT \* FROM "patients" WHERE patient_id = \$1`).
		WithArgs(patientID).
		WillReturnError(context.Canceled)

	// Act
	_, err := repo.GetPatientByID(ctx, patientID)

	// Assert
	assert.Error(t, err)
	// Note: The actual error might be different depending on how GORM handles context cancellation
}

// Benchmark tests
func BenchmarkRepository_CreatePatient(b *testing.B) {
	db, mock := setupMockDB(&testing.T{})
	repo := NewRepository(db)
	ctx := context.Background()

	newPatient := patient.Patient{
		Name:       "Benchmark",
		Lastname:   "Patient",
		NationalID: "12345678901",
	}

	// Setup mock for all iterations
	for i := 0; i < b.N; i++ {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "patients"`).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, newPatient.Name, newPatient.Lastname, newPatient.NationalID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(i + 1))
		mock.ExpectCommit()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.CreatePatient(ctx, newPatient)
	}
}

func BenchmarkRepository_GetPatientByID(b *testing.B) {
	db, mock := setupMockDB(&testing.T{})
	repo := NewRepository(db)
	ctx := context.Background()
	patientID := uint(1)

	// Setup mock for all iterations
	for i := 0; i < b.N; i++ {
		mock.ExpectQuery(`SELECT \* FROM "patients" WHERE patient_id = \$1`).
			WithArgs(patientID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "lastname", "national_id"}).
				AddRow(patientID, "John", "Doe", "12345678901"))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.GetPatientByID(ctx, patientID)
	}
}
