package utils

import (
	"path/filepath"
	"strings"
)

func ChangeExt(filename, newExt string) string {
	base := filepath.Base(filename)
	ext := filepath.Ext(base)
	return strings.TrimSuffix(filename, ext) + newExt
}
