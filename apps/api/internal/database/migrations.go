package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type migrationFile struct {
	name    string
	version string
	content string
}

// RunMigrations applies SQL migration files in version order and records each
// successful application in schema_migrations.
func RunMigrations(ctx context.Context, db *sql.DB, dir string) error {
	if db == nil {
		return errors.New("database is nil")
	}

	if _, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`); err != nil {
		return fmt.Errorf("ensure schema_migrations: %w", err)
	}

	applied, err := appliedMigrations(ctx, db)
	if err != nil {
		return err
	}

	files, err := readMigrationFiles(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if _, ok := applied[file.version]; ok {
			continue
		}

		if err := applyMigration(ctx, db, file); err != nil {
			return err
		}
	}

	return nil
}

func appliedMigrations(ctx context.Context, db *sql.DB) (map[string]struct{}, error) {
	rows, err := db.QueryContext(ctx, `SELECT version FROM schema_migrations`)
	if err != nil {
		return nil, fmt.Errorf("load applied migrations: %w", err)
	}
	defer rows.Close()

	applied := make(map[string]struct{})
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, fmt.Errorf("scan applied migration: %w", err)
		}
		applied[version] = struct{}{}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate applied migrations: %w", err)
	}

	return applied, nil
}

func readMigrationFiles(dir string) ([]migrationFile, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read migrations directory %q: %w", dir, err)
	}

	files := make([]migrationFile, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".sql" {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read migration %q: %w", path, err)
		}

		files = append(files, migrationFile{
			name:    entry.Name(),
			version: strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name())),
			content: string(content),
		})
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].name < files[j].name
	})

	return files, nil
}

func applyMigration(ctx context.Context, db *sql.DB, file migrationFile) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin migration %q: %w", file.name, err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	statements := splitSQLStatements(file.content)
	for _, statement := range statements {
		if _, err := tx.ExecContext(ctx, statement); err != nil {
			return fmt.Errorf("apply migration %q: %w", file.name, err)
		}
	}

	if _, err := tx.ExecContext(ctx, `INSERT INTO schema_migrations (version) VALUES ($1)`, file.version); err != nil {
		return fmt.Errorf("record migration %q: %w", file.name, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit migration %q: %w", file.name, err)
	}

	return nil
}

func splitSQLStatements(input string) []string {
	statements := make([]string, 0)
	var current strings.Builder
	inSingleQuote := false

	for i := 0; i < len(input); i++ {
		ch := input[i]

		if ch == '\'' {
			current.WriteByte(ch)
			if inSingleQuote && i+1 < len(input) && input[i+1] == '\'' {
				current.WriteByte(input[i+1])
				i++
				continue
			}
			inSingleQuote = !inSingleQuote
			continue
		}

		if ch == ';' && !inSingleQuote {
			statement := strings.TrimSpace(current.String())
			if statement != "" {
				statements = append(statements, statement)
			}
			current.Reset()
			continue
		}

		current.WriteByte(ch)
	}

	if statement := strings.TrimSpace(current.String()); statement != "" {
		statements = append(statements, statement)
	}

	return statements
}
