package util

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

//Checks if 3 tables need for execution are initialised properly
func IsStateDBInitialised(db *sql.DB) bool {
	required := map[string]bool{
		"files":    false,
		"remotes":  false,
		"metadata": false,
	}

	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		return false
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return false
		}
		if _, ok := required[name]; ok {
			required[name] = true
		}
	}

	for _, ok := range required {
		if !ok {
			return false
		}
	}

	return true
}


func IsStateDBInitialisedFromDirPath(path string) (bool, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return false, err
	}
	if dbExists:=PathExists(filepath.Join(path,"state.db")); !dbExists {
		return false, os.ErrNotExist
	}
	dbPath := filepath.Join(path, "state.db")
	p := filepath.ToSlash(dbPath)
	dsn := fmt.Sprintf("file:%s?mode=ro&_busy_timeout=5000", p)
	db, err := sql.Open("sqlite",dsn)
	if err != nil {
		return false, err
	}
	if err:= db.Ping(); err!=nil{
		return false, err
	}
	defer db.Close()

	return IsStateDBInitialised(db), nil
}