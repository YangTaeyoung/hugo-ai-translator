package file

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/YangTaeyoung/hugo-ai-translator/config"
	"github.com/pkg/errors"
)

type MarkdownFile struct {
	FileName  string
	OriginDir string
	Language  config.LanguageCode
	Content   Markdown
}

type MarkdownFiles []MarkdownFile

type Writer interface {
	Write(ctx context.Context, files MarkdownFiles) ([]string, error)
}

type WriterConfig struct {
	ContentDir     string
	TargetPathRule string
}

type writer struct {
	cfg WriterConfig
}

func NewWriter(cfg WriterConfig) Writer {
	return &writer{
		cfg: cfg,
	}
}

func (w writer) Write(ctx context.Context, files MarkdownFiles) ([]string, error) {
	var (
		paths []string
		err   error
	)

	for _, file := range files {
		targetPath := TargetFilePath(w.cfg.ContentDir, w.cfg.TargetPathRule, file.OriginDir, string(file.Language), file.FileName)
		slog.DebugContext(ctx, "output path for translated markdown", "path", targetPath)

		parent := filepath.Dir(targetPath)

		if err = os.MkdirAll(parent, os.ModePerm); err != nil {
			return nil, errors.Wrap(err, "failed to create parent directory")
		}

		if err = os.WriteFile(targetPath, []byte(file.Content), os.ModePerm); err != nil {
			return nil, errors.Wrap(err, "failed to write translated markdown file")
		}

		paths = append(paths, targetPath)
	}

	return paths, nil
}
