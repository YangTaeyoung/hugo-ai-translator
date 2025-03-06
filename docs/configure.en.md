---
translated: true
---
# Configuration
This section provides an explanation of the configuration file for the Hugo AI Translator.

## Location
The configuration file for Hugo AI Translator can be found at `~/.hugo-ai-translator/config.yaml`.

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
- `model`: Specifies the model to be used with the OpenAI API. For more information about the model, please refer to [Open AI Models](https://platform.openai.com/docs/models).
- `api_key`: Specifies the API key to use with the OpenAI API.

## `translator`
- `content_dir`: Selects the content directory of Hugo. It serves as the root path where the translated results will be stored.
- `source`
    - `source_language`: Specifies the original language of the markdown to be translated.
    - `ignore_rules`: Specifies which markdown files should not be translated. Wildcards such as *, ** can be used.
      - ex) `ignore_rules: ["*.en.md", "*.ko.md"]`, `ignore_rules: ["some/path/**"]`
- `target`
  - `target_languages`: Specifies the languages to translate into. Multiple languages can be specified. For supported languages, please refer to [Supported Languages](../README.md#supported-languages).
  - ex) `target_languages: ["en", "ja", "fr", "de"]`

### `translator.target_path_rule`
Specifies the path where the translated results will be saved. You can use the reserved words `{origin}`, `{fileName}`, and `{language}`.
- `{origin}`: Refers to the directory path of the original file starting from `translator.content_dir`. For example, in `~/dev/personal/YangTaeyoung.github.io/content/some/index.md`, `~/dev/personal/YangTaeyoung.github.io/content` is `content_dir`, and `some` is `origin`.
- `{fileName}`: Refers to the file name without the extension.
- `{language}`: Refers to the code of the language into which it will be translated.

#### Example
The configuration values are as follows:
```yaml
target:
    content_dir: /content
    target_path_rule: '{origin}/{fileName}.{language}.md'
    target_languages:
        - en
        - ja
```

The file schema before translation is structured as follows:
```
/content
    /some
        index.md
    index.md
```

The file schema after translation is structured as follows:
```
/content
    /some
        hamburger.md
        hamburger.en.md <-- Translated file
        hamburger.ja.md <-- Translated file
    index.md
    index.en.md <-- Translated file
    index.ja.md <-- Translated file
```