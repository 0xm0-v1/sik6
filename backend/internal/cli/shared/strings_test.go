package shared

import "testing"

func TestSanitizeList(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input []string
		want  []string
	}{
		"empty":         {input: nil, want: nil},
		"spaces-only":   {input: []string{"   "}, want: nil},
		"trimmed":       {input: []string{" foo ", "bar"}, want: []string{"foo", "bar"}},
		"drop-empty":    {input: []string{"foo", "", "bar"}, want: []string{"foo", "bar"}},
		"tabs":          {input: []string{"\tfoo\t", "\tbar\t"}, want: []string{"foo", "bar"}},
		"already-clean": {input: []string{"foo", "bar"}, want: []string{"foo", "bar"}},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := SanitizeList(tc.input)
			if len(got) != len(tc.want) {
				t.Fatalf("expected %d items, got %d (%v)", len(tc.want), len(got), got)
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Fatalf("index %d: expected %q, got %q", i, tc.want[i], got[i])
				}
			}
		})
	}
}
