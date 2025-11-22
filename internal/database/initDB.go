package database

import (
	"database/sql"
	"fmt"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func AttempInitSQLLite(metaDataFolderPath string) error {
	dbPath := filepath.Join(metaDataFolderPath, "state.db")
	db, err := sql.Open("sqlite",dbPath)
	if err != nil {
		return err
	}
	
	if err:= db.Ping(); err!=nil{
		return err
	}
	defer db.Close()

	if isInitalised(db) {
		return fmt.Errorf("already initialised repository")
	}

	schema:= `
		CREATE TABLE IF NOT EXISTS files (
			uuid TEXT PRIMARY KEY,
			name TEXT,
			relative_path TEXT,
			last_modified INTEGER,
			hash TEXT
		);
		CREATE TABLE IF NOT EXISTS remotes (
			name TEXT,
			link TEXT
		);
		CREATE TABLE IF NOT EXISTS metadata (
			key TEXT,
			value TEXT
		);
	`

	_, err = db.Exec(schema)
	if err != nil {
		return err
	}
	return nil
}

func isInitalised(db *sql.DB) bool{
	var name string
	err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='files'").Scan(&name)
	if err == sql.ErrNoRows {
        return false
    }
    if err != nil {
        return false
    }
    return true
}