package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// LoadDevDotEnv loads the .env.development file when ENV is unset or set to dev.
// The lookup walks up from the working directory (useful for go run) and
// falls back to the executable location for packaged binaries.
func LoadDevDotEnv() error {
	env := os.Getenv("ENV")
	if env != "" && env != "dev" {
		return nil
	}

	if env == "" {
		// Default to dev so downstream code knows we are in development mode.
		_ = os.Setenv("ENV", "dev")
	}

	path, err := locateEnvFile(".env.development")
	if err != nil {
		return err
	}

	return godotenv.Load(path)
}

func locateEnvFile(name string) (string, error) {
	if path, ok := findFileUpwards(name, workingDirectory); ok {
		return path, nil
	}

	if path, ok := findFileUpwards(name, executableDirectory); ok {
		return path, nil
	}

	return "", fmt.Errorf("%s not found in reachable directories", name)
}

type baseDirFunc func() (string, error)

func workingDirectory() (string, error) {
	return os.Getwd()
}

func executableDirectory() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exe), nil
}

func findFileUpwards(name string, fn baseDirFunc) (string, bool) {
	start, err := fn()
	if err != nil || start == "" {
		return "", false
	}

	dir := start
	for {
		candidate := filepath.Join(dir, name)
		if _, err := os.Stat(candidate); err == nil {
			return candidate, true
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", false
}
