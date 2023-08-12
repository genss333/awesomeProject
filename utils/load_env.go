package utils

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func LoadDotEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func GoDotEnvVariable(key string) string {
	LoadDotEnv()
	return os.Getenv(key)
}
