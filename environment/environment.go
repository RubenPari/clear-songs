package environment

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

// LoadEnvVariables load environment variables from a file path
func LoadEnvVariables() {
	// get current working directory
	cwd, errCwd := os.Getwd()

	if errCwd != nil {
		log.Fatalf("error getting current working directory: %v", errCwd)
	}

	// check if on Linux/macOS
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		// move up one level folder
		cwd = filepath.Dir(cwd)
	}

	envPath := filepath.Join(cwd, ".env")

	errLoadFilePath := godotenv.Load(envPath)

	if errLoadFilePath != nil {
		log.Fatalf("error loading .env file: %v", errLoadFilePath)
	}

	log.Println("Loaded environment variables from .env file")
}
