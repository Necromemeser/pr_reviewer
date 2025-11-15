package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	_ "github.com/lib/pq"
)

func ConnectAndMigrate() (*sql.DB, error) {
	dsn := "postgres://user:pass@db:5432/pr_reviewer?sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	var pingErr error
	for i := 0; i < 30; i++ {
		pingErr = db.Ping()
		if pingErr == nil {
			break
		}
		time.Sleep(time.Second)
	}
	if pingErr != nil {
		return nil, fmt.Errorf("ping db: %w", pingErr)
	}

	log.Println("Connected to database")

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
