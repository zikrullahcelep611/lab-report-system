package token

import "gorm.io/gorm"

type ExpiredToken struct {
	gorm.Model
	Token      string `json:"token" gorm:"uniqueIndex"`
	ExpireTime int64  `json:"expire_time"`
}