package config

import (
	"log"
	"os"
)

type Config struct {
	JWTSecret string
}

func LoadConfig() *Config {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Println("Warning: JWT_SECRET environment variable not set, using default")
		jwtSecret = "qhx5shDMjBChXdPrxzFmCr+W09Dz4uPvUGUp66xCrCs="
	}

	return &Config{
		JWTSecret: jwtSecret,
	}
}
