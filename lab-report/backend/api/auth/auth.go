package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"gitbub.com/zikrullahcelep611/lab-report/backend/models/auth"
)

type AuthService interface {
	Login(ctx context.Context, email string, password string) (auth.Login, error)
	Logout(ctx context.Context, tokenString string) error
}

type JwtService interface {
	GenerateJwtToken(email string, expirationTime time.Time) (string, error)
}

type AuthHandler struct {
	authService AuthService
	jwtService  JwtService
}

func NewAuthController(authService AuthService, jwtService JwtService) *AuthHandler {
	return &AuthHandler{authService: authService, jwtService: jwtService}
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var creds auth.Login
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user, err := a.authService.Login(ctx, creds.Email, creds.Password)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(time.Hour * 24)
	tokenString, err := a.jwtService.GenerateJwtToken(user.Email, expirationTime)
	if err != nil {
		http.Error(w, "Could not create token", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
	w.WriteHeader(http.StatusOK)
}

func (a *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "Token not found", http.StatusBadRequest)
		return
	}
	tokenString := cookie.Value

	err = a.authService.Logout(ctx, tokenString)
	if err != nil {
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logged out successfully",
	})
}
