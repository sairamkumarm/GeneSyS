package database

import (
	"database/sql"
	"genesys/internal/util"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func AttemptInitSQLite(metaDataFolderPath string) (int,error) {
	dbPath := filepath.Join(metaDataFolderPath, "state.db")
	db, err := sql.Open("sqlite",dbPath)
	if err != nil {
		return 1,err
	}
	if err:= db.Ping(); err!=nil{
		return 1,err
	}
	defer db.Close()

	if util.IsStateDBInitialised(db) {
		return 2,nil
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
		return 1,err
	}
	return 0,nil
}
