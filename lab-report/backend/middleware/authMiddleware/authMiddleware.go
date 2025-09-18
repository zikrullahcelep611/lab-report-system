package authmiddleware

import (
	"context"
	"errors"
	"net/http"

	"gitbub.com/zikrullahcelep611/lab-report/backend/models/claims"
	"github.com/golang-jwt/jwt/v5"
)

type TokenService interface{
	IsTokenBlackListed(ctx context.Context, tokenString string) bool
}

type JwtService interface{
	GetJwtKey() []byte
}

type AuthMiddleware struct{
	jwtService JwtService
	tokenService TokenService
}

func NewAuthMiddleware(tokenService TokenService, jwtService JwtService) *AuthMiddleware{
	return &AuthMiddleware{tokenService: tokenService, jwtService: jwtService}
}

func (auth *AuthMiddleware) Authenticate(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		ctx := r.Context()
		cookie, err := r.Cookie("token")
		if err != nil{
			if errors.Is(err, http.ErrNoCookie){
				w.WriteHeader(http.StatusUnauthorized)
				return 
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if auth.tokenService.IsTokenBlackListed(ctx, cookie.Value){
			w.WriteHeader(http.StatusUnauthorized)
			return 
		}

		claims := &claims.Claims{}
		token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token)(interface{}, error){
			return auth.jwtService.GetJwtKey(), nil
		})

		if err != nil{
			if errors.Is(err, jwt.ErrSignatureInvalid){
				w.WriteHeader(http.StatusUnauthorized)
				return 
			}
			w.WriteHeader(http.StatusBadRequest)
			return 
		}

		if !token.Valid{
			w.WriteHeader(http.StatusUnauthorized)
			return 
		}
		next.ServeHTTP(w, r)
	})
}
