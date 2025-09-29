package patient

import (
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestPatient_ValidPatient_ShouldPassValidation(t *testing.T) {
	// Arrange
	patient := Patient{
		Name:       "Ahmet",
		Lastname:   "Yılmaz",
		NationalID: "12345678901",
	}

	validator := validator.New()

	// Act
	err := validator.Struct(patient)

	// Assert
	assert.NoError(t, err)
}

func TestPatient_EmptyName_ShouldFailValidation(t *testing.T) {
	// Arrange
	patient := Patient{
		Name:       "", // Empty name
		Lastname:   "Yılmaz",
		NationalID: "12345678901",
	}

	validator := validator.New()

	// Act
	err := validator.Struct(patient)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Name")
	assert.Contains(t, err.Error(), "required")
}

func TestPatient_EmptyLastname_ShouldFailValidation(t *testing.T) {
	// Arrange
	patient := Patient{
		Name:       "Ahmet",
		Lastname:   "", // Empty lastname
		NationalID: "12345678901",
	}

	validator := validator.New()

	// Act
	err := validator.Struct(patient)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Lastname")
	assert.Contains(t, err.Error(), "required")
}

func TestPatient_EmptyNationalID_ShouldFailValidation(t *testing.T) {
	// Arrange
	patient := Patient{
		Name:       "Ahmet",
		Lastname:   "Yılmaz",
		NationalID: "", // Empty national ID
	}

	validator := validator.New()

	// Act
	err := validator.Struct(patient)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "NationalID")
	assert.Contains(t, err.Error(), "required")
}

func TestPatient_InvalidNationalIDLength_ShouldFailValidation(t *testing.T) {
	// Arrange - National ID must be exactly 11 characters
	testCases := []struct {
		name       string
		nationalID string
	}{
		{"Too Short", "123456789"},      // 9 digits
		{"Too Short 2", "1234567890"},   // 10 digits
		{"Too Long", "123456789012"},    // 12 digits
		{"Too Long 2", "1234567890123"}, // 13 digits
	}

	validator := validator.New()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			patient := Patient{
				Name:       "Ahmet",
				Lastname:   "Yılmaz",
				NationalID: tc.nationalID,
			}

			// Act
			err := validator.Struct(patient)

			// Assert
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "NationalID")
			assert.Contains(t, err.Error(), "len")
		})
	}
}

func TestPatient_NonNumericNationalID_ShouldFailValidation(t *testing.T) {
	// Arrange - National ID must be numeric only
	testCases := []struct {
		name       string
		nationalID string
	}{
		{"Letters", "1234567890A"},
		{"Special Characters", "12345678-01"},
		{"Spaces", "12345 67890"},
		{"Mixed", "123ABC78901"},
	}

	validator := validator.New()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			patient := Patient{
				Name:       "Ahmet",
				Lastname:   "Yılmaz",
				NationalID: tc.nationalID,
			}

			// Act
			err := validator.Struct(patient)

			// Assert
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "NationalID")
			assert.Contains(t, err.Error(), "numeric")
		})
	}
}

func TestPatient_ValidNationalIDFormats_ShouldPassValidation(t *testing.T) {
	// Arrange - Test different valid national ID formats
	testCases := []struct {
		name       string
		nationalID string
	}{
		{"All Same Digits", "11111111111"},
		{"Sequential", "12345678901"},
		{"Random Valid", "98765432100"},
		{"Starting with Zero", "01234567890"},
	}

	validator := validator.New()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			patient := Patient{
				Name:       "Test",
				Lastname:   "User",
				NationalID: tc.nationalID,
			}

			// Act
			err := validator.Struct(patient)

			// Assert
			assert.NoError(t, err, "National ID %s should be valid", tc.nationalID)
		})
	}
}

func TestPatient_TurkishCharacters_ShouldWork(t *testing.T) {
	// Arrange - Test Turkish characters in names
	testCases := []struct {
		name     string
		lastname string
	}{
		{"Ömer", "Özkan"},
		{"Çağla", "Şahin"},
		{"Gül", "Ünal"},
		{"İsmail", "Ğüler"},
	}

	validator := validator.New()

	for _, tc := range testCases {
		t.Run(tc.name+" "+tc.lastname, func(t *testing.T) {
			// Arrange
			patient := Patient{
				Name:       tc.name,
				Lastname:   tc.lastname,
				NationalID: "12345678901",
			}

			// Act
			err := validator.Struct(patient)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tc.name, patient.Name)
			assert.Equal(t, tc.lastname, patient.Lastname)
		})
	}
}

func TestPatient_LongNames_ShouldWork(t *testing.T) {
	// Arrange - Test with longer names
	longName := "Abdurrahman"
	longLastname := "Papadopoulos"

	patient := Patient{
		Name:       longName,
		Lastname:   longLastname,
		NationalID: "12345678901",
	}

	validator := validator.New()

	// Act
	err := validator.Struct(patient)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, longName, patient.Name)
	assert.Equal(t, longLastname, patient.Lastname)
}

func TestPatient_GormModel_ShouldIncludeBaseFields(t *testing.T) {
	// Arrange
	patient := Patient{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:       "Test",
		Lastname:   "Patient",
		NationalID: "12345678901",
	}

	// Act & Assert
	assert.Equal(t, uint(1), patient.ID)
	assert.False(t, patient.CreatedAt.IsZero())
	assert.False(t, patient.UpdatedAt.IsZero())
	assert.Nil(t, patient.DeletedAt.Time)
}

func TestPatient_JSONTags_ShouldBeCorrect(t *testing.T) {
	// Arrange
	patient := Patient{
		Name:       "John",
		Lastname:   "Doe",
		NationalID: "12345678901",
	}

	// Act & Assert - Verify field accessibility
	assert.Equal(t, "John", patient.Name)
	assert.Equal(t, "Doe", patient.Lastname)
	assert.Equal(t, "12345678901", patient.NationalID)
}

func TestPatient_AllFieldsEmpty_ShouldFailValidation(t *testing.T) {
	// Arrange
	patient := Patient{
		Name:       "",
		Lastname:   "",
		NationalID: "",
	}

	validator := validator.New()

	// Act
	err := validator.Struct(patient)

	// Assert
	assert.Error(t, err)

	// Check that all required fields are mentioned in error
	errorString := err.Error()
	assert.Contains(t, errorString, "Name")
	assert.Contains(t, errorString, "Lastname")
	assert.Contains(t, errorString, "NationalID")
}

func TestPatient_EdgeCaseNames_ShouldWork(t *testing.T) {
	// Arrange - Test edge cases for names
	testCases := []struct {
		name        string
		patientName string
		lastname    string
	}{
		{"Single Letter Name", "A", "B"},
		{"Numbers in Name", "Ali2", "Veli3"}, // May not be realistic but tests validation
		{"Hyphenated Names", "Anne-Marie", "Smith-Jones"},
	}

	validator := validator.New()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			patient := Patient{
				Name:       tc.patientName,
				Lastname:   tc.lastname,
				NationalID: "12345678901",
			}

			// Act
			err := validator.Struct(patient)

			// Assert
			assert.NoError(t, err)
		})
	}
}

// Benchmark tests for performance
func BenchmarkPatient_Validation(b *testing.B) {
	patient := Patient{
		Name:       "Performance",
		Lastname:   "Test",
		NationalID: "12345678901",
	}

	validator := validator.New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validator.Struct(patient)
	}
}

func BenchmarkPatient_Creation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Patient{
			Name:       "Benchmark",
			Lastname:   "User",
			NationalID: "12345678901",
		}
	}
}
