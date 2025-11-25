package util

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func IsValidGeneSySFolder(path string) (bool, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false, errors.New(fmt.Sprint("error access path:", path))
	}
	exists:= PathExists(absPath)
	if !exists {
		return false, os.ErrNotExist
	} 
	dirPath:= filepath.Join(absPath, ".genesys")
	exists = PathExists(dirPath)
	if !exists {
		return false, os.ErrNotExist
	}
	err = IsValidDir(dirPath)
	return err==nil,err
}

func IsInitialisedGeneSySRepo(path string) (bool, error) {
	folderExists, err := IsValidGeneSySFolder(path)
	if err != nil {
		return false, err
	}

	dbReady, err := IsStateDBInitialisedFromDirPath(filepath.Join(path,".genesys"))
	if err != nil {
		return false, err
	}
	return folderExists && dbReady, nil
}