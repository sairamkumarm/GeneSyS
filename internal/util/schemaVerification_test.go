package util

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

func TestIsStateDBInitialised(t *testing.T) {

	t.Run("EmptyDBReturnsFalse", func(t *testing.T) {
		db, err := sql.Open("sqlite", ":memory:")
		if err != nil {
			t.Fatalf("failed to open in-memory DB: %v", err)
		}
		defer db.Close()

		if IsStateDBInitialised(db) {
			t.Errorf("expected false, got true for empty DB")
		}
	})

	t.Run("PartialDBReturnsFalse", func(t *testing.T) {
		db, err := sql.Open("sqlite", ":memory:")
		if err != nil {
			t.Fatalf("failed to open in-memory DB: %v", err)
		}
		defer db.Close()

		_, err = db.Exec(`CREATE TABLE files (id TEXT);`)
		if err != nil {
			t.Fatalf("failed to create table: %v", err)
		}

		if IsStateDBInitialised(db) {
			t.Errorf("expected false, got true for partial DB")
		}
	})

	t.Run("FullDBReturnsTrue", func(t *testing.T) {
		db, err := sql.Open("sqlite", ":memory:")
		if err != nil {
			t.Fatalf("failed to open in-memory DB: %v", err)
		}
		defer db.Close()

		_, err = db.Exec(`
			CREATE TABLE files (id TEXT);
			CREATE TABLE remotes (id TEXT);
			CREATE TABLE metadata (id TEXT);
		`)
		if err != nil {
			t.Fatalf("failed to create tables: %v", err)
		}

		if !IsStateDBInitialised(db) {
			t.Errorf("expected true, got false for full DB")
		}
	})
}


func TestIsStateDBInitialisedFromDirPath(t *testing.T) {

	t.Run("InvalidPathReturnsError", func(t *testing.T) {
		_, err := IsStateDBInitialisedFromDirPath("/this/does/not/exist")
		if err == nil {
			t.Fatalf("expected error for invalid path, got nil")
		}
	})

	t.Run("NoDBFileReturnsError", func(t *testing.T) {
		dir := t.TempDir()

		ok, err := IsStateDBInitialisedFromDirPath(dir)
		if err == nil {
			t.Fatalf("expected error for missing DB file, got nil")
		}
		if ok {
			t.Fatalf("expected false, got true")
		}
	})

	t.Run("CorruptDBReturnsError", func(t *testing.T) {
		dir := t.TempDir()
		f, err := os.CreateTemp(dir,"state.db")
		if err != nil {
			t.Fatalf("error creating temporary corrupted database file")
		}
		_, err = f.Write([]byte("corrupted file"))
		if err != nil {
			t.Fatalf("error writing in file")
		}
		f.Close()

		ok, err := IsStateDBInitialisedFromDirPath(dir)
		if err == nil {
			t.Fatalf("expected error for corrupt db, got nil")
		}
		if ok {
			t.Fatalf("expected false, got true")
		}
	})

	t.Run("EmptyDBReturnsFalse", func(t *testing.T) {
		dir := t.TempDir()
		dbPath := filepath.Join(dir, "state.db")

		db, err := sql.Open("sqlite", dbPath)
		if err != nil {
			t.Fatalf("failed to create empty DB: %v", err)
		}
		db.Close()

		ok, err := IsStateDBInitialisedFromDirPath(dir)
		if err == nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ok {
			t.Fatalf("expected false, got true")
		}
	})

	t.Run("PartialDBReturnsFalse", func(t *testing.T) {
		dir := t.TempDir()
		dbPath := filepath.Join(dir, "state.db")

		db, err := sql.Open("sqlite", dbPath)
		if err != nil {
			t.Fatalf("failed to open DB: %v", err)
		}

		_, err = db.Exec(`CREATE TABLE files (id TEXT);`)
		if err != nil {
			t.Fatalf("failed to create table: %v", err)
		}
		db.Close()

		ok, err := IsStateDBInitialisedFromDirPath(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ok {
			t.Fatalf("expected false for partial DB, got true")
		}
	})

	t.Run("FullDBReturnsTrue", func(t *testing.T) {
		dir := t.TempDir()
		dbPath := filepath.Join(dir, "state.db")

		db, err := sql.Open("sqlite", dbPath)
		if err != nil {
			t.Fatalf("failed to open DB: %v", err)
		}

		_, err = db.Exec(`
			CREATE TABLE files (id TEXT);
			CREATE TABLE remotes (id TEXT);
			CREATE TABLE metadata (id TEXT);
		`)
		if err != nil {
			t.Fatalf("failed to create tables: %v", err)
		}
		db.Close()

		ok, err := IsStateDBInitialisedFromDirPath(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !ok {
			t.Fatalf("expected true for full DB, got false")
		}
	})
}