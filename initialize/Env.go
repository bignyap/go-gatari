package initialize

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func GetEnvVals() error {
	// Check if .env file exists
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		fmt.Println(".env file not found, skipping loading environment variables from file")
		return nil
	}

	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}

	return nil
}
