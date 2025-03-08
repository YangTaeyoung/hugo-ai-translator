---
translated: true
---
[![Go Test Action](https://github.com/YangTaeyoung/hugo-ai-translator/actions/workflows/test-ci.yaml/badge.svg)](https://github.com/YangTaeyoung/hugo-ai-translator/actions/workflows/test-ci.yaml)  
![GitHub Release](https://img.shields.io/github/v/release/YangTaeyoung/hugo-ai-translator)  
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)  

# Hello Hugo AI Translator! 👋

This document is an AI translator that translates content stored in Hugo blogs.

It uses a model from [OpenAI](https://openai.com), and the translated results are saved according to the naming rules specified by the user.

# Translation

This translation has been done through `hugo-ai-translator`.

- [Korean](/README.md)  
- [English](/README.en.md)  
- [日本語](/README.ja.md)  
- [中文](/README.cn.md)  
- [Español](/README.es.md)  
- [Français](/README.fr.md)  
- [Deutsch](/README.de.md)  

# Supported Languages

| Language | Code |
|----------|------|
| Korean   | `ko` |
| English  | `en` |
| 日本語    | `ja` |
| 中文     | `cn` |
| Español  | `es` |
| Français | `fr` |
| Deutsch  | `de` |

# Installation

You can install it with a simple command.

```shell
go install github.com/YangTaeyoung/hugo-ai-translator@v1.1.0
```

# Configure

You can configure to run hugo-ai-translator with the following command.

```shell
hugo-ai-translator configure
```

For more detailed information about configuration, please refer to the [configuration](docs/configure.md) document.

# Usage

## Simple Translation

You can translate all markdown files in the current folder with simple rules.

### Quick Start

```shell
cd path/to/markdown-directory

hugo-ai-translator simple --source-language en \
  --target-language ko \
  --target-languages all \
  --model gpt-4o-mini \
  --api-key {open ai api key}
```  

## Rule Based Translation

You can translate by applying specific rules.

[Configure](docs/configure.md) should be done first, and the translation will proceed according to the set rules.

### Quick Start

If configured, translation can be done without using other options as follows:

```shell
hugo-ai-translator
```