package util

import (
	"path/filepath"
	"strings"
)

func PathWithoutExt(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}
