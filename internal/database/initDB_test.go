package database

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"genesys/internal/util"

	_ "modernc.org/sqlite"
)

func TestAttemptInitSQLite(t *testing.T) {
	t.Run("FreshFolderCreatesDB", func(t *testing.T) {
		dir := t.TempDir()
		status, err := AttemptInitSQLite(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if status != 0 {
			t.Fatalf("expected status 0, got %d", status)
		}

		dbPath := filepath.Join(dir, "state.db")
		if !util.PathExists(dbPath) {
			t.Fatalf("state.db was not created")
		}
		time.Sleep(50 * time.Millisecond)
	})

	t.Run("AlreadyInitializedDB", func(t *testing.T) {
		dir := t.TempDir()
		_, err := AttemptInitSQLite(dir)
		if err != nil {
			t.Fatalf("setup failed: %v", err)
		}

		status, err := AttemptInitSQLite(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if status != 2 {
			t.Fatalf("expected status 2, got %d", status)
		}

	})

	t.Run("CorruptedFileOverwritten", func(t *testing.T) {
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
		status, err := AttemptInitSQLite(dir)
		if err != nil {
			t.Fatalf("corrupted files, should've been overwritten, it wasn't ")
		}
		if status == 1 {
			t.Fatalf("expected status 0, got %d", status)
		}

	})

	t.Run("MissingFolderReturnsError", func(t *testing.T) {
		dir := t.TempDir()
		nonExistentFolder := filepath.Join(dir, "nonexistentfolder")

		status, err := AttemptInitSQLite(nonExistentFolder) // pass folder, not file
		if err == nil {
			t.Fatalf("expected an error, got nil")
		}
		if status != 1 {
			t.Fatalf("expected status 1, got %d", status)
		}

	})
}
