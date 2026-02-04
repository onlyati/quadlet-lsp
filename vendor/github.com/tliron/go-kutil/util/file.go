package util

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func DoesFileExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {
		return false, err
	}
}

func IsFileHidden(path string) bool {
	return strings.HasPrefix(filepath.Base(path), ".")
}

func IsFileExecutable(mode fs.FileMode) bool {
	return mode&0100 != 0
}

func Touch(path string, permissions fs.FileMode, dirPermissions fs.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), dirPermissions); err == nil {
		if file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, permissions); err == nil {
			return file.Close()
		} else {
			return err
		}
	} else {
		return err
	}
}
