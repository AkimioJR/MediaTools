package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func ChangeExt(filename, newExt string) string {
	base := filepath.Base(filename)
	ext := filepath.Ext(base)
	return strings.TrimSuffix(filename, ext) + newExt
}

func CreateFile(filePath string) (*os.File, error) {
	err := os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		return nil, err
	}
	return os.Create(filePath)
}
