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
	FileName   string
	OriginDir  string
	Language   config.LanguageCode
	Content    Markdown
	Translated Markdown
}

type MarkdownFiles []MarkdownFile

type Writer interface {
	Write(ctx context.Context, file MarkdownFile) error
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

func (w writer) Write(ctx context.Context, file MarkdownFile) error {
	targetPath := TargetFileContentPath(w.cfg.ContentDir, w.cfg.TargetPathRule, file.OriginDir, file.Language.String(), file.FileName)
	slog.DebugContext(ctx, "output path for translated markdown", "path", targetPath)

	parent := filepath.Dir(targetPath)

	if err := os.MkdirAll(parent, os.ModePerm); err != nil {
		return errors.Wrap(err, "failed to create parent directory")
	}

	if err := WriteMarkdownWithFrontmatter(targetPath, []byte(file.Translated), os.ModePerm,
		"translated", true,
	); err != nil {
		return err
	}

	return nil
}
