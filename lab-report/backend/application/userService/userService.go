package userservice

import (
	"context"
	"errors"

	customErrors "gitbub.com/zikrullahcelep611/lab-report/backend/models/errors"
	"gitbub.com/zikrullahcelep611/lab-report/backend/models/user"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUser(ctx context.Context, id uint) (user.User, error)
	CreateUser(ctx context.Context, newUser user.User) (user.User, error)
	UpdateUser(ctx context.Context, updateUser user.User) (user.User, error)
	DeleteUser(ctx context.Context, id uint) error
	CheckUserExist(ctx context.Context, hospitalID uint) (bool, error)
}

type UserService struct {
	userRepository UserRepository
}

func NewUserService(userRepository UserRepository) *UserService{
	return &UserService{userRepository: userRepository}
}

func (u *UserService) GetUser(ctx context.Context, id uint) (user.User, error){
	usr, err := u.userRepository.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound){
			log.Warn().Str("operation", "GetUser").Err(err).Uint("user_id", id).Msg("User not found")
			return user.User{}, &customErrors.UserNotFoundError{Message:"User not found"}
		}
		log.Error().Str("operation", "GetUser").Err(err).Uint("user_id", id).Msg("Failed to retrieve user")
		return user.User{}, err
	}

	log.Info().Str("operation", "GetUser").Uint("user_id", id).Msg("User retrieved successfully")
	return usr, nil
}

func (u *UserService) CreateUser(ctx context.Context, newUser user.User) (user.User, error){
	if newUser.Email == "" || newUser.Password == "" || newUser.HospitalID == 0{
        return user.User{}, errors.New("email, password and hospital_id are required")
    }

    log.Info().
        Str("operation", "CreateUser").Str("email", newUser.Email).Msg("Creating new user")

	newUser.Password = u.HashPassword(newUser.Password)
	usr, err := u.userRepository.CreateUser(ctx, newUser)
	if err != nil {
		log.Error().
			Str("operation", "CreateUser").
			Err(err).
			Msg("Failed to create user")
		return user.User{}, err
	}

	log.Info().
		Str("operation", "CreateUser").
		Uint("user_id", usr.ID).
		Msg("User created successfully")

	return usr, nil
}

func (u *UserService) HashPassword(password string) (string){
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword)
} 

func (u *UserService) DeleteUser(ctx context.Context, id uint) error {
	log.Info().Str("operation", "DeleteUser").Uint("user_id", id).Msg("Deleting user")

	err := u.userRepository.DeleteUser(ctx, id)
	if err != nil {
		log.Error().
			Str("operation", "DeleteUser").
			Err(err).
			Uint("user_id", id).
			Msg("Failed to delete user")
		return err
	}

	log.Info().
		Str("operation", "DeleteUser").
		Uint("user_id", id).
		Msg("User deleted successfully")

	return nil
}

func (u *UserService) UpdateUser(ctx context.Context, updateUser user.User) (user.User, error) {
    log.Info().Str("operation", "UpdateUser").Uint("user_id", updateUser.ID).Msg("Updating user")
    updateUser.Password = u.HashPassword(updateUser.Password)

    usr, err := u.userRepository.UpdateUser(ctx, updateUser)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            log.Warn().Str("operation", "UpdateUser").Uint("user_id", updateUser.ID).Msg("User not found for update")
            return user.User{}, &customErrors.UserNotFoundError{Message: "User not found"}
        }
        log.Error().Str("operation", "UpdateUser").Err(err).Uint("user_id", updateUser.ID).Msg("Failed to update user")
        return user.User{}, err
    }

    log.Info().Str("operation", "UpdateUser").Uint("user_id", usr.ID).Msg("User updated successfully")
    return usr, nil
}

func (s *UserService) CheckUserExist(ctx context.Context, userModel user.User) (bool, error) {
	log.Info().
		Str("operation", "CheckUserExist").
		Str("email", userModel.Email).
		Uint("hospital_id", userModel.HospitalID).
		Uint("id", userModel.ID).
		Msg("Checking if user exists")

	exists, err := s.userRepository.CheckUserExist(ctx, userModel.HospitalID)
	if err != nil {
		log.Error().
			Str("operation", "CheckUserExist").
			Err(err).
			Msg("Failed to check user existence")
		return false, err
	}

	if exists {
		log.Info().
			Str("operation", "CheckUserExist").
			Msg("User exists")
	} else {
		log.Info().
			Str("operation", "CheckUserExist").
			Msg("User does not exist")
	}

	return exists, nil
}
