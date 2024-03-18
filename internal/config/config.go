package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecretKey string
}

func LoadConfig() *Config {
	// Load .env file
	err := godotenv.Load() // This will look for a .env file in the current directory
	if err != nil {
		log.Println("No .env file found or error loading .env file")
	}

	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	if jwtSecretKey == "" {
		log.Fatal("JWT_SECRET_KEY environment variable is not set")
	}

	return &Config{
		JWTSecretKey: jwtSecretKey,
	}
}
