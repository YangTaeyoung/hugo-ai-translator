package file

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// fixMismatchedQuotes는 front matter의 각 key: value 줄에서
// 값이 따옴표(" 또는 ')로 시작했으나 동일한 따옴표로 종료되지 않는 경우
// 내부 문자열을 추출해 올바른 double-quoted 문자열로 재생성합니다.
func fixMismatchedQuotes(frontMatter string) string {
	lines := strings.Split(frontMatter, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		// 주석이나 ':'가 없는 줄은 건너뜁니다.
		if strings.HasPrefix(trimmed, "#") || !strings.Contains(line, ":") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := parts[0]
		valuePart := strings.TrimSpace(parts[1])
		if len(valuePart) == 0 {
			continue
		}
		// 값이 " 또는 '로 시작하면
		if valuePart[0] == '"' || valuePart[0] == '\'' {
			quoteChar := valuePart[0]
			// 만약 끝나는 문자가 동일한 따옴표가 아니라면 수정 시도
			if valuePart[len(valuePart)-1] != quoteChar {
				// 만약 끝 문자가 다른 종류의 따옴표라면 제거
				if valuePart[len(valuePart)-1] == '"' || valuePart[len(valuePart)-1] == '\'' {
					valuePart = valuePart[:len(valuePart)-1]
				}
				inner := valuePart[1:]
				// strconv.Quote는 내부의 필요한 문자들을 escape 처리한 double-quoted 문자열을 생성합니다.
				newValue := strconv.Quote(inner)
				lines[i] = fmt.Sprintf("%s: %s", key, newValue)
			}
		}
	}
	return strings.Join(lines, "\n")
}

// convertSingleToDoubleQuoted는 front matter 내 단일 인용부호(')로 묶인 스칼라 값을
// double-quoted 스타일로 변환합니다.
// 단일 인용부호 내에서 두 번 연속된 인용부호(”)는 단일 인용부호 리터럴임을 감안합니다.
func convertSingleToDoubleQuoted(frontMatter string) string {
	// 정규표현식: key: 'value'
	// (:\s*) 캡처: 콜론과 공백
	// ((?:[^']|(?:''))*) 캡처: 내부 내용 (''는 리터럴 ')
	re := regexp.MustCompile(`(:\s*)'((?:[^']|(?:''))*)'`)
	return re.ReplaceAllStringFunc(frontMatter, func(match string) string {
		groups := re.FindStringSubmatch(match)
		if len(groups) < 3 {
			return match
		}
		prefix := groups[1]
		content := groups[2]
		// 단일 인용부호에서는 두 번 연속된 인용부호('')가 리터럴 '임을 의미하므로 복원합니다.
		unescapedContent := strings.ReplaceAll(content, "''", "'")
		quoted := strconv.Quote(unescapedContent)
		return prefix + quoted
	})
}

// toYamlNode은 임의의 value를 YAML 노드로 변환합니다.
func toYamlNode(value interface{}) (*yaml.Node, error) {
	if str, ok := value.(string); ok {
		// 내부에 특별한 문자(쌍따옴표, 따옴표, 탭, 콜론 등)가 포함되면 double-quoted 스타일을 강제합니다.
		if strings.ContainsAny(str, "\"'\t:") {
			return &yaml.Node{
				Kind:  yaml.ScalarNode,
				Tag:   "!!str",
				Value: str,
				Style: yaml.DoubleQuotedStyle,
			}, nil
		}
	}
	data, err := yaml.Marshal(value)
	if err != nil {
		return nil, errors.Wrapf(err, "YAML marshalling failed. value: %+v", value)
	}
	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, errors.Wrapf(err, "YAML unmarshalling failed. data: %s", data)
	}
	if len(doc.Content) > 0 {
		return doc.Content[0], nil
	}
	return nil, errors.New("Empty YAML node")
}

// updateFrontmatterPreserveOrder는 기존 front matter를 YAML 노드로 파싱하여
// updates에 담긴 key-value 쌍을 순서를 유지하면서 업데이트(또는 추가)합니다.
func updateFrontmatterPreserveOrder(frontMatter string, updates map[string]interface{}) (string, error) {
	// 먼저 quoting 문제를 해결
	fixed := fixMismatchedQuotes(frontMatter)
	// 단일 인용부호로 묶인 값들을 double-quoted 스타일로 변환
	fixed = convertSingleToDoubleQuoted(fixed)

	var node yaml.Node
	if err := yaml.Unmarshal([]byte(fixed), &node); err != nil {
		return "", errors.Wrapf(err, "YAML unmarshalling failed. frontMatter: %+v", fixed)
	}
	if len(node.Content) == 0 {
		return "", errors.New("Invalid YAML frontmatter")
	}
	mapping := node.Content[0]
	if mapping.Kind != yaml.MappingNode {
		return "", errors.Errorf("Frontmatter must be a mapping node but got %+v", mapping.Value)
	}

	updated := make(map[string]bool)
	// 기존 순서를 유지하며 업데이트 진행
	for i := 0; i < len(mapping.Content); i += 2 {
		keyNode := mapping.Content[i]
		if newVal, ok := updates[keyNode.Value]; ok {
			newValueNode, err := toYamlNode(newVal)
			if err != nil {
				return "", err
			}
			mapping.Content[i+1] = newValueNode
			updated[keyNode.Value] = true
		}
	}
	// 존재하지 않는 key는 맨 뒤에 추가
	for k, v := range updates {
		if !updated[k] {
			keyNode := &yaml.Node{
				Kind:  yaml.ScalarNode,
				Value: k,
				Tag:   "!!str",
			}
			valueNode, err := toYamlNode(v)
			if err != nil {
				return "", err
			}
			mapping.Content = append(mapping.Content, keyNode, valueNode)
		}
	}
	out, err := yaml.Marshal(&node)
	if err != nil {
		return "", errors.Wrap(err, "YAML marshalling 실패")
	}
	return string(out), nil
}

// WriteMarkdownWithFrontmatter는 파일 내용을 업데이트합니다.
// keyValues는 key와 value 쌍으로 전달되며, key는 string, value는 any 타입이어야 합니다.
func WriteMarkdownWithFrontmatter(path string, file []byte, perm os.FileMode, keyValues ...interface{}) error {
	if len(keyValues)%2 != 0 {
		return errors.New("keyValues must be provided key, value pairs")
	}
	updates := make(map[string]interface{})
	for i := 0; i < len(keyValues); i += 2 {
		key, ok := keyValues[i].(string)
		if !ok {
			return errors.New("key in keyValues must be string")
		}
		updates[key] = keyValues[i+1]
	}

	contentStr := string(file)
	var newContent string

	if strings.HasPrefix(contentStr, "---") {
		// 기존 front matter가 있는 경우
		parts := strings.SplitN(contentStr, "---", 3)
		if len(parts) < 3 {
			return errors.Errorf("invalid front matter format. contentStr: %s", contentStr)
		}
		// parts[1]에는 YAML front matter, parts[2]에는 본문 내용이 있음
		newYaml, err := updateFrontmatterPreserveOrder(parts[1], updates)
		if err != nil {
			return err
		}
		newContent = fmt.Sprintf("---\n%s---%s", newYaml, parts[2])
	} else {
		// front matter가 없는 경우, 새로 생성
		newYaml, err := yaml.Marshal(updates)
		if err != nil {
			return errors.Wrap(err, "YAML marshalling 실패")
		}
		newContent = fmt.Sprintf("---\n%s---\n%s", newYaml, contentStr)
	}

	if err := os.WriteFile(path, []byte(newContent), perm); err != nil {
		return errors.Wrap(err, "파일 쓰기 실패")
	}

	return nil
}
