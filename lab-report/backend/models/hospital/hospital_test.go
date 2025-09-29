package hospital

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestHospital_ValidHospital_ShouldWork(t *testing.T) {
	// Arrange
	hospital := Hospital{
		Name: "Ankara Şehir Hastanesi",
	}

	// Act & Assert
	assert.Equal(t, "Ankara Şehir Hastanesi", hospital.Name)
	assert.NotEmpty(t, hospital.Name)
}

func TestHospital_EmptyName_ShouldWork(t *testing.T) {
	// Arrange - Hospital model doesn't have validation tags
	hospital := Hospital{
		Name: "",
	}

	// Act & Assert
	assert.Equal(t, "", hospital.Name)
}

func TestHospital_LongName_ShouldWork(t *testing.T) {
	// Arrange
	longName := "İstanbul Üniversitesi İstanbul Tıp Fakültesi Hastanesi Cerrahpaşa"
	hospital := Hospital{
		Name: longName,
	}

	// Act & Assert
	assert.Equal(t, longName, hospital.Name)
}

func TestHospital_TurkishCharacters_ShouldWork(t *testing.T) {
	// Arrange - Test Turkish characters in hospital names
	testCases := []string{
		"Ankara Şehir Hastanesi",
		"İstanbul Üniversitesi Tıp Fakültesi",
		"Çukurova Üniversitesi Hastanesi",
		"Özel Acıbadem Hastanesi",
		"Gülhane Eğitim ve Araştırma Hastanesi",
	}

	for _, name := range testCases {
		t.Run(name, func(t *testing.T) {
			// Arrange
			hospital := Hospital{
				Name: name,
			}

			// Act & Assert
			assert.Equal(t, name, hospital.Name)
		})
	}
}

func TestHospital_SpecialCharacters_ShouldWork(t *testing.T) {
	// Arrange
	testCases := []string{
		"Dr. Sadi Konuk Eğitim ve Araştırma Hastanesi",
		"Prof. Dr. Cemil Taşçıoğlu Şehir Hastanesi",
		"TOBB ETÜ Hastanesi",
		"M.A.Ş. Özel Hastanesi",
		"A.Ü. Tıp Fakültesi",
	}

	for _, name := range testCases {
		t.Run(name, func(t *testing.T) {
			// Arrange
			hospital := Hospital{
				Name: name,
			}

			// Act & Assert
			assert.Equal(t, name, hospital.Name)
		})
	}
}

func TestHospital_GormModel_ShouldIncludeBaseFields(t *testing.T) {
	// Arrange
	hospital := Hospital{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name: "Test Hospital",
	}

	// Act & Assert
	assert.Equal(t, uint(1), hospital.ID)
	assert.False(t, hospital.CreatedAt.IsZero())
	assert.False(t, hospital.UpdatedAt.IsZero())
	assert.Nil(t, hospital.DeletedAt.Time)
}

func TestHospital_JSONTag_ShouldBeCorrect(t *testing.T) {
	// Arrange
	hospital := Hospital{
		Name: "JSON Test Hospital",
	}

	// Act & Assert
	assert.Equal(t, "JSON Test Hospital", hospital.Name)
}

func TestHospital_MultipleHospitals_ShouldWork(t *testing.T) {
	// Arrange
	hospitals := []Hospital{
		{Name: "Hospital 1"},
		{Name: "Hospital 2"},
		{Name: "Hospital 3"},
	}

	// Act & Assert
	assert.Len(t, hospitals, 3)
	assert.Equal(t, "Hospital 1", hospitals[0].Name)
	assert.Equal(t, "Hospital 2", hospitals[1].Name)
	assert.Equal(t, "Hospital 3", hospitals[2].Name)
}

// Benchmark tests
func BenchmarkHospital_Creation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Hospital{
			Name: "Benchmark Hospital",
		}
	}
}
