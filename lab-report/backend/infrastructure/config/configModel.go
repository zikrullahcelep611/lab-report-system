package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

type ConfigModel struct {
	Server   ServerConfig   `yaml:"server" validate:"required"`
	Database DatabaseConfig `yaml:"database" validate:"required"`
	Email    EmailConfig    `yaml:"email" validate:"required"`
	Log      LogConfig      `yaml:"log" validate:"required"`
	Redis    RedisConfig    `yaml:"redis" validate:"required"`
	JWT      JWTConfig      `validate:"required"`
}

type ServerConfig struct {
	Port int `yaml:"port" validate:"min=0"`
}

type LogConfig struct {
	Level zerolog.Level `yaml:"level" validate:"required"`
}

type DatabaseConfig struct {
	DNS string `yaml:"dns" validate:"required"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr" validate:"required"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db" validate:"min=0"`
}

type EmailConfig struct {
	Host     string `json:"host" validate:"required"`
	Port     int    `json:"port" validate:"required, gt=0"`
	User     string `json:"user" validate:"required, email"`
	Password string `json:"password" validate:"required"`
}

type JWTConfig struct {
	SecretKey string `validate:"required"`
}

func (c *ConfigModel) ValidateConfig() error {
	validate := validator.New()
	err := validate.Struct(c)
	if err != nil {
		return err
	}

	return nil
}
