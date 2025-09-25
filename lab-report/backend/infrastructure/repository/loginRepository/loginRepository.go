package loginrepository

import (
	"context"
	"errors"

	"gitbub.com/zikrullahcelep611/lab-report/backend/models/auth"
	"gitbub.com/zikrullahcelep611/lab-report/backend/models/user"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (l *Repository) Login(ctx context.Context, email string, password string) (auth.Login, error) {
	var usr user.User
	result := l.DB.WithContext(ctx).Where("email = ?", email).First(&usr)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Warn().Str("email", email).Msg("Login attempt with non existent email")
			return auth.Login{}, result.Error
		}
		log.Error().Err(result.Error).Str("email", email).Msg("Failed to retreive user during login")
		return auth.Login{}, result.Error
	}

	err := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			log.Warn().Str("email", email).Str("Hashed_password", usr.Password).Str("password", password).Msg("Incorrect password attempt")
			return auth.Login{}, err
		}
		log.Error().Err(err).Str("email", email).Msg("Error comparing passwords")
		return auth.Login{}, err
	}

	log.Info().Str("email", email).Msg("User authenticated successfully")
	return auth.Login{
		Email: usr.Email,
	}, nil
}
