---
translated: true
---
# 配置
本节将介绍 Hugo AI Translator 的配置文件。

## 位置
Hugo AI Translator 的配置文件位于 `~/.hugo-ai-translator/config.yaml`。

```shell
cat ~/.hugo-ai-translator/config.yaml
```

## Schema 
```yaml
openai:
    model: gpt-4o-mini
    api_key: {your-openai-api-key}
translator:
    content_dir: ~/dev/personal/YangTaeyoung.github.io/content
    source:
        source_language: ko
        ignore_rules: []
    target:
        target_languages:
            - en
            - ja
            - fr
            - de
        target_path_rule: '{origin}/{fileName}.{language}.md'
```

## `openai`
- `model`: 指定 OpenAI API 中使用的模型，关于模型的信息请参考 [Open AI Models](https://platform.openai.com/docs/models)  
- `api_key`: 指定用于使用 OpenAI API 的 API 密钥

## `translator`
- `content_dir`: 选择 Hugo 的内容目录，作为翻译结果保存至根路径的角色。
- `source`
    - `source_language`: 指定要翻译的 markdown 的原始语言。
    - `ignore_rules`: 指定不进行翻译的 markdown 文件，可以使用 *, ** 等通配符。  
      - 例如 `ignore_rules: ["*.en.md", "*.ko.md"]`， `ignore_rules: ["some/path/**"]`
- `target`
  - `target_languages`: 指定要翻译的语言，可以指定多种语言。支持的语言请参考 [Supported Languages](../README.md#supported-languages)。  
  - 例如 `target_languages: ["en", "ja", "fr", "de"]`

### `translator.target_path_rule`
指定翻译结果保存的路径，可以使用 `{origin}`、`{fileName}` 和 `{language}` 的保留字。
- `{origin}`: 指原始文件的目录路径，从 `translator.content_dir` 开始。对于 `~/dev/personal/YangTaeyoung.github.io/content/some/index.md`， `~/dev/personal/YangTaeyoung.github.io/content` 是 `content_dir`， `some` 是 `origin` 。
- `{fileName}`: 指文件名（不含扩展名）。
- `{language}`: 指待翻译语言的代码。

#### 示例
设置值如下所示。
```yaml
target:
    content_dir: /content
    target_path_rule: '{origin}/{fileName}.{language}.md'
    target_languages:
        - en
        - ja
```

翻译前文件架构如下所示。
```
/content
    /some
        index.md
    index.md
```

翻译后文件架构如下所示。
```
/content
    /some
        hamburger.md
        hamburger.en.md <-- 翻译的文件
        hamburger.ja.md <-- 翻译的文件
    index.md
    index.en.md <-- 翻译的文件
    index.ja.md <-- 翻译的文件
```