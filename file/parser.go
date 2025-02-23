package file

import (
	"context"
	"io/fs"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

type Markdown string

func (m Markdown) String() string {
	return string(m)
}

type ParsedMarkdownFile struct {
	Path     string
	Markdown Markdown
}

type ParsedMarkdownFiles []ParsedMarkdownFile

type ParserConfig struct {
	ContentDir  string
	IgnoreRules []string
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
		// contentDir 기준의 상대 경로 계산
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
	filePaths, err := p.listMarkdownFilePaths()
	if err != nil {
		return nil, err
	}
	slog.DebugContext(ctx, "file path pattern match finished", "count", len(filePaths))

	var markdownFiles ParsedMarkdownFiles

	for _, filePath := range filePaths {
		var file []byte

		file, err = os.ReadFile(path.Join(p.cfg.ContentDir, filePath))
		if err != nil {
			return nil, err
		}

		markdownFiles = append(markdownFiles, ParsedMarkdownFile{
			Path:     filePath,
			Markdown: Markdown(file),
		})
		slog.DebugContext(ctx, "markdown file parsed", "path", filePath)
	}

	return markdownFiles, nil
}
