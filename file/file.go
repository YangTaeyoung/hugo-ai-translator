package file

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

var (
	ErrEmptyPath = errors.New("empty path")
)

func TargetFilePath(contentDir string,
	targetFilePathRule string,
	origin string,
	language string,
	fileName string,
) string {
	replacer := strings.NewReplacer(
		"{origin}", origin,
		"{language}", language,
		"{fileName}", fileName,
	)

	return path.Join(contentDir, replacer.Replace(targetFilePathRule))
}

func FileNameWithoutExtension(path string) (string, error) {
	if path == "" {
		return "", ErrEmptyPath
	}

	fileName := filepath.Base(path)

	if pos := strings.LastIndexByte(fileName, '.'); pos != -1 {
		return fileName[:pos], nil
	}

	return fileName, nil
}
