package configs

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
)

type AppConfig struct {
	Port           string
	DataSourceName string
	JWTSecret      string
	AppSecret      string
	StripeSecret   string
	PubKey         string
	SuccessURL     string
	CancelURL      string
	UserServiceURL string
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

	UserServiceURL := os.Getenv("USER_SERVICE_URL")
	if len(UserServiceURL) < 1 {
		return AppConfig{}, errors.New("user service URL not found")
	}

	AuthURL := os.Getenv("AUTH_SERVICE_URL")
	if len(AuthURL) < 1 {
		return AppConfig{}, errors.New("auth service URL not found")
	}
	return AppConfig{Port: httpPort, DataSourceName: dsn, AppSecret: appSecret,
		JWTSecret:      jwtSecret,
		StripeSecret:   os.Getenv("STRIPE_API_KEY"),
		PubKey:         os.Getenv("STRIPE_PUB_KEY"),
		SuccessURL:     os.Getenv("SUCCESS_URL"),
		CancelURL:      os.Getenv("CANCEL_URL"),
		UserServiceURL: UserServiceURL,
		AuthURL:        AuthURL}, nil

}
