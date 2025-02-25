package file

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// toYamlNode은 임의의 value를 YAML 노드로 변환합니다.
func toYamlNode(value interface{}) (*yaml.Node, error) {
	data, err := yaml.Marshal(value)
	if err != nil {
		return nil, errors.Wrapf(err, "YAML marshalling failed. value: %+v", value)
	}
	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, errors.Wrapf(err, "YAML unmarshalling failed. data: %s", data)
	}
	// doc은 DocumentNode이며, 실제 값은 Content의 첫번째 요소에 있음
	if len(doc.Content) > 0 {
		return doc.Content[0], nil
	}
	return nil, fmt.Errorf("빈 YAML 노드")
}

// updateFrontmatterPreserveOrder는 기존 frontmatter를 YAML 노드로 파싱하여
// updates에 담긴 key-value 쌍을 순서를 유지하면서 업데이트(또는 추가)합니다.
func updateFrontmatterPreserveOrder(frontMatter string, updates map[string]interface{}) (string, error) {
	var node yaml.Node
	if err := yaml.Unmarshal([]byte(frontMatter), &node); err != nil {
		return "", errors.Wrapf(err, "YAML unmarshalling failed. frontMatter: %+v", frontMatter)
	}
	if len(node.Content) == 0 {
		return "", errors.New("Invalid YAML frontmatter")
	}
	mapping := node.Content[0]
	if mapping.Kind != yaml.MappingNode {
		return "", errors.Errorf("frontmatter must be a mapping node but got %+v", mapping.Value)
	}

	updated := make(map[string]bool)
	// 기존 노드 순서 유지하며 업데이트 진행
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

// UpdateMarkdownFrontmatterPreserveOrder는 파일 내용을 업데이트합니다.
// 파라미터 keyValues는 key와 value 쌍으로 전달되며, key는 string, value는 any 타입이어야 합니다.
func WriteMarkdownWithFrontmatter(path string, file []byte, perm os.FileMode, keyValues ...interface{}) error {
	if len(keyValues)%2 != 0 {
		return fmt.Errorf("keyValues는 반드시 key, value 쌍으로 제공되어야 합니다")
	}
	updates := make(map[string]interface{})
	for i := 0; i < len(keyValues); i += 2 {
		key, ok := keyValues[i].(string)
		if !ok {
			return fmt.Errorf("keyValues의 key는 반드시 string이어야 합니다")
		}
		updates[key] = keyValues[i+1]
	}

	contentStr := string(file)
	var newContent string

	if strings.HasPrefix(contentStr, "---") {
		// 기존 frontmatter가 있는 경우
		parts := strings.SplitN(contentStr, "---", 3)
		if len(parts) < 3 {
			return fmt.Errorf("유효하지 않은 frontmatter 형식")
		}
		// parts[1]에는 YAML frontmatter, parts[2]에는 본문 내용이 있음
		newYaml, err := updateFrontmatterPreserveOrder(parts[1], updates)
		if err != nil {
			return err
		}
		newContent = fmt.Sprintf("---\n%s---%s", newYaml, parts[2])
	} else {
		// frontmatter가 없는 경우, 새로 생성
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
