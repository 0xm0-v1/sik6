// config/envload.go
package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// LoadDevDotEnv loads the .env.development file only when ENV=dev.
// It resolves the executable's directory, goes two levels up,
// and attempts to load the environment variables from there.
func LoadDevDotEnv() error {
	// Skip if not running in development mode
	if os.Getenv("ENV") != "dev" {
		return nil
	}

	// Get the absolute path of the running executable
	exe, err := os.Executable()
	if err != nil {
		return err
	}

	// Construct the path to .env.development relative to the executable
	envPath := filepath.Join(filepath.Dir(exe), "..", "..", ".env.development")

	// Load the environment variables from the file
	if err := godotenv.Load(envPath); err != nil {
		return err
	}

	return nil
}
