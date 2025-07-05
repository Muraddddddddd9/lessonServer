package config

import (
	"fmt"
	"lesson_server/constants"
	"os"

	"github.com/joho/godotenv"
)

type ConfigStruct struct {
	ORIGIN      string `env:"ORIGIN"`
	DB_NAME     string `env:"DB_NAME"`
	DB_USERNAME string `env:"DB_USERNAME"`
	DB_PASSWORD string `env:"DB_PASSWORD"`
	DB_HOST     string `env:"DB_HOST"`
	DB_PORT     string `env:"DB_PORT"`
}

func ConfigLoad() (*ConfigStruct, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf(constants.ErrLoadEnv)
	}

	cfg := &ConfigStruct{
		ORIGIN:      getEnv("ORIGIN"),
		DB_NAME:     getEnv("DB_NAME"),
		DB_USERNAME: getEnv("DB_USERNAME"),
		DB_PASSWORD: getEnv("DB_PASSWORD"),
		DB_HOST:     getEnv("DB_HOST"),
		DB_PORT:     getEnv("DB_PORT"),
	}

	return cfg, nil
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		return ""
	}

	return value
}
