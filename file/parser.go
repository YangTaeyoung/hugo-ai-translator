package file

import (
	"bytes"
	"context"
	"io/fs"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/YangTaeyoung/hugo-ai-translator/config"
	"github.com/adrg/frontmatter"
	"github.com/bmatcuk/doublestar/v4"
)

type Markdown string

type Depth int

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
	IgnoreRules     []string
	TargetLanguages config.LanguageCodes
	SourceLanguage  config.LanguageCode
	TargetPathRule  string
}

type Parser interface {
	Parse(ctx context.Context) (MarkdownFiles, error)
	Simple(ctx context.Context) (MarkdownFiles, error)
}

type parser struct {
	cfg ParserConfig
}

func NewParser(cfg ParserConfig) Parser {
	return &parser{
		cfg: cfg,
	}
}

func (p parser) Simple(_ context.Context) (MarkdownFiles, error) {
	// ContentDir에 있는 모든 .md파일을 읽어서 반환
	var markdownFiles MarkdownFiles

	paths, err := os.ReadDir(p.cfg.ContentDir)
	if err != nil {
		return nil, err
	}

	for _, path := range paths {
		var (
			file     []byte
			fileName string
		)

		if path.IsDir() {
			continue
		}

		if !strings.HasSuffix(path.Name(), ".md") {
			continue
		}

		file, err = os.ReadFile(filepath.Join(p.cfg.ContentDir, path.Name()))
		if err != nil {
			return nil, err
		}

		fileName, err = FileNameWithoutExtension(path.Name())
		if err != nil {
			return nil, err
		}

		for _, language := range p.cfg.TargetLanguages {
			markdownFiles = append(markdownFiles, MarkdownFile{
				OriginDir: filepath.Dir(path.Name()),
				FileName:  fileName,
				Content:   Markdown(file),
				Language:  language,
			})
		}

	}

	return markdownFiles, nil
}

func (p parser) listMarkdownFilePaths() ([]string, error) {
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

func (p parser) Parse(ctx context.Context) (MarkdownFiles, error) {
	var (
		markdownFiles MarkdownFiles
		translatedMap = make(map[string]bool)
	)

	filePaths, err := p.listMarkdownFilePaths()
	if err != nil {
		return nil, err
	}
	slog.DebugContext(ctx, "file path pattern match finished", "count", len(filePaths))

	// 파일 경로를 순회하면서 front matter를 파싱
	// 파싱한 front matter에 translated가 true로 설정되어 있으면 이미 번역된 파일로 간주,
	// translatedMap에 파일 경로를 키로 추가하고 이미 번역된 파일을 다시 번역하지 않기 위해 skip
	for _, filePath := range filePaths {
		var (
			file     []byte
			fileName string
		)

		file, err = os.ReadFile(path.Join(p.cfg.ContentDir, filePath))
		if err != nil {
			return nil, err
		}

		var frontMatter struct {
			Translated bool `yaml:"translated"`
		}

		if _, err = frontmatter.Parse(bytes.NewReader(file), &frontMatter); err != nil {
			return nil, err
		}

		if frontMatter.Translated {
			translatedMap[filePath] = true
			slog.DebugContext(ctx, "skip already translated file", "path", filePath)
			continue
		}

		fileName, err = FileNameWithoutExtension(filePath)
		if err != nil {
			return nil, err
		}

		for _, lang := range p.cfg.TargetLanguages {
			targetFilePath := strings.TrimPrefix(
				TargetFilePath(p.cfg.TargetPathRule, filepath.Dir(filePath), lang.String(), fileName),
				"./",
			)

			slog.Debug("output path for translated markdown", "path", targetFilePath)

			if _, ok := translatedMap[targetFilePath]; ok {
				slog.DebugContext(ctx, "skip already translated language", "path", targetFilePath, "language", lang)
				continue
			}

			originDir := filepath.Dir(filePath)

			// OriginDir이 SourceLanguage를 포함하고 있는 경우 제거
			// ex) /en/docs -> /docs
			fragments := strings.Split(originDir, "/")
			if len(fragments) > 0 {
				if i := slices.Index(fragments, p.cfg.SourceLanguage.String()); i >= 0 {
					fragments = append(fragments[:i], fragments[i+1:]...)
				}

				originDir = filepath.Join(fragments...)
			}

			markdownFiles = append(markdownFiles, MarkdownFile{
				OriginDir: originDir,
				Content:   Markdown(file),
				FileName:  fileName,
				Language:  lang,
			})
		}
	}

	return markdownFiles, nil
}
