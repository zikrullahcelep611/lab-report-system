package tokenrepository

import (
	"context"
	"errors"
	"time"

	"gitbub.com/zikrullahcelep611/lab-report/backend/models/token"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)


type Repository struct{
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) (*Repository){
	return &Repository{DB: db}
}

func (t *Repository) DeleteExpiredTokens(ctx context.Context){
	result := t.DB.WithContext(ctx).Unscoped().Where("expire_time < ?", time.Now().Unix()).
	Delete(&token.ExpiredTokens{})

	if result.Error != nil{
		log.Error().Str("operation", "DeleteExpiredTokens").Err(result.Error).
		Msg("Failed to delete expired tokens")
	}

	log.Info().Str("operation", "DeleteExpiredTokens").Int64("deletec_count", result.RowsAffected).
	Msg("Expired token deleted successfully")
}

func (t *Repository) AddTokenToBlackList(ctx context.Context, tokenStr string, expireTime time.Time) error{
	newToken := token.ExpiredTokens{
		Token: tokenStr,
		ExpireTime: expireTime.Unix(),
	}

	result := t.DB.WithContext(ctx).Create(&newToken)
	if result.Error != nil{
		log.Error().Str("operation", "AddTokenToBlackList").Err(result.Error).Msg("Failed to add token to blacklist")
		return errors.New("redis set errors")
	}

	log.Info().Str("operation", "AddTokenToBlackList").Str("token", tokenStr).Int64("expire_time", newToken.ExpireTime).
	Msg("Token added to blacklist successfully")
   	return nil
}

func (t *Repository) IsTokenBlackListed(ctx context.Context, tokenStr string) bool{
	var count int64
	result := t.DB.WithContext(ctx).Model(&token.ExpiredTokens{}).Where("token = ?", tokenStr).
	Count(&count)

	if result.Error != nil{
		log.Error().Str("operation", "IsTokenBlackListed").Err(result.Error).Str("token", tokenStr).
		Msg("Failed to check if token blacklisted")
		return false
	}

	if count > 0{
		log.Info().Str("operation", "IsTokenBlackListed").Str("token", tokenStr).
		Msg("Token is blacklisted")
	}else{
		log.Info().
			Str("operation", "IsTokenBlacklisted").
			Str("token", tokenStr).
			Msg("Token is not blacklisted")
	}

	return count > 0
}