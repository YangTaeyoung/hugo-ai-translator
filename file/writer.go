package file

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"

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
	Write(ctx context.Context, files MarkdownFiles) error
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

func (w writer) Write(ctx context.Context, files MarkdownFiles) error {
	var (
		err error
	)

	for i, file := range files {
		replacer := strings.NewReplacer(
			"{origin}", file.OriginDir,
			"{language}", string(file.Language),
			"{fileName}", file.FileName,
		)

		targetPath := path.Join(w.cfg.ContentDir, replacer.Replace(w.cfg.TargetPathRule))
		slog.DebugContext(ctx, "output path for translated markdown", "path", targetPath)

		parent := filepath.Dir(targetPath)

		if err = os.MkdirAll(parent, os.ModePerm); err != nil {
			return errors.Wrap(err, "failed to create parent directory")
		}

		if err = os.WriteFile(targetPath, []byte(file.Content), os.ModePerm); err != nil {
			return errors.Wrap(err, "failed to write translated markdown file")
		}

		slog.InfoContext(ctx, "file written",
			"path", targetPath,
			"progress", fmt.Sprintf("%d/%d", i+1, len(files)),
		)
	}

	return nil
}
