// config/config.go
package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func ReadEnv(value string) string {
	rootDir := os.Getenv("ROOT_DIR")
	err := godotenv.Load(filepath.Join(rootDir, ".env"))
	if err != nil {
		log.Printf("Error loading .env file: %s", err)
	}
	return os.Getenv(value)
}
