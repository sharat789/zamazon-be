package configs

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
)

type AppConfig struct {
	Port           string
	DataSourceName string
	AppSecret      string
	JWTSecret      string
	CatalogURL     string
	AuthURL        string
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

	appSecret := os.Getenv("APP_SECRET")
	if len(appSecret) < 1 {
		return AppConfig{}, errors.New("app secret variable not found")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if len(jwtSecret) < 1 {
		return AppConfig{}, errors.New("jwt secret variable not found")
	}

	catalogURL := os.Getenv("CATALOG_URL")
	if len(catalogURL) < 1 {
		return AppConfig{}, errors.New("catalog url variable not found")
	}
	authURL := os.Getenv("AUTH_URL")
	if len(authURL) < 1 {
		return AppConfig{}, errors.New("auth url variable not found")
	}
	return AppConfig{Port: httpPort, DataSourceName: dsn, AppSecret: appSecret, JWTSecret: jwtSecret, CatalogURL: catalogURL, AuthURL: authURL}, nil
}
