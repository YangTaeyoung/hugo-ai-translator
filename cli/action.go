package cli

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/YangTaeyoung/hugo-ai-translator/config"
	"github.com/YangTaeyoung/hugo-ai-translator/environment"
	"github.com/YangTaeyoung/hugo-ai-translator/file"
	"github.com/YangTaeyoung/hugo-ai-translator/translator"
	"github.com/urfave/cli/v3"
)

func TranslateAction(ctx context.Context, cmd *cli.Command) error {
	cfgPath := cmd.String("config")

	cfg, err := config.New(cfgPath)
	if err != nil {
		return err
	}
	slog.DebugContext(ctx, "config parsed", "path", cfgPath)

	env := environment.New(cfg)
	slog.DebugContext(ctx, "environment created")

	p := file.NewParser(file.ParserConfig{
		ContentDir:  cfg.Translator.ContentDir,
		IgnoreRules: cfg.Translator.Source.IgnoreRules,
	})

	markdownFiles, err := p.Parse(ctx)
	if err != nil {
		return err
	}
	slog.DebugContext(ctx, "markdown files parsed", "count", len(markdownFiles))

	t := translator.New(env.Client, translator.Config{
		SourceLanguage:  cfg.Translator.Source.SourceLanguage,
		TargetLanguages: cfg.Translator.Target.TargetLanguages,
	})

	w := file.NewWriter(file.WriterConfig{
		ContentDir:     cfg.Translator.ContentDir,
		TargetPathRule: cfg.Translator.Target.TargetPathRule,
	})

	for i, markdownFile := range markdownFiles {
		var translated translator.Results

		translated, err = t.Translate(ctx, markdownFile)
		if err != nil {
			return err
		}

		if err = w.Write(ctx, translated.MarkdownFiles()); err != nil {
			return err
		}

		slog.InfoContext(ctx, "translated",
			"path", markdownFile.Path,
			"progress", fmt.Sprintf("%d/%d", i+1, len(markdownFiles)),
		)
	}

	return nil
}
