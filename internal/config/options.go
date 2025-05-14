// options - contain property of file for config
package config

import (
	"errors"
	"os"
	"path/filepath"
)

var (
	// ErrOptionsExtansion - no extansion of file
	ErrOptionsExtansion = errors.New("extansion file empty")

	ErrOptionsEmptyFile = errors.New("empty path of file")
)

type options struct {
	pathOfFile string

	// file extansion without dot
	fileExt string

	// file name without extension
	fileName string
}

// parsePath - member of options
// parse file after get 'pathOfFile'
func (op *options) parsePath() error {
	if op.pathOfFile == "" {
		return ErrOptionsEmptyFile
	}
	if _, err := os.Stat(op.pathOfFile); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// find data from ENV
			op.pathOfFile = ""
			return nil
		}
		return err
	}
	op.fileName = filepath.Base(op.pathOfFile)
	op.fileExt = filepath.Ext(op.fileName)
	if op.fileExt == "" {
		return ErrOptionsExtansion
	}
	op.fileName = op.fileName[len(op.fileExt):]
	if op.fileExt[0] == '.' {
		op.fileExt = op.fileExt[1:]
	}
	return nil
}
