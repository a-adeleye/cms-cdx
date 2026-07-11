package database

import "testing"

func TestSplitSQLStatements(t *testing.T) {
	t.Parallel()

	input := `
		CREATE TABLE example (id TEXT NOT NULL);
		INSERT INTO example (id) VALUES ('a;b');
		INSERT INTO example (id) VALUES ('It''s ready');
	`

	got := splitSQLStatements(input)
	want := []string{
		"CREATE TABLE example (id TEXT NOT NULL)",
		"INSERT INTO example (id) VALUES ('a;b')",
		"INSERT INTO example (id) VALUES ('It''s ready')",
	}

	if len(got) != len(want) {
		t.Fatalf("splitSQLStatements() len = %d, want %d", len(got), len(want))
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("splitSQLStatements()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
