package userRepository

import (
	"context"

	"gitbub.com/zikrullahcelep611/lab-report/backend/models/user"
	"github.com/rs/zerolog/log"

	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (u *Repository) GetUser(ctx context.Context, id uint) (user.User, error) {
	var usr user.User
	result := u.DB.WithContext(ctx).Where("id = ?", id).First(&usr)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Warn().Str("operation", "GetUser").Err(result.Error).
				Uint("user_id", id).Msg("User not found")
		} else {
			log.Error().Str("operation", "GetUser").Err(result.Error).
				Uint("user_id", id).Msg("Failed to retreive user")
		}
		return user.User{}, result.Error
	}

	log.Info().Str("operation", "GetUser").Uint("user_id", id).Msg("Retreieved user successfully")
	return usr, nil
}

func (u *Repository) CreateUser(ctx context.Context, newUser user.User) (user.User, error) {
	result := u.DB.WithContext(ctx).Create(&newUser)
	if result.Error != nil {
		log.Error().Str("operation", "CreateUser").Err(result.Error).Msg("Failed to create user")
		return user.User{}, result.Error
	}
	log.Info().Str("operation", "CreateUser").Uint("user_id", newUser.ID).Msg("User created successfully")
	return newUser, nil
}

func (u *Repository) UpdateUser(ctx context.Context, updateUser user.User) (user.User, error) {
	result := u.DB.WithContext(ctx).Save(&updateUser)
	if result.Error != nil {
		log.Error().Str("operation", "UpdateUser").Err(result.Error).Uint("user_id", updateUser.ID).Msg("Failed to update user")
		return user.User{}, nil
	}

	log.Info().Str("operation", "UpdateUser").Uint("user_id", updateUser.ID).Msg("User updated successfully")
	return updateUser, nil
}

func (u *Repository) DeleteUser(ctx context.Context, id uint) error {
	result := u.DB.WithContext(ctx).Delete(&user.User{}, id)
	if result.Error != nil {
		log.Error().Str("operation", "DeleteUser").Err(result.Error).Uint("user_id", id).Msg("Failed to delete user")
		return result.Error
	}

	log.Info().Str("operation", "DeleteUser").Uint("user_id", id).Msg("User deleted successfully")
	return nil
}

func (u *Repository) CheckUserExist(ctx context.Context, hospitalID uint) (bool, error) {
	var count int64
	result := u.DB.WithContext(ctx).Model(&user.User{}).Where("hospital_id = ?", hospitalID).Count(&count)
	if result.Error != nil {
		log.Error().Str("operation", "CheckUserExist").Err(result.Error).Uint("hospital_id", hospitalID).
			Msg("Failed to check user existince")
		return false, result.Error
	}

	exists := count > 0
	if exists {
		log.Info().Str("operation", "CheckUserExist").
			Msgf("User exists with Hospital ID = %d", hospitalID)
	} else {
		log.Info().Str("operation", "CheckUserExist").Msg("User does not exist")
	}
	return exists, nil
}
