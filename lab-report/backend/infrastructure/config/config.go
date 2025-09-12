package config

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const (
	ConfigName = "application"
	ConfigType = "yml"
	DefaultEnv = "qa"
)

func LoadConfig(configPath string) *ConfigModel {
	env := os.Getenv("ENV")
	if env == "" {
		log.Warn().Msg("ENV is not set, using default env")
		env = DefaultEnv
	}
	viper.AddConfigPath(configPath)
	viper.SetConfigType(ConfigType)

	data := readConfig(env)
	return data
}

func readConfig(env string) *ConfigModel {
	viper.SetConfigName(ConfigName)
	readConfigErr := viper.ReadInConfig()
	if readConfigErr != nil {
		log.Fatal().Err(readConfigErr).Msg("Error reading config file")
		return nil
	}

	config := &ConfigModel{}
	v := viper.Sub(env)

	unMarshallErr := v.Unmarshal(config)
	if unMarshallErr != nil {
		log.Fatal().Err(unMarshallErr).Msg("Error unmarshalling config file")
		return nil
	}

	return config
}