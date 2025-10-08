package fixtures

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	fixtureDir      = "testdata/fixtures"
	defaultErrorMsg = "fixture %q not found; looked in %s"
)

// Seed executes the requested SQL fixture files using the supplied connection pool.
// Filenames should reference resources located under backend/testdata/fixtures.
func Seed(ctx context.Context, pool *pgxpool.Pool, filenames ...string) error {
	if len(filenames) == 0 {
		return nil
	}

	for _, name := range filenames {
		if err := executeFixture(ctx, pool, name); err != nil {
			return err
		}
	}

	return nil
}

func executeFixture(ctx context.Context, pool *pgxpool.Pool, filename string) error {
	name := strings.TrimSpace(filename)
	if name == "" {
		return fmt.Errorf("fixture filename is required")
	}

	content, err := loadFixture(name)
	if err != nil {
		return err
	}

	statements := strings.TrimSpace(string(content))
	if statements == "" {
		return nil
	}

	if _, err := pool.Exec(ctx, statements); err != nil {
		return fmt.Errorf("execute fixture %q: %w", name, err)
	}

	return nil
}

func loadFixture(filename string) ([]byte, error) {
	bases, paths := candidateFixtureBases()

	for _, base := range bases {
		path := filepath.Join(base, fixtureDir, filename)
		data, err := os.ReadFile(path)
		if err == nil {
			return data, nil
		}

		if errors.Is(err, os.ErrNotExist) {
			continue
		}

		return nil, fmt.Errorf("read fixture %q: %w", filename, err)
	}

	return nil, fmt.Errorf(defaultErrorMsg, filename, strings.Join(paths, ", "))
}

func candidateFixtureBases() ([]string, []string) {
	var (
		seen  = make(map[string]struct{})
		bases []string
		paths []string
	)

	add := func(dir string) {
		if dir == "" {
			return
		}
		if _, ok := seen[dir]; ok {
			return
		}
		seen[dir] = struct{}{}
		bases = append(bases, dir)
		paths = append(paths, filepath.Join(dir, fixtureDir))
	}

	if wd, err := os.Getwd(); err == nil {
		if base := findFixtureBase(wd); base != "" {
			add(base)
		}
	}

	if exeDir, err := executableDirectory(); err == nil {
		if base := findFixtureBase(exeDir); base != "" {
			add(base)
		}
	}

	return bases, paths
}

func findFixtureBase(start string) string {
	dir := start
	for {
		candidate := filepath.Join(dir, fixtureDir)
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}

func executableDirectory() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exe), nil
}
