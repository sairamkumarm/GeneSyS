package util

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestIsValidGeneSySFolder(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(base string) string
		wantBool bool
		wantErr  error
	}{
		{
			name: "BasePathDoesNotExist",
			setup: func(base string) string {
				return filepath.Join(base, "does_not_exist")
			},
			wantBool: false,
			wantErr:  os.ErrNotExist,
		},
		{
			name: "BaseExistsButGenesysMissing",
			setup: func(base string) string {
				return base
			},
			wantBool: false,
			wantErr:  os.ErrNotExist,
		},
		{
			name: "GenesysIsFileNotDir",
			setup: func(base string) string {
				p := filepath.Join(base, ".genesys")
				if err := os.WriteFile(p, []byte("not a dir"), 0644); err != nil {
					t.Fatalf("setup error: %v", err)
				}
				return base
			},
			wantBool: false,
			wantErr:  ErrNotDir,
		},
		{
			name: "ValidGenesysFolder",
			setup: func(base string) string {
				p := filepath.Join(base, ".genesys")
				if err := os.Mkdir(p, 0755); err != nil {
					t.Fatalf("setup error: %v", err)
				}
				return base
			},
			wantBool: true,
			wantErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base := t.TempDir()
			path := tt.setup(base)

			gotBool, gotErr := IsValidGeneSySFolder(path)

			if gotBool != tt.wantBool {
				t.Errorf("expected bool %v, got %v", tt.wantBool, gotBool)
			}

			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("expected error %v, got %v", tt.wantErr, gotErr)
			}
		})
	}
}
func TestIsInitialisedGeneSySRepo(t *testing.T) {

	t.Run("InvalidBaseFolder", func(t *testing.T) {
		_, err := IsInitialisedGeneSySRepo("does/not/exist")
		if !errors.Is(err, os.ErrNotExist) {
			t.Errorf("expected os.ErrNotExist, got %v", err)
		}
	})

	t.Run("ValidFolderButNoDB", func(t *testing.T) {
		base := t.TempDir()
		genesys := filepath.Join(base, ".genesys")
		if err := os.Mkdir(genesys, 0755); err != nil {
			t.Fatalf("setup error: %v", err)
		}

		_, err := IsInitialisedGeneSySRepo(base)
		if !errors.Is(err, os.ErrNotExist) {
			t.Errorf("expected os.ErrNotExist, got %v", err)
		}
	})

	t.Run("DBExistsButNotInitialised", func(t *testing.T) {
		base := t.TempDir()

		genesys := filepath.Join(base, ".genesys")
		if err := os.Mkdir(genesys, 0755); err != nil {
			t.Fatalf("setup error: %v", err)
		}

		// dbPath := filepath.Join(genesys, "state.db")
		// db, err := sql.Open("sqlite", dbPath)
		// if err != nil {
		// 	t.Fatalf("setup error: %v", err)
		// }
		// if _, err := db.Exec(`CREATE TABLE placeholder (id TEXT)`); err != nil {
		// 	t.Fatalf("setup error: %v", err)
		// }
		// db.Close()
		f, err :=os.Create(filepath.Join(genesys,"state.db"))
		if err != nil {
			t.Fatalf("setup error: %v", err)
		}
		f.Close()
		ok, err := IsInitialisedGeneSySRepo(base)
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if ok {
			t.Fatalf("expected false, got true")
		}
	})

	t.Run("DBNonExistant", func(t *testing.T) {
		base := t.TempDir()

		genesys := filepath.Join(base, ".genesys")
		if err := os.Mkdir(genesys, 0755); err != nil {
			t.Fatalf("setup error: %v", err)
		}

		ok, err := IsInitialisedGeneSySRepo(base)
		if err == nil {
			t.Fatalf("expected error, got %v", err)
		}
		if ok {
			t.Fatalf("expected false, got true")
		}
	})

	t.Run("FullyInitialisedRepo", func(t *testing.T) {
		base := t.TempDir()

		genesys := filepath.Join(base, ".genesys")
		if err := os.Mkdir(genesys, 0755); err != nil {
			t.Fatalf("setup error: %v", err)
		}

		dbPath := filepath.Join(genesys, "state.db")
		db, err := sql.Open("sqlite", dbPath)
		if err != nil {
			t.Fatalf("setup error: %v", err)
		}

		_, err = db.Exec(`
			CREATE TABLE files (id TEXT);
			CREATE TABLE remotes (id TEXT);
			CREATE TABLE metadata (id TEXT);
		`)
		if err != nil {
			t.Fatalf("setup error: %v", err)
		}
		db.Close()

		ok, err := IsInitialisedGeneSySRepo(base)
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if !ok {
			t.Fatalf("expected true, got false")
		}
	})
}
