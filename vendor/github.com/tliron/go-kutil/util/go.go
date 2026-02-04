package util

import (
	"os"
	"path/filepath"
)

func GetGoPath() (string, error) {
	if gopath := os.Getenv("GOPATH"); gopath != "" {
		return gopath, nil
	} else {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, "go"), nil
		} else {
			return "", err
		}
	}
}

func GetGoBin() (string, error) {
	if gobin := os.Getenv("GOBIN"); gobin != "" {
		return gobin, nil
	} else if gopath, err := GetGoPath(); err == nil {
		return filepath.Join(gopath, "bin"), nil
	} else {
		return "", err
	}
}
