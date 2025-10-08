package shared

import "strings"

// SanitizeList trims whitespace from each entry and drops empty values.
func SanitizeList(values []string) []string {
	out := make([]string, 0, len(values))
	for _, v := range values {
		if trimmed := strings.TrimSpace(v); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}
