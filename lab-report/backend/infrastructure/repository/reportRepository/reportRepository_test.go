package reportrepository

import (
	"context"
	"testing"
	"time"

	"gitbub.com/zikrullahcelep611/lab-report/backend/models/report"
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

func TestRepository_CreateReport_Success(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	newReport := report.Report{
		DiagnosisTitle:   "Blood Test Results",
		DiagnosisDetails: "All parameters within normal range",
		ImagePath:        "/uploads/blood_test.pdf",
		UserID:           1,
		PatientID:        2,
	}

	// Mock the INSERT query
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "reports"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil,
			newReport.DiagnosisTitle, newReport.DiagnosisDetails, newReport.ImagePath,
			newReport.UserID, newReport.PatientID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	// Mock the Preload query for relationships
	mock.ExpectQuery(`SELECT \* FROM "reports" WHERE`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "diagnosis_title", "diagnosis_details", "image_path", "user_id", "patient_id"}).
			AddRow(1, newReport.DiagnosisTitle, newReport.DiagnosisDetails, newReport.ImagePath, newReport.UserID, newReport.PatientID))

	// Mock User preload
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "firstname", "lastname", "email"}).
			AddRow(1, "Dr. John", "Doe", "john@hospital.com"))

	// Mock Patient preload
	mock.ExpectQuery(`SELECT \* FROM "patients" WHERE`).
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "lastname", "national_id"}).
			AddRow(2, "Jane", "Smith", "12345678901"))

	// Act
	result, err := repo.CreateReport(ctx, newReport)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, newReport.DiagnosisTitle, result.DiagnosisTitle)
	assert.Equal(t, newReport.DiagnosisDetails, result.DiagnosisDetails)
	assert.Equal(t, uint(1), result.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_CreateReport_InsertFailure(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	newReport := report.Report{
		DiagnosisTitle:   "Test Report",
		DiagnosisDetails: "Test Details",
		UserID:           1,
		PatientID:        2,
	}

	// Mock the INSERT query to fail
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "reports"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil,
			newReport.DiagnosisTitle, newReport.DiagnosisDetails, newReport.ImagePath,
			newReport.UserID, newReport.PatientID).
		WillReturnError(gorm.ErrInvalidTransaction)
	mock.ExpectRollback()

	// Act
	result, err := repo.CreateReport(ctx, newReport)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrInvalidTransaction, err)
	assert.Equal(t, report.Report{}, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetReportByID_Success(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()
	reportID := uint(1)

	expectedReport := report.Report{
		Model:            gorm.Model{ID: reportID},
		DiagnosisTitle:   "X-Ray Results",
		DiagnosisDetails: "No abnormalities detected",
		ImagePath:        "/uploads/xray.jpg",
		UserID:           1,
		PatientID:        2,
	}

	// Mock the SELECT query with preloads
	mock.ExpectQuery(`SELECT \* FROM "reports" WHERE id = \$1`).
		WithArgs(reportID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "diagnosis_title", "diagnosis_details", "image_path", "user_id", "patient_id"}).
			AddRow(expectedReport.ID, expectedReport.DiagnosisTitle, expectedReport.DiagnosisDetails,
				expectedReport.ImagePath, expectedReport.UserID, expectedReport.PatientID))

	// Mock Patient preload
	mock.ExpectQuery(`SELECT \* FROM "patients" WHERE`).
		WithArgs(expectedReport.PatientID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "lastname", "national_id"}).
			AddRow(2, "John", "Doe", "12345678901"))

	// Mock User preload
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE`).
		WithArgs(expectedReport.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "firstname", "lastname", "email"}).
			AddRow(1, "Dr. Jane", "Smith", "jane@hospital.com"))

	// Act
	result, err := repo.GetReportByID(ctx, reportID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedReport.ID, result.ID)
	assert.Equal(t, expectedReport.DiagnosisTitle, result.DiagnosisTitle)
	assert.Equal(t, expectedReport.DiagnosisDetails, result.DiagnosisDetails)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetReportByID_NotFound(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()
	reportID := uint(999)

	// Mock the SELECT query to return no rows
	mock.ExpectQuery(`SELECT \* FROM "reports" WHERE id = \$1`).
		WithArgs(reportID).
		WillReturnError(gorm.ErrRecordNotFound)

	// Act
	result, err := repo.GetReportByID(ctx, reportID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.Equal(t, report.Report{}, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetReportsWithPatientName_Success(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	firstName := "John"
	lastName := "Doe"

	// Mock the JOIN query
	mock.ExpectQuery(`SELECT "reports"\."id","reports"\."created_at","reports"\."updated_at","reports"\."deleted_at","reports"\."diagnosis_title","reports"\."diagnosis_details","reports"\."image_path","reports"\."user_id","reports"\."patient_id" FROM "reports" JOIN patients ON reports\.patient_id = patients\.id WHERE \(patients\.name = \$1 AND patients\.lastname = \$2\)`).
		WithArgs(firstName, lastName).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "diagnosis_title", "diagnosis_details", "image_path", "user_id", "patient_id"}).
			AddRow(1, time.Now(), time.Now(), nil, "Test Report", "Test Details", "/test.pdf", 1, 2))

	// Mock User preload
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "firstname", "lastname", "email"}).
			AddRow(1, "Dr. Jane", "Smith", "jane@hospital.com"))

	// Mock Patient preload
	mock.ExpectQuery(`SELECT \* FROM "patients" WHERE`).
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "lastname", "national_id"}).
			AddRow(2, firstName, lastName, "12345678901"))

	// Act
	result, err := repo.GetReportsWithPatientName(ctx, firstName, lastName)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Test Report", result[0].DiagnosisTitle)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetReportsWithPatientNationalID_Success(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	nationalID := "12345678901"

	// Mock the JOIN query
	mock.ExpectQuery(`SELECT "reports"\."id","reports"\."created_at","reports"\."updated_at","reports"\."deleted_at","reports"\."diagnosis_title","reports"\."diagnosis_details","reports"\."image_path","reports"\."user_id","reports"\."patient_id" FROM "reports" JOIN patients ON reports\.patient_id = patients\.id WHERE patients\.national_id = \$1`).
		WithArgs(nationalID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "diagnosis_title", "diagnosis_details", "image_path", "user_id", "patient_id"}).
			AddRow(1, time.Now(), time.Now(), nil, "Blood Test", "Normal results", "/blood.pdf", 1, 2))

	// Mock Patient preload
	mock.ExpectQuery(`SELECT \* FROM "patients" WHERE`).
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "lastname", "national_id"}).
			AddRow(2, "Ali", "Veli", nationalID))

	// Mock User preload
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "firstname", "lastname", "email"}).
			AddRow(1, "Dr. Mehmet", "Yılmaz", "mehmet@hospital.com"))

	// Act
	result, err := repo.GetReportsWithPatientNationalID(ctx, nationalID)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Blood Test", result[0].DiagnosisTitle)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetAllReportsOrdered_Success(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	// Mock the SELECT query with ORDER BY
	mock.ExpectQuery(`SELECT \* FROM "reports" WHERE "reports"\."deleted_at" IS NULL ORDER BY created_at DESC`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "diagnosis_title", "diagnosis_details", "image_path", "user_id", "patient_id"}).
			AddRow(2, time.Now().Add(-time.Hour), time.Now(), nil, "Recent Report", "Recent Details", "/recent.pdf", 1, 2).
			AddRow(1, time.Now().Add(-2*time.Hour), time.Now(), nil, "Older Report", "Older Details", "/older.pdf", 1, 2))

	// Mock User preloads
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE`).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "firstname", "lastname", "email"}).
			AddRow(1, "Dr. Test", "Doctor", "test@hospital.com"))

	// Mock Patient preloads
	mock.ExpectQuery(`SELECT \* FROM "patients" WHERE`).
		WithArgs(2, 2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "lastname", "national_id"}).
			AddRow(2, "Test", "Patient", "12345678901"))

	// Act
	result, err := repo.GetAllReportsOrdered(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Recent Report", result[0].DiagnosisTitle)
	assert.Equal(t, "Older Report", result[1].DiagnosisTitle)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_UpdateReport_Success(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	updateReport := report.Report{
		Model:            gorm.Model{ID: 1},
		DiagnosisTitle:   "Updated Report Title",
		DiagnosisDetails: "Updated Details",
		ImagePath:        "/updated.pdf",
		UserID:           1,
		PatientID:        2,
	}

	// Mock the UPDATE query (GORM Save method)
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "reports" SET`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil,
			updateReport.DiagnosisTitle, updateReport.DiagnosisDetails, updateReport.ImagePath,
			updateReport.UserID, updateReport.PatientID, updateReport.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Act
	result, err := repo.UpdateReport(ctx, updateReport)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, updateReport.ID, result.ID)
	assert.Equal(t, updateReport.DiagnosisTitle, result.DiagnosisTitle)
	assert.Equal(t, updateReport.DiagnosisDetails, result.DiagnosisDetails)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_DeleteReport_Success(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()
	reportID := uint(1)

	// Mock the DELETE query (soft delete)
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "reports" SET "deleted_at"`).
		WithArgs(sqlmock.AnyArg(), reportID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Act
	result, err := repo.DeleteReport(ctx, reportID)

	// Assert
	assert.NoError(t, err)
	assert.True(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_DeleteReport_Failure(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()
	reportID := uint(1)

	// Mock the DELETE query to fail
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "reports" SET "deleted_at"`).
		WithArgs(sqlmock.AnyArg(), reportID).
		WillReturnError(gorm.ErrInvalidTransaction)
	mock.ExpectRollback()

	// Act
	result, err := repo.DeleteReport(ctx, reportID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrInvalidTransaction, err)
	assert.False(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetReportsWithPatientName_EmptyResult(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	firstName := "NonExistent"
	lastName := "Patient"

	// Mock the JOIN query to return no rows
	mock.ExpectQuery(`SELECT "reports"\."id","reports"\."created_at","reports"\."updated_at","reports"\."deleted_at","reports"\."diagnosis_title","reports"\."diagnosis_details","reports"\."image_path","reports"\."user_id","reports"\."patient_id" FROM "reports" JOIN patients ON reports\.patient_id = patients\.id WHERE \(patients\.name = \$1 AND patients\.lastname = \$2\)`).
		WithArgs(firstName, lastName).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "diagnosis_title", "diagnosis_details", "image_path", "user_id", "patient_id"}))

	// Act
	result, err := repo.GetReportsWithPatientName(ctx, firstName, lastName)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 0)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Context_ShouldPassContext(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)

	// Create a context with cancel to test context propagation
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	reportID := uint(1)

	// Mock query that should respect context
	mock.ExpectQuery(`SELECT \* FROM "reports" WHERE id = \$1`).
		WithArgs(reportID).
		WillReturnError(context.Canceled)

	// Act
	_, err := repo.GetReportByID(ctx, reportID)

	// Assert
	assert.Error(t, err)
	// Note: The actual error might be different depending on how GORM handles context cancellation
}

func TestRepository_TurkishCharacters_ShouldWork(t *testing.T) {
	// Arrange
	db, mock := setupMockDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	firstName := "Ömer"
	lastName := "Çelik"

	// Mock the JOIN query with Turkish characters
	mock.ExpectQuery(`SELECT "reports"\."id","reports"\."created_at","reports"\."updated_at","reports"\."deleted_at","reports"\."diagnosis_title","reports"\."diagnosis_details","reports"\."image_path","reports"\."user_id","reports"\."patient_id" FROM "reports" JOIN patients ON reports\.patient_id = patients\.id WHERE \(patients\.name = \$1 AND patients\.lastname = \$2\)`).
		WithArgs(firstName, lastName).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "diagnosis_title", "diagnosis_details", "image_path", "user_id", "patient_id"}).
			AddRow(1, time.Now(), time.Now(), nil, "Türkçe Rapor", "Türkçe açıklama", "/turkce.pdf", 1, 2))

	// Mock User preload
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "firstname", "lastname", "email"}).
			AddRow(1, "Dr. Özge", "Şahin", "ozge@hastane.com"))

	// Mock Patient preload
	mock.ExpectQuery(`SELECT \* FROM "patients" WHERE`).
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "lastname", "national_id"}).
			AddRow(2, firstName, lastName, "12345678901"))

	// Act
	result, err := repo.GetReportsWithPatientName(ctx, firstName, lastName)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Türkçe Rapor", result[0].DiagnosisTitle)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// Benchmark tests
func BenchmarkRepository_CreateReport(b *testing.B) {
	db, mock := setupMockDB(&testing.T{})
	repo := NewRepository(db)
	ctx := context.Background()

	newReport := report.Report{
		DiagnosisTitle:   "Benchmark Report",
		DiagnosisDetails: "Benchmark Details",
		UserID:           1,
		PatientID:        2,
	}

	// Setup mock for all iterations
	for i := 0; i < b.N; i++ {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "reports"`).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil,
				newReport.DiagnosisTitle, newReport.DiagnosisDetails, newReport.ImagePath,
				newReport.UserID, newReport.PatientID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(i + 1))
		mock.ExpectCommit()

		// Mock preload queries
		mock.ExpectQuery(`SELECT \* FROM "reports" WHERE`).
			WithArgs(i + 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "diagnosis_title", "diagnosis_details", "image_path", "user_id", "patient_id"}).
				AddRow(i+1, newReport.DiagnosisTitle, newReport.DiagnosisDetails, newReport.ImagePath, newReport.UserID, newReport.PatientID))
		mock.ExpectQuery(`SELECT \* FROM "users" WHERE`).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "firstname", "lastname", "email"}).
				AddRow(1, "Dr. Benchmark", "User", "benchmark@hospital.com"))
		mock.ExpectQuery(`SELECT \* FROM "patients" WHERE`).
			WithArgs(2).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "lastname", "national_id"}).
				AddRow(2, "Benchmark", "Patient", "12345678901"))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.CreateReport(ctx, newReport)
	}
}
