package util

import (
	"errors"
	"os"
	"testing"
)

func TestPathExists(t *testing.T) {
	t.Run("ExistingFileReturnsTrue", func(t *testing.T) {
		f, err := os.CreateTemp("", "existing_path")
		if err != nil {
			t.Fatalf("temp file creation failed: %v", err)
		}
		path := f.Name()
		f.Close()

		if !PathExists(path) {
			t.Errorf("expected true, got false")
		}
	})
	t.Run("NonExistingFileReturnsFalse", func(t *testing.T) {

		if PathExists("randompaththatshouldnotexist") {
			t.Errorf("expected false, got true")
		}
	})
	t.Run("InvalidPathNameReturnsFalse", func(t *testing.T) {

		if PathExists("randompaththat / shouldnotexist") {
			t.Errorf("expected false, got true")
		}
	})
	t.Run("PathReturnsTrue", func(t *testing.T) {
		f, err := os.CreateTemp("", "existing path")
		if err != nil {
			t.Fatalf("temp file creation failed: %v", err)
		}
		path := f.Name()
		f.Close()

		if !PathExists(path) {
			t.Errorf("expected true, got false")
		}
	})
}

func TestIsValidDir(t *testing.T) {

	tempFile, err:=os.CreateTemp("","temp_file")
	if err != nil {
		t.Fatalf("Error creating test file")
	}
	tempFilePath:=tempFile.Name()
	tempFile.Close()

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		path    string
		err error
	}{
		// TODO: Add test cases.
		{"ExistingDirectoryReturnsTrue",t.TempDir(), nil},
		{"NonExistantPathReturnsError","randompath", os.ErrNotExist},
		{"ExistingPathNotDirectoryReturnsError",tempFilePath, ErrNotDir},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := IsValidDir(tt.path)
			if !errors.Is(gotErr, tt.err) {
				t.Errorf("IsValidDir() failed: recieved '%v', expected '%v'", gotErr, tt.err)
			}
		})
	}
}
