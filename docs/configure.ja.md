---
translated: true
---
# 設定
このセクションでは、Hugo AI Translatorの設定ファイルについて説明します。

## 場所
Hugo AI Translatorの設定ファイルは`~/.hugo-ai-translator/config.yaml`にあります。

```shell
cat ~/.hugo-ai-translator/config.yaml
```

## スキーマ  
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
- `model`: OpenAI APIで使用するモデルを指定します。モデルの情報は[Open AI Models](https://platform.openai.com/docs/models)を参照してください。
- `api_key`: OpenAI APIを使用するためのAPIキーを指定します。

## `translator`
- `content_dir`: Hugoのコンテンツディレクトリを選択します。翻訳された結果が保存されるルートパスの役割を果たします。
- `source`
    - `source_language`: 翻訳するマークダウンの元の言語を指定します。
    - `ignore_rules`: 翻訳しないマークダウンファイルを指定します。*, ** などのワイルドカードを使用できます。
      - 例) `ignore_rules: ["*.en.md", "*.ko.md"]`, `ignore_rules: ["some/path/**"]`
- `target`
  - `target_languages`: 翻訳する言語を指定します。複数の言語を指定できます。対応言語は[Supported Languages](../README.md#supported-languages)を参照してください。
  - 例) `target_languages: ["en", "ja", "fr", "de"]`

### `translator.target_path_rule`
翻訳された結果が保存されるパスを指定します。`{origin}`、`{fileName}`、`{language}`の予約語を利用できます。
- `{origin}`: `translator.content_dir`からの元のファイルのディレクトリパスを意味します。`~/dev/personal/YangTaeyoung.github.io/content/some/index.md`の場合、`~/dev/personal/YangTaeyoung.github.io/content`が`content_dir`、`some`が`origin`になります。
- `{fileName}`: 拡張子を除いたファイル名を意味します。
- `{language}`: 翻訳される言語のコードを意味します。

#### 例
設定値は次のようになります。
```yaml
target:
    content_dir: /content
    target_path_rule: '{origin}/{fileName}.{language}.md'
    target_languages:
        - en
        - ja
```

翻訳前のファイルスキーマは次のように構成されています。

```
/content
    /some
        index.md
    index.md
```

翻訳後のファイルスキーマは次のように構成されます。
```
/content
    /some
        hamburger.md
        hamburger.en.md <-- 翻訳されたファイル
        hamburger.ja.md <-- 翻訳されたファイル
    index.md
    index.en.md <-- 翻訳されたファイル
    index.ja.md <-- 翻訳されたファイル
```
