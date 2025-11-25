package util

import (
	"errors"
	"os"
	"path/filepath"
)

var(
	ErrPathAccess = errors.New("error accessing path")
	ErrNotDir = errors.New("is not a directory")
)

// Checks if path exists, only non erronous access to path stat, results in true
func PathExists(path string) bool {
	path, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	_, err = os.Stat(path)
	return err==nil
}


// Checks validity as a directory, true if path exists, and is a directory, false in every other case
func IsValidDir(path string) (error) {

	path, err := filepath.Abs(path)
	if err != nil {
		return ErrPathAccess
	}
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return ErrNotDir
	}
	return nil
}