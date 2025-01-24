package configs

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
)

type AppConfig struct {
	Port string
}

func EnvSetup() (cfg AppConfig, err error) {
	if os.Getenv("APP_ENV") == "dev" {
		envErr := godotenv.Load()
		if envErr != nil {
			return AppConfig{}, envErr
		}
	}
	httpPort := os.Getenv("HTTP_PORT")
	if len(httpPort) < 1 {
		return AppConfig{}, errors.New("no env variables found")
	}
	return AppConfig{Port: httpPort}, nil
}
