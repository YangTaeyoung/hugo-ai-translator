---
translated: true
---
![Go Version](https://img.shields.io/badge/Go-1.24-%23007d9c)
[![GoDoc](https://godoc.org/github.com/YangTaeyoung/hugo-ai-translator?status.svg)](https://pkg.go.dev/github.com/YangTaeyoung/hugo-ai-translator)
[![Go Test Action](https://github.com/YangTaeyoung/hugo-ai-translator/actions/workflows/test-ci.yaml/badge.svg)](https://github.com/YangTaeyoung/hugo-ai-translator/actions/workflows/test-ci.yaml)
![GitHub Release](https://img.shields.io/github/v/release/YangTaeyoung/hugo-ai-translator)
![GitHub License](https://img.shields.io/github/license/YangTaeyoung/hugo-ai-translator)
[![Go report](https://goreportcard.com/badge/github.com/YangTaeyoung/hugo-ai-translator)](https://goreportcard.com/report/github.com/YangTaeyoung/hugo-ai-translator)

# Hello Hugo AI Translator! 👋

This document is an AI translator that translates content stored in Hugo blogs.

It uses a model from [OpenAI](https://openai.com), and the translated results are saved according to the naming rules specified by the user.

# Translation

This translation has been processed through `hugo-ai-translator`.

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

You can install it with a simple command.

```shell
go install github.com/YangTaeyoung/hugo-ai-translator@v1.1.1
```

# Configure

You can set up configurations to run hugo-ai-translator using the following command.

```shell
hugo-ai-translator configure
```

For more details on the configuration, please refer to the [Configuration](docs/configure.md) document.

# Usage

## Simple Translation

You can perform a simple translation of all markdown files in the current folder using basic rules.

### Quick Start

```shell
cd path/to/markdown-directory

hugo-ai-translator simple --source-language en \
  --target-language ko \
  --target-languages all \
  --model gpt-4o-mini \
  --api-key {open ai api key}
``` 

## Rule Base Translation

You can apply specific rules for translation.

[Configure](docs/configure.md) needs to be done first, and the translation will proceed according to the established rules.

### Quick Start

With configurations in place, you can translate without using other options as follows:

```shell
hugo-ai-translator
``` 
