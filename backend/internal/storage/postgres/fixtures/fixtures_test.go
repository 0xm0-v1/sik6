package fixtures

import (
	"reflect"
	"testing"
)

func TestSplitStatements(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		sql     string
		want    []string
		wantErr bool
	}{
		{
			name: "basic multi statement",
			sql: `-- comment
TRUNCATE TABLE recipes RESTART IDENTITY;
/* block comment */ INSERT INTO recipes (recipe_name) VALUES ('Classic Pancakes');
`,
			want: []string{
				"TRUNCATE TABLE recipes RESTART IDENTITY",
				"INSERT INTO recipes (recipe_name) VALUES ('Classic Pancakes')",
			},
		},
		{
			name: "semicolon inside literals",
			sql: `INSERT INTO sample(text) VALUES ('keep; inside');
INSERT INTO sample(text) VALUES ("double quotes; ok");
`,
			want: []string{
				"INSERT INTO sample(text) VALUES ('keep; inside')",
				`INSERT INTO sample(text) VALUES ("double quotes; ok")`,
			},
		},
		{
			name: "dollar quoted block",
			sql: `DO $$BEGIN
  PERFORM 1;
END$$;
TRUNCATE table demo;
`,
			want: []string{
				"DO $$BEGIN\n  PERFORM 1;\nEND$$",
				"TRUNCATE table demo",
			},
		},
		{
			name:    "unterminated string",
			sql:     `INSERT INTO sample(text) VALUES ('missing terminator);`,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := splitStatements(tc.sql)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error but got nil (statements: %v)", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("splitStatements returned error: %v", err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("unexpected statements:\nwant: %#v\n got: %#v", tc.want, got)
			}
		})
	}
}
