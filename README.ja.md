---
translated: true
---
[![Go Test Action](https://github.com/YangTaeyoung/hugo-ai-translator/actions/workflows/test-ci.yaml/badge.svg)](https://github.com/YangTaeyoung/hugo-ai-translator/actions/workflows/test-ci.yaml)
![GitHub Release](https://img.shields.io/github/v/release/YangTaeyoung/hugo-ai-translator)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

# こんにちはHugo AI翻訳者! 👋

この文書はHugoブログに保存されたコンテンツを翻訳するAI翻訳ツールです。

[OpenAI](https://openai.com)のモデルを使用しており、翻訳結果はユーザーが指定した命名ルールに従って保存されます。

# 翻訳

この翻訳文は`hugo-ai-translator`を通じて翻訳されました。

- [한국어](/README.md)
- [English](/README.en.md)
- [日本語](/README.ja.md)
- [中文](/README.cn.md)
- [Español](/README.es.md)
- [Français](/README.fr.md)
- [Deutsch](/README.de.md)


# 対応言語

| 言語     | コード |
|----------|------|
| 한국어     | `ko` |
| English  | `en` |
| 日本語     | `ja` |
| 中文      | `cn` |
| Español  | `es` |
| Français | `fr` |
| Deutsch  | `de` |

# インストール

簡単なコマンドでインストールできます。

```shell
go install github.com/YangTaeyoung/hugo-ai-translator@v1.1.0
```

# 設定

次のコマンドでhugo-ai-translatorを起動するための設定ができます。

```shell
hugo-ai-translator configure
```

設定についての詳細は[設定](docs/configure.md)の文書を参照してください。

# 使用法

## 単純な翻訳

現在のフォルダ内のすべてのマークダウンを単純なルールで翻訳できます。

### クイックスタート

```shell
cd path/to/markdown-directory

hugo-ai-translator simple --source-language en \
  --target-language ko \
  --target-languages all \
  --model gpt-4o-mini \
  --api-key {open ai api key}
``` 

## ルールベース翻訳

特定のルールを適用して翻訳できます。

[Configure](docs/configure.md)が先行される必要があり、設定されたルールに従って翻訳が行われます。

### クイックスタート

設定があれば、以下のように他のオプションを使用せずに翻訳できます。

```shell
hugo-ai-translator
```