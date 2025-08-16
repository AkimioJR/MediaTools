package utils

import (
	pathlib "path"
	"strings"
)

func ChangeExt(filename, newExt string) string {
	base := pathlib.Base(filename)
	ext := pathlib.Ext(base)
	return strings.TrimSuffix(filename, ext) + newExt
}
