package auth

import (
	"gitbub.com/zikrullahcelep611/lab-report/backend/models/role"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type Login struct {
	gorm.Model
	Email    string `json:"email"`
	Role 	 role.Role `json:"role" `
	Password string `json:"passsword"`
	jwt.RegisteredClaims
}