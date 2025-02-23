package translator

import (
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

var (
	ErrEmptyPath = errors.New("empty path")
)

func fileNameWithoutExtension(path string) (string, error) {
	if path == "" {
		return "", ErrEmptyPath
	}

	fileName := filepath.Base(path)

	if pos := strings.LastIndexByte(fileName, '.'); pos != -1 {
		return fileName[:pos], nil
	}

	return fileName, nil
}
