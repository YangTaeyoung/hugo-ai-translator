package file

import (
	"bytes"
	"context"
	"io/fs"
	"log/slog"
	"os"
	"path"
	"path/filepath"
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
	TargetPathRule  string
}

type TranslateFrontMatter struct {
	Translated bool `yaml:"translated"`
}

type TranslateFrontMatters []TranslateFrontMatter

type Parser interface {
	Parse(ctx context.Context) (ParsedMarkdownFiles, error)
	Simple(ctx context.Context) (ParsedMarkdownFiles, error)
}

type parser struct {
	cfg ParserConfig
}

func NewParser(cfg ParserConfig) Parser {
	return &parser{
		cfg: cfg,
	}
}

func (p parser) Simple(_ context.Context) (ParsedMarkdownFiles, error) {
	// ContentDir에 있는 모든 .md파일을 읽어서 반환
	var markdownFiles ParsedMarkdownFiles

	paths, err := os.ReadDir(p.cfg.ContentDir)
	if err != nil {
		return nil, err
	}

	for _, path := range paths {
		var file []byte

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

		markdownFiles = append(markdownFiles, ParsedMarkdownFile{
			Path:            path.Name(),
			Markdown:        Markdown(file),
			TargetLanguages: p.cfg.TargetLanguages,
		})
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

func (p parser) Parse(ctx context.Context) (ParsedMarkdownFiles, error) {
	var (
		originMarkdownFiles ParsedMarkdownFiles
		markdownFiles       ParsedMarkdownFiles
		translatedMap       = make(map[string]bool)
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
		var file []byte

		file, err = os.ReadFile(path.Join(p.cfg.ContentDir, filePath))
		if err != nil {
			return nil, err
		}

		var frontMatter TranslateFrontMatter

		if _, err = frontmatter.Parse(bytes.NewReader(file), &frontMatter); err != nil {
			return nil, err
		}

		if frontMatter.Translated {
			translatedMap[filePath] = true
			slog.DebugContext(ctx, "skip already translated file", "path", filePath)
			continue
		}

		originMarkdownFiles = append(originMarkdownFiles, ParsedMarkdownFile{
			Path:     filePath,
			Markdown: Markdown(file),
		})
	}

	// Origin Markdown 파일을 기반으로 번역할 언어를 설정. 이전 스텝에서 이미 번역된 파일의 경우 해당 파일의 번역 언어를 설정하지 않음
	for _, originMarkdownFile := range originMarkdownFiles {
		for _, lang := range p.cfg.TargetLanguages {
			var fileName string

			fileName, err = FileNameWithoutExtension(originMarkdownFile.Path)
			if err != nil {
				return nil, err
			}

			targetFilePath := TargetFilePath(p.cfg.TargetPathRule, filepath.Dir(originMarkdownFile.Path), lang.String(), fileName)
			if strings.HasPrefix(targetFilePath, "./") {
				targetFilePath = targetFilePath[2:]
			}
			slog.Debug("output path for translated markdown", "path", targetFilePath)

			if _, ok := translatedMap[targetFilePath]; ok {
				slog.DebugContext(ctx, "skip already translated language", "path", targetFilePath, "language", lang)
				continue
			}

			originMarkdownFile.TargetLanguages = append(originMarkdownFile.TargetLanguages, lang)
		}

		if len(originMarkdownFile.TargetLanguages) == 0 {
			slog.DebugContext(ctx, "skip origin markdown file because no target languages", "path", originMarkdownFile.Path)
			continue
		}

		markdownFiles = append(markdownFiles, originMarkdownFile)
	}

	return markdownFiles, nil
}
