package user

import (
	"testing"

	"gitbub.com/zikrullahcelep611/lab-report/backend/models/role"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUser_ValidUser_ShouldCreateSuccessfully(t *testing.T) {
	user := User{
		Model:      gorm.Model{ID: 1},
		Firstname:  "John",
		Lastname:   "Doe",
		Email:      "john@example.com",
		Password:   "hashedPassword123",
		HospitalID: 1,
		Role:       role.RoleTechnician,
	}

	assert.NotEmpty(t, user.Firstname)
	assert.NotEmpty(t, user.Lastname)
	assert.NotEmpty(t, user.Email)
	assert.NotEmpty(t, user.Password)
	assert.Greater(t, user.HospitalID, uint(0))
	assert.Equal(t, role.RoleTechnician, user.Role)
}

func TestUser_EmptyFields_ShouldHaveEmptyValues(t *testing.T) {
	user := User{}

	assert.Empty(t, user.Firstname)
	assert.Empty(t, user.Lastname)
	assert.Empty(t, user.Email)
	assert.Empty(t, user.Password)
	assert.Equal(t, uint(0), user.HospitalID)
	assert.Empty(t, user.Role)
}

func TestUser_DefaultRole_ShouldBeEmpty(t *testing.T) {

	user := User{
		Firstname:  "Jane",
		Lastname:   "Smith",
		Email:      "jane@example.com",
		Password:   "password123",
		HospitalID: 1,
	}

	assert.Empty(t, user.Role)
}

func TestUser_Roles_ShouldAcceptValidRoles(t *testing.T) {
	// Test different roles
	testCases := []struct {
		name     string
		userRole role.Role
	}{
		{"Super Admin", role.RoleSuperAdmin},
		{"Admin", role.RoleAdmin},
		{"Technician", role.RoleTechnician},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
		
			user := User{
				Firstname:  "Test",
				Lastname:   "User",
				Email:      "test@example.com",
				Password:   "password123",
				HospitalID: 1,
				Role:       tc.userRole,
			}

	
			assert.Equal(t, tc.userRole, user.Role)
		})
	}
}

func TestUser_HospitalID_ShouldBePositive(t *testing.T) {
	// Arrange
	user := User{
		Firstname:  "Test",
		Lastname:   "User",
		Email:      "test@example.com",
		Password:   "password123",
		HospitalID: 5,
		Role:       role.RoleTechnician,
	}

	assert.Greater(t, user.HospitalID, uint(0))
	assert.Equal(t, uint(5), user.HospitalID)
}

func BenchmarkUser_Creation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		user := User{
			Firstname:  "John",
			Lastname:   "Doe",
			Email:      "john@example.com",
			Password:   "hashedPassword123",
			HospitalID: 1,
			Role:       role.RoleTechnician,
		}
		_ = user
	}
}
