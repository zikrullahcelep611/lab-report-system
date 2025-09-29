package report

import (
	"testing"
	"time"

	"gitbub.com/zikrullahcelep611/lab-report/backend/models/patient"
	"gitbub.com/zikrullahcelep611/lab-report/backend/models/user"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestReport_ValidReport_ShouldPassValidation(t *testing.T) {
	
	report := Report{
		DiagnosisTitle:   "Blood Test Results",
		DiagnosisDetails: "All parameters are within normal range",
		ImagePath:        "/uploads/blood_test_2024.png",
		UserID:           6,
		PatientID:        1,
	}

	validator := validator.New()

	// Act
	err := validator.Struct(report)

	// Assert
	assert.NoError(t, err)
}

func TestReport_EmptyDiagnosisTitle_ShouldFailValidation(t *testing.T) {
	
	report := Report{
		DiagnosisTitle:   "", // Empty title
		DiagnosisDetails: "Some details",
		UserID:           1,
		PatientID:        1,
	}

	validator := validator.New()

	// Act
	err := validator.Struct(report)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "DiagnosisTitle")
	assert.Contains(t, err.Error(), "required")
}

func TestReport_EmptyDiagnosisDetails_ShouldFailValidation(t *testing.T) {
	
	report := Report{
		DiagnosisTitle:   "Blood Test",
		DiagnosisDetails: "", // Empty details
		UserID:           1,
		PatientID:        1,
	}

	validator := validator.New()

	// Act
	err := validator.Struct(report)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "DiagnosisDetails")
	assert.Contains(t, err.Error(), "required")
}

func TestReport_EmptyImagePath_ShouldPassValidation(t *testing.T) {
	 
	report := Report{
		DiagnosisTitle:   "X-Ray Results",
		DiagnosisDetails: "No abnormalities detected",
		ImagePath:        "", // Empty image path (optional)
		UserID:           1,
		PatientID:        1,
	}

	validator := validator.New()

	// Act
	err := validator.Struct(report)

	// Assert
	assert.NoError(t, err, "ImagePath should be optional")
}

func TestReport_ZeroUserID_ShouldPassValidation(t *testing.T) {
	
	report := Report{
		DiagnosisTitle:   "MRI Results",
		DiagnosisDetails: "Brain scan shows normal activity",
		UserID:           0, // Zero UserID
		PatientID:        1,
	}

	validator := validator.New()

	// Act
	err := validator.Struct(report)

	// Assert
	assert.NoError(t, err, "UserID validation should be handled at business logic level")
}

func TestReport_ZeroPatientID_ShouldPassValidation(t *testing.T) {

	report := Report{
		DiagnosisTitle:   "CT Scan Results",
		DiagnosisDetails: "Chest CT shows clear lungs",
		UserID:           1,
		PatientID:        0, // Zero PatientID
	}

	validator := validator.New()

	// Act
	err := validator.Struct(report)

	// Assert
	assert.NoError(t, err, "PatientID validation should be handled at business logic level")
}

func TestReport_LongImagePath_ShouldWork(t *testing.T) {
	
	longPath := "/uploads/reports/very/long/path/to/medical/image/file/blood_test_results_2024_09_29_patient_12345.pdf"

	report := Report{
		DiagnosisTitle:   "Blood Analysis",
		DiagnosisDetails: "Complete blood count analysis results",
		ImagePath:        longPath,
		UserID:           1,
		PatientID:        1,
	}

	validator := validator.New()

	// Act
	err := validator.Struct(report)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, longPath, report.ImagePath)
}

func TestReport_WithRelationships_ShouldInitializeCorrectly(t *testing.T) {
	
	testUser := user.User{
		Model:     gorm.Model{ID: 1},
		Firstname: "Dr. John",
		Lastname:  "Doe",
		Email:     "john@hospital.com",
	}

	testPatient := patient.Patient{
		Model:      gorm.Model{ID: 1},
		Name:       "Jane",
		Lastname:   "Smith",
		NationalID: "12345678901",
	}

	report := Report{
		DiagnosisTitle:   "Comprehensive Health Check",
		DiagnosisDetails: "Annual health screening results",
		UserID:           1,
		User:             testUser,
		PatientID:        1,
		Patient:          testPatient,
	}

	// Act & Assert
	assert.Equal(t, uint(1), report.UserID)
	assert.Equal(t, uint(1), report.PatientID)
	assert.Equal(t, "Dr. John", report.User.Firstname)
	assert.Equal(t, "Jane", report.Patient.Name)
	assert.Equal(t, "12345678901", report.Patient.NationalID)
}

func TestReport_GormModel_ShouldIncludeBaseFields(t *testing.T) {
	
	report := Report{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		DiagnosisTitle:   "Test Report",
		DiagnosisDetails: "Test Details",
		UserID:           1,
		PatientID:        1,
	}

	// Act & Assert
	assert.Equal(t, uint(1), report.ID)
	assert.False(t, report.CreatedAt.IsZero())
	assert.False(t, report.UpdatedAt.IsZero())
}

// Benchmark tests for performance
func BenchmarkReport_Validation(b *testing.B) {
	report := Report{
		DiagnosisTitle:   "Performance Test Report",
		DiagnosisDetails: "Testing validation performance",
		UserID:           1,
		PatientID:        1,
	}

	validator := validator.New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validator.Struct(report)
	}
}
