// common - generate utiles
package common

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// CreatePathWithFile - create a file and all directories to it if they do not exist
func CreatePathWithFile(partOfFilePath string) error {
	fileName := filepath.Base(partOfFilePath)
	if fileName == "" {
		return errors.New("incorrect path of file")
	}
	if fileExtension := filepath.Ext(fileName); fileExtension != ".db" {
		return errors.New("invalid file extension")
	}
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	fullPath := filepath.Join(currentDir, partOfFilePath)
	onlyDir := strings.Replace(fullPath, fileName, "", 1)
	if err := os.MkdirAll(onlyDir, 0o755); err != nil {
		return err
	}
	_, err = os.Create(fullPath)
	return err
}

// Abs - absolute value
func Abs(val int) int {
	if val < 0 {
		return -val
	}
	return val
}
