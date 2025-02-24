package file

import (
	"context"
	"io/fs"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/YangTaeyoung/hugo-ai-translator/config"
	"github.com/bmatcuk/doublestar/v4"
)

type Markdown string

func (m Markdown) String() string {
	return string(m)
}

type ParsedMarkdownFile struct {
	Path            string
	Markdown        Markdown
	TargetLanguages config.LanguageCodes
}

type ParsedMarkdownFiles []ParsedMarkdownFile

type ParserConfig struct {
	ContentDir      string
	TranslatedPaths []string
	IgnoreRules     []string
	TargetLanguages config.LanguageCodes
	TargetPathRule  string
}

type Parser interface {
	Parse(ctx context.Context) (ParsedMarkdownFiles, error)
}

type parser struct {
	cfg ParserConfig
}

func NewParser(cfg ParserConfig) Parser {
	return &parser{
		cfg: cfg,
	}
}
func (p parser) listMarkdownFilePaths(ctx context.Context) ([]string, error) {
	var results []string

	if err := filepath.WalkDir(p.cfg.ContentDir, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// 디렉터리는 스킵
		if d.IsDir() {
			return nil
		}
		// 확장자가 .md가 아니면 무시
		if !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}
		// contentDir 기준의 상대 경로 계산

		absPath, err := filepath.Abs(filePath)
		if err != nil {
			return err
		}

		// 이미 번역된 내역인 경우 무시함
		if slices.Contains(p.cfg.TranslatedPaths, absPath) {
			slog.DebugContext(ctx, "skip already translated file ", "path", absPath)
			return nil
		}

		relPath, err := filepath.Rel(p.cfg.ContentDir, filePath)
		if err != nil {
			return err
		}

		// glob 패턴 매칭은 Unix 스타일 경로 구분자를 사용하는 것이 좋으므로 변환
		relPathUnix := filepath.ToSlash(relPath)
		// ignoreRules와 매칭되는지 확인
		for _, rule := range p.cfg.IgnoreRules {
			match, err := doublestar.PathMatch(rule, relPathUnix)
			if err != nil {
				return err
			}
			if match {
				// ignoreRules와 일치하면 결과에 추가하지 않음
				return nil
			}
		}

		// 해당 파일은 포함
		results = append(results, relPath)

		return nil
	}); err != nil {
		return nil, err
	}

	return results, nil
}

func (p parser) Parse(ctx context.Context) (ParsedMarkdownFiles, error) {
	filePaths, err := p.listMarkdownFilePaths(ctx)
	if err != nil {
		return nil, err
	}
	slog.DebugContext(ctx, "file path pattern match finished", "count", len(filePaths))

	var markdownFiles ParsedMarkdownFiles

	for _, filePath := range filePaths {
		var file []byte

		targetLanguages := make(config.LanguageCodes, 0, len(p.cfg.TargetLanguages))
		for _, targetLanguage := range p.cfg.TargetLanguages {
			var (
				originDir string
				fileName  string
			)
			originDir = filepath.Dir(filePath)
			fileName, err = FileNameWithoutExtension(filePath)
			if err != nil {
				return nil, err
			}

			targetFilePath := TargetFilePath(p.cfg.ContentDir, p.cfg.TargetPathRule, originDir, targetLanguage.String(), fileName)

			// 번역이 되지 않은 경우 번역 대상에 추가
			if slices.Contains(p.cfg.TranslatedPaths, targetFilePath) {
				slog.DebugContext(ctx, "skip already translated file",
					"path", targetFilePath,
					"language", targetLanguage.Name(),
				)
				continue
			}

			targetLanguages = append(targetLanguages, targetLanguage)
		}

		// 모든 언어가 이미 번역된 경우 해당 파일은 스킵
		if len(targetLanguages) == 0 {
			slog.DebugContext(ctx, "skip already translated file", "path", filePath)
			continue
		}

		file, err = os.ReadFile(path.Join(p.cfg.ContentDir, filePath))
		if err != nil {
			return nil, err
		}

		markdownFiles = append(markdownFiles, ParsedMarkdownFile{
			Path:            filePath,
			Markdown:        Markdown(file),
			TargetLanguages: targetLanguages,
		})
		slog.DebugContext(ctx, "markdown file parsed", "path", filePath)
	}

	return markdownFiles, nil
}
