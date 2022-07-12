package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

/* Get the environment value */
func GoDotEnvValue(key string) string {
	err := godotenv.Load("./config/.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	value := os.Getenv(key)

	return value
}
