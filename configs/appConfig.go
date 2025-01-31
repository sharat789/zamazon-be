package configs

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
)

type AppConfig struct {
	Port           string
	DataSourceName string
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

	dsn := os.Getenv("DSN")
	if len(dsn) < 1 {
		return AppConfig{}, errors.New("no env variables found")
	}
	return AppConfig{Port: httpPort, DataSourceName: dsn}, nil
}
