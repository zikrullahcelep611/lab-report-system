package role

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRole_Constants_ShouldHaveCorrectValues(t *testing.T) {
	// Act & Assert
	assert.Equal(t, Role("super_admin"), RoleSuperAdmin)
	assert.Equal(t, Role("admin"), RoleAdmin)
	assert.Equal(t, Role("technician"), RoleTechnician)
}

func TestRole_StringConversion_ShouldWork(t *testing.T) {
	// Arrange
	testCases := []struct {
		role     Role
		expected string
	}{
		{RoleSuperAdmin, "super_admin"},
		{RoleAdmin, "admin"},
		{RoleTechnician, "technician"},
	}

	for _, tc := range testCases {
		t.Run(string(tc.role), func(t *testing.T) {
			// Act
			result := string(tc.role)

			// Assert
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestRole_Comparison_ShouldWork(t *testing.T) {
	// Arrange
	role1 := RoleSuperAdmin
	role2 := RoleSuperAdmin
	role3 := RoleAdmin

	// Act & Assert
	assert.True(t, role1 == role2, "Same roles should be equal")
	assert.False(t, role1 == role3, "Different roles should not be equal")
}

func TestRole_Assignment_ShouldWork(t *testing.T) {
	// Arrange
	var userRole Role

	// Act
	userRole = RoleAdmin

	// Assert
	assert.Equal(t, RoleAdmin, userRole)
	assert.Equal(t, "admin", string(userRole))
}

func TestRole_FromString_ShouldWork(t *testing.T) {
	// Arrange
	testCases := []struct {
		input    string
		expected Role
	}{
		{"super_admin", RoleSuperAdmin},
		{"admin", RoleAdmin},
		{"technician", RoleTechnician},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			// Act
			role := Role(tc.input)

			// Assert
			assert.Equal(t, tc.expected, role)
		})
	}
}

func TestRole_InvalidRole_ShouldNotMatchConstants(t *testing.T) {
	// Arrange
	invalidRole := Role("invalid_role")

	// Act & Assert
	assert.NotEqual(t, RoleSuperAdmin, invalidRole)
	assert.NotEqual(t, RoleAdmin, invalidRole)
	assert.NotEqual(t, RoleTechnician, invalidRole)
}

func TestRole_EmptyRole_ShouldWork(t *testing.T) {
	// Arrange
	emptyRole := Role("")

	// Act & Assert
	assert.Equal(t, "", string(emptyRole))
	assert.NotEqual(t, RoleSuperAdmin, emptyRole)
}

func TestRole_CaseSensitive_ShouldWork(t *testing.T) {
	// Arrange
	upperCaseRole := Role("ADMIN")
	mixedCaseRole := Role("Admin")

	// Act & Assert
	assert.NotEqual(t, RoleAdmin, upperCaseRole)
	assert.NotEqual(t, RoleAdmin, mixedCaseRole)
	assert.Equal(t, "admin", string(RoleAdmin))
}

func TestRole_AllRoles_ShouldBeDifferent(t *testing.T) {
	// Arrange
	roles := []Role{RoleSuperAdmin, RoleAdmin, RoleTechnician}

	// Act & Assert
	for i, role1 := range roles {
		for j, role2 := range roles {
			if i != j {
				assert.NotEqual(t, role1, role2, "Roles should be different: %s vs %s", role1, role2)
			}
		}
	}
}

func TestRole_IsValidRole_ShouldWork(t *testing.T) {
	// Helper function to check if role is valid
	isValidRole := func(r Role) bool {
		return r == RoleSuperAdmin || r == RoleAdmin || r == RoleTechnician
	}

	// Arrange
	testCases := []struct {
		role     Role
		expected bool
	}{
		{RoleSuperAdmin, true},
		{RoleAdmin, true},
		{RoleTechnician, true},
		{Role("invalid"), false},
		{Role(""), false},
		{Role("user"), false},
	}

	for _, tc := range testCases {
		t.Run(string(tc.role), func(t *testing.T) {
			// Act
			result := isValidRole(tc.role)

			// Assert
			assert.Equal(t, tc.expected, result)
		})
	}
}

// Benchmark tests
func BenchmarkRole_StringConversion(b *testing.B) {
	role := RoleAdmin

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = string(role)
	}
}

func BenchmarkRole_Comparison(b *testing.B) {
	role1 := RoleAdmin
	role2 := RoleAdmin

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = role1 == role2
	}
}
