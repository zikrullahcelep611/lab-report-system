package claims

import (
	"gitbub.com/zikrullahcelep611/lab-report/backend/models/role"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Email string `json:"email"`
	Role role.Role `json:"role"`
	jwt.RegisteredClaims
}