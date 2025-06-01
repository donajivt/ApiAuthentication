package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type JwtOptions struct {
	Issuer   string
	Audience string
	Secret   string
}

type Config struct {
	DSN        string
	JwtOptions JwtOptions
}

var Cfg Config

func Load() {
	_ = godotenv.Load()
	Cfg = Config{
		DSN: os.Getenv("DB_DSN"),
		JwtOptions: JwtOptions{
			Issuer:   os.Getenv("JWT_ISSUER"),
			Audience: os.Getenv("JWT_AUDIENCE"),
			Secret:   os.Getenv("JWT_SECRET"),
		},
	}
	if Cfg.DSN == "" {
		log.Fatal("DB_DSN not set")
	}
}
