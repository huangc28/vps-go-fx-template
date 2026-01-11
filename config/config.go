package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AppName         string
	Env             string
	Port            int
	LogLevel        string
	ShutdownTimeout time.Duration

	DB DBConfig
}

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	Name     string
}

func NewViper() *viper.Viper {
	v := viper.New()

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	v.SetDefault("APP_NAME", "go-vps-service")
	v.SetDefault("APP_ENV", "development")
	v.SetDefault("APP_PORT", 8080)
	v.SetDefault("LOG_LEVEL", "info")
	v.SetDefault("APP_SHUTDOWN_TIMEOUT", "10s")

	v.SetDefault("DB_USER", "")
	v.SetDefault("DB_PASSWORD", "")
	v.SetDefault("DB_HOST", "")
	v.SetDefault("DB_PORT", 5432)
	v.SetDefault("DB_NAME", "")

	return v
}

func NewConfig(v *viper.Viper) (Config, error) {
	timeout, err := time.ParseDuration(v.GetString("APP_SHUTDOWN_TIMEOUT"))
	if err != nil {
		return Config{}, err
	}

	return Config{
		AppName:         v.GetString("APP_NAME"),
		Env:             v.GetString("APP_ENV"),
		Port:            v.GetInt("APP_PORT"),
		LogLevel:        v.GetString("LOG_LEVEL"),
		ShutdownTimeout: timeout,
		DB: DBConfig{
			User:     v.GetString("DB_USER"),
			Password: v.GetString("DB_PASSWORD"),
			Host:     v.GetString("DB_HOST"),
			Port:     v.GetInt("DB_PORT"),
			Name:     v.GetString("DB_NAME"),
		},
	}, nil
}
