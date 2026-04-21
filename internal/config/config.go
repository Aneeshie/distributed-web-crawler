package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	RedisAddr   string
	PostgresURL string
}

func Load() Config {

	// not a reference to config because its a small struct

	err := godotenv.Load()
	if err != nil {
		log.Println(".env not found, using system env var")
	}

	cfg := Config{
		RedisAddr:   os.Getenv("REDIS_ADDR"),
		PostgresURL: os.Getenv("POSTGRES_URL"),
	}

	if cfg.RedisAddr == "" || cfg.PostgresURL == "" {
		log.Fatal("missing environment variables")
	}

	return cfg

}
