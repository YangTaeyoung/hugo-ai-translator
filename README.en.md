---
translated: true
---
# Hello Hugo AI Translator! 👋

This document is an AI translator that translates content stored in Hugo blogs.

It uses the model from [OpenAI](https://openai.com), and the translated results are saved according to the naming rules specified by the user.

# Translation

This translation was performed via `hugo-ai-translator`.

- [한국어](/README.md)
- [English](/README.en.md)
- [日本語](/README.ja.md)
- [中文](/README.cn.md)
- [Español](/README.es.md)
- [Français](/README.fr.md)
- [Deutsch](/README.de.md)


# Supported Languages

| Language   | Code |
|------------|------|
| 한국어       | `ko` |
| English    | `en` |
| 日本語      | `ja` |
| 中文        | `cn` |
| Español    | `es` |
| Français   | `fr` |
| Deutsch    | `de` |

# Installation

It can be installed with a simple command.

```shell
go install github.com/YangTaeyoung/hugo-ai-translator@v1.0.0
```

# Configure

You can set up the hugo-ai-translator to run using the following command.

```shell
hugo-ai-translator configure
```

For more details on the configuration, please refer to the [configuration](docs/configure.en.md) document.

# Usage

## Simple Translation

You can translate all markdowns in the current folder using a simple rule.

### Quick Start

```shell
cd path/to/markdown-directory

hugo-ai-translator simple --source-language en \
  --target-language ko \
  --target-language ja \
  --model gpt-4 \
  --api-key {open ai api key}
``` 

## Rule-Based Translation

You can translate by applying specific rules.

[Configure](docs/configure.md) must precede this, and the translation will proceed according to the configured rules.

### Quick Start

If the configuration is set, translation can be done without using other options as follows.

```shell
hugo-ai-translator
```