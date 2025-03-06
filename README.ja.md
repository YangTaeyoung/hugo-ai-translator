---
translated: true
---
# Hello Hugo AI Translator! 👋

この文書はHugoブログに保存されたコンテンツを翻訳するAI翻訳機です。

[OpenAI](https://openai.com)のモデルを使用しており、翻訳結果はユーザーが指定した命名ルールに従って保存されます。

# Translation

この翻訳は`hugo-ai-translator`を通じて行われました。

- [한국어](/README.md)
- [English](/README.en.md)
- [日本語](/README.ja.md)
- [中文](/README.cn.md)
- [Español](/README.es.md)
- [Français](/README.fr.md)
- [Deutsch](/README.de.md)


# Supported Languages

| Language | Code |
|----------|------|
| 한국어      | `ko` |
| English  | `en` |
| 日本語      | `ja` |
| 中文       | `cn` |
| Español  | `es` |
| Français | `fr` |
| Deutsch  | `de` |

# Installation

簡単なコマンドでインストールできます。

```shell
go install github.com/YangTaeyoung/hugo-ai-translator@v1.0.1
```

# Configure

次のコマンドを使用してhugo-ai-translatorを実行するために設定できます。

```shell
hugo-ai-translator configure
```

設定についての詳細は[設定](docs/configure.md)文書を参照してください。

# Usage

## Simple Translation

現在のフォルダー内にあるすべてのMarkdownを単純なルールで翻訳できます。

### Quick Start

```shell
cd path/to/markdown-directory

hugo-ai-translator simple --source-language en \
  --target-language ko \
  --target-language ja \
  --model gpt-4 \
  --api-key {open ai api key}
``` 

## Rull Base Translation

特定のルールを適用して翻訳できます。

[Configure](docs/configure.ja.md)が先行する必要があり、設定されたルールに従って翻訳が進められます。

### Quick Start

設定があれば、以下のように他のオプションを使わずに翻訳ができます。

```shell
hugo-ai-translator
```