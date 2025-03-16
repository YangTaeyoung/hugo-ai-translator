---
translated: true
---
![Go Version](https://img.shields.io/badge/Go-1.24-%23007d9c)  
[![GoDoc](https://godoc.org/github.com/YangTaeyoung/hugo-ai-translator?status.svg)](https://pkg.go.dev/github.com/YangTaeyoung/hugo-ai-translator)  
[![Go Test Action](https://github.com/YangTaeyoung/hugo-ai-translator/actions/workflows/test-ci.yaml/badge.svg)](https://github.com/YangTaeyoung/hugo-ai-translator/actions/workflows/test-ci.yaml)  
![GitHub Release](https://img.shields.io/github/v/release/YangTaeyoung/hugo-ai-translator)  
![GitHub License](https://img.shields.io/github/license/YangTaeyoung/hugo-ai-translator)  
[![Go report](https://goreportcard.com/badge/github.com/YangTaeyoung/hugo-ai-translator)](https://goreportcard.com/report/github.com/YangTaeyoung/hugo-ai-translator)  

# こんにちはHugo AI翻訳者！👋

この文書は、Hugoブログに保存されたコンテンツを翻訳するAI翻訳器です。

[OpenAI](https://openai.com)のモデルを使用しており、翻訳された結果はユーザーが指定した名前のルールに従って保存されます。

# 翻訳

この翻訳版は`hugo-ai-translator`を通じて翻訳されました。

- [한국어](/README.md)  
- [English](/README.en.md)  
- [日本語](/README.ja.md)  
- [中文](/README.cn.md)  
- [Español](/README.es.md)  
- [Français](/README.fr.md)  
- [Deutsch](/README.de.md)  


# 対応言語

| 言語    | コード |  
|---------|-------|  
| 한국어      | `ko`  |  
| English  | `en`  |  
| 日本語      | `ja`  |  
| 中文       | `cn`  |  
| Español  | `es`  |  
| Français | `fr`  |  
| Deutsch  | `de`  |

# インストール

簡単なコマンドでインストールできます。

```shell
go install github.com/YangTaeyoung/hugo-ai-translator@v1.1.1
```

# 設定

次のコマンドを使用してhugo-ai-translatorを実行するための設定を行うことができます。

```shell
hugo-ai-translator configure
```

設定に関する詳細は[設定](docs/configure.md)文書を参照してください。

# 使用法

## 簡単な翻訳

現在のフォルダー内にあるすべてのMarkdownを簡単なルールで翻訳できます。

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

[Configure](docs/configure.md)が事前に行われる必要があり、設定されたルールに従って翻訳が進行します。

### クイックスタート

設定がある場合、以下のように他のオプションを使用せずに翻訳が可能です。

```shell
hugo-ai-translator
```