package fixtures

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

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

	raw := strings.TrimSpace(string(content))
	if raw == "" {
		return nil
	}

	statements, err := splitStatements(raw)
	if err != nil {
		return fmt.Errorf("split fixture %q: %w", name, err)
	}

	for idx, stmt := range statements {
		if stmt == "" {
			continue
		}
		if _, err := pool.Exec(ctx, stmt); err != nil {
			return fmt.Errorf("execute fixture %q statement %d: %w", name, idx+1, err)
		}
	}

	return nil
}

func splitStatements(sql string) ([]string, error) {
	var (
		statements     []string
		current        strings.Builder
		inSingleQuote  bool
		inDoubleQuote  bool
		inLineComment  bool
		inBlockComment bool
		dollarTag      string
	)

	flush := func() {
		stmt := strings.TrimSpace(current.String())
		if stmt != "" {
			statements = append(statements, stmt)
		}
		current.Reset()
	}

	readDollarTag := func(src string, start int) (string, int) {
		end := start + 1
		for end < len(src) {
			ch := src[end]
			if ch == '$' {
				return src[start : end+1], end
			}
			if ch == '_' || unicode.IsLetter(rune(ch)) || unicode.IsDigit(rune(ch)) {
				end++
				continue
			}
			break
		}
		return "", start
	}

	for i := 0; i < len(sql); i++ {
		ch := sql[i]
		var next byte
		if i+1 < len(sql) {
			next = sql[i+1]
		}

		if inLineComment {
			if ch == '\n' {
				inLineComment = false
				current.WriteByte(ch)
			}
			continue
		}

		if inBlockComment {
			if ch == '*' && next == '/' {
				inBlockComment = false
				i++
			}
			continue
		}

		if dollarTag != "" {
			if ch == '$' && len(sql)-i >= len(dollarTag) && sql[i:i+len(dollarTag)] == dollarTag {
				current.WriteString(dollarTag)
				i += len(dollarTag) - 1
				dollarTag = ""
				continue
			}
			current.WriteByte(ch)
			continue
		}

		if !inSingleQuote && !inDoubleQuote {
			if ch == '-' && next == '-' {
				inLineComment = true
				i++
				continue
			}
			if ch == '/' && next == '*' {
				inBlockComment = true
				i++
				continue
			}
			if ch == '$' {
				tag, end := readDollarTag(sql, i)
				if tag != "" {
					dollarTag = tag
					current.WriteString(tag)
					i = end
					continue
				}
			}
		}

		switch {
		case inSingleQuote:
			current.WriteByte(ch)
			if ch == '\'' {
				if next == '\'' {
					current.WriteByte(next)
					i++
				} else {
					inSingleQuote = false
				}
			}
		case inDoubleQuote:
			current.WriteByte(ch)
			if ch == '"' {
				if next == '"' {
					current.WriteByte(next)
					i++
				} else {
					inDoubleQuote = false
				}
			}
		default:
			if ch == '\'' {
				inSingleQuote = true
				current.WriteByte(ch)
				continue
			}
			if ch == '"' {
				inDoubleQuote = true
				current.WriteByte(ch)
				continue
			}
			if ch == ';' {
				flush()
				continue
			}
			current.WriteByte(ch)
		}
	}

	if inSingleQuote {
		return nil, errors.New("unterminated single-quoted string literal")
	}
	if inDoubleQuote {
		return nil, errors.New("unterminated double-quoted identifier")
	}
	if dollarTag != "" {
		return nil, fmt.Errorf("unterminated dollar-quoted literal %s", dollarTag)
	}
	if inBlockComment {
		return nil, errors.New("unterminated block comment")
	}

	flush()
	return statements, nil
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
