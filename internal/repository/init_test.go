package repository_test

import (
	"os"
	"path/filepath"
	"testing"

	"genesys/internal/repository"
	"genesys/internal/util"
)

// helper to check folder existence
func folderExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func TestInitOrLoadRepo(t *testing.T) {

	t.Run("EmptyPathInitializesCurrentDir", func(t *testing.T) {
		tempDir := t.TempDir()
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tempDir)

		status, err := repository.InitOrLoadRepo("")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if status != 0 {
			t.Fatalf("expected status 0, got %d", status)
		}

		genesysPath := filepath.Join(tempDir, ".genesys")
		if !folderExists(genesysPath) {
			t.Fatalf(".genesys folder not created")
		}
	})

	t.Run("NonExistentPathCreatesFolderAndDB", func(t *testing.T) {
		tempDir := t.TempDir()
		targetPath := filepath.Join(tempDir, "nested")
		status, err := repository.InitOrLoadRepo(targetPath)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if status != 0 {
			t.Fatalf("expected status 0, got %d", status)
		}

		genesysPath := filepath.Join(targetPath, ".genesys")
		if !folderExists(genesysPath) {
			t.Fatalf(".genesys folder not created")
		}

		// check state.db created
		dbPath := filepath.Join(genesysPath, "state.db")
		if !util.PathExists(dbPath) {
			t.Fatalf("state.db not created")
		}
	})

	t.Run("AlreadyInitializedRepo", func(t *testing.T) {
		tempDir := t.TempDir()
		genesysPath := filepath.Join(tempDir, ".genesys")
		os.MkdirAll(genesysPath, 0755)

		// first init
		status, err := repository.InitOrLoadRepo(tempDir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if status != 0 {
			t.Fatalf("expected status 0, got %d", status)
		}

		// second init should detect existing DB
		status, err = repository.InitOrLoadRepo(tempDir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if status != 2 {
			t.Fatalf("expected status 2, got %d", status)
		}
	})

	t.Run("InvalidPathReturnsError", func(t *testing.T) {
		// Use a path that cannot be accessed, e.g., invalid characters
		invalidPath := string([]byte{0x00})
		status, err := repository.InitOrLoadRepo(invalidPath)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if status != 1 {
			t.Fatalf("expected status 1, got %d", status)
		}
	})
	t.Run("InitInsideFileReturnsError", func(t *testing.T) {
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "not_a_dir")

		// create a file
		f, err := os.Create(filePath)
		if err != nil {
			t.Fatalf("could not create temp file: %v", err)
		}
		f.Close()

		status, err := repository.InitOrLoadRepo(filePath)
		if err == nil {
			t.Fatalf("expected error when path is a file, got nil")
		}
		if status != 1 {
			t.Fatalf("expected status 1, got %d", status)
		}
	})

}
