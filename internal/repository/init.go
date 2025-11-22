package repository

import (
	"errors"
	"fmt"
	"genesys/internal/database"
	"os"
	"path/filepath"
)



func InitOrLoadRepo(path string) error {
	if path=="" {
		path="."
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	folderPath := filepath.Join(abs,".genesys")
	info, err := os.Stat(folderPath)
	if err != nil {
		if	errors.Is(err, os.ErrNotExist){
			if err := os.MkdirAll(folderPath, 0755); err!=nil {
				return err
			}
		} else {
			return err
		}
	} else if !info.IsDir(){
		return errors.New( fmt.Sprint(folderPath,"is not a directory, metadata cannot be stored when .genesys is not a folder"))
	}

	err = database.AttempInitSQLLite(folderPath)
	if err != nil {
		return err
	}
	return nil
}