package repository

import (
	"errors"
	"fmt"
	"genesys/internal/database"
	"genesys/internal/util"
	"os"
	"path/filepath"
)
 
func InitOrLoadRepo(path string) (int, error) {
	if path=="" {
		path="."
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return 1, err
	}
	folderPath:=filepath.Join(absPath,".genesys")
	exists, err := util.IsValidGeneSySFolder(folderPath)
	if err != nil {
		if	errors.Is(err, os.ErrNotExist){
			if err := os.MkdirAll(folderPath, 0755); err!=nil {
				return 1, errors.New("error creating .genesys/ folder")
			}
		} else {
			return 1, errors.New(fmt.Sprint(err, folderPath))	
		}
	} else if !exists {
		return 1, errors.New("invalid state: .genesys is nonexistant, without validation error")
	}
	status, err := database.AttemptInitSQLite(folderPath)
	if err != nil {
		return 1, err
	}
	return status, nil
}