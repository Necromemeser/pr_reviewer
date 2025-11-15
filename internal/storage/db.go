package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	_ "github.com/lib/pq"
)

func ConnectAndMigrate() (*sql.DB, error) {
	dsn := "postgres://user:pass@db:5432/pr_reviewer?sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	log.Println("Connected to database")

	// sqlBytes, err := os.ReadFile("migrations/001_init.sql")
	// if err != nil {
	// 	return nil, fmt.Errorf("read migration file: %w", err)
	// }

	// _, err = db.Exec(string(sqlBytes))
	// if err != nil {
	// 	return nil, fmt.Errorf("apply migration: %w", err)
	// }

	// log.Println("Migration applied")

	dir := "migrations"

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read migrations dir: %w", err)
	}

	var migrationFiles []string
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".sql" {
			migrationFiles = append(migrationFiles, f.Name())
		}
	}

	sort.Strings(migrationFiles)

	for _, file := range migrationFiles {
		path := filepath.Join(dir, file)

		sqlBytes, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read migration file %s: %w", file, err)
		}

		_, err = db.Exec(string(sqlBytes))
		if err != nil {
			return nil, fmt.Errorf("apply migration %s: %w", file, err)
		}

		fmt.Println("Applied migration:", file)
	}

	return db, nil
}
