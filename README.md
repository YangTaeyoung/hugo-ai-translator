![Go Version](https://img.shields.io/badge/Go-1.24-%23007d9c)
[![GoDoc](https://godoc.org/github.com/YangTaeyoung/hugo-ai-translator?status.svg)](https://pkg.go.dev/github.com/YangTaeyoung/hugo-ai-translator)
[![Go Test Action](https://github.com/YangTaeyoung/hugo-ai-translator/actions/workflows/test-ci.yaml/badge.svg)](https://github.com/YangTaeyoung/hugo-ai-translator/actions/workflows/test-ci.yaml)
![GitHub Release](https://img.shields.io/github/v/release/YangTaeyoung/hugo-ai-translator)
![GitHub License](https://img.shields.io/github/license/YangTaeyoung/hugo-ai-translator)
[![Go report](https://goreportcard.com/badge/github.com/YangTaeyoung/hugo-ai-translator)](https://goreportcard.com/report/github.com/YangTaeyoung/hugo-ai-translator)

# Hello Hugo AI Translator! 👋

해당 문서는 Hugo 블로그에 저장된 컨텐트를 번역하는 AI 번역기입니다.

[OpenAI](https://openai.com)의 모델을 사용하며, 번역된 결과는 사용자가 지정한 네이밍 룰에 따라 저장됩니다.

# Translation

해당 번역본은 `hugo-ai-translator`를 통해 번역되었습니다.

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

간단한 명령어로 설치할 수 있습니다.

```shell
go install github.com/YangTaeyoung/hugo-ai-translator@v1.1.1
```

# Configure

다음 커맨드를 통해 hugo-ai-translator를 구동하기 위한 설정을 할 수 있습니다.

```shell
hugo-ai-translator configure
```

설정에 대해 보다 자세한 내용은 [설정](docs/configure.md) 문서를 참고해주세요.

# Usage

## Simple Translation

현재 폴더 내에 있는 모든 마크다운을 단순한 룰로 번역할 수 있습니다.

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

특정한 룰을 적용하여 번역할 수 있습니다.

[Configure](docs/configure.md)가 선행되어야 하며, 설정된 룰에 따라 번역이 진행됩니다.

### Quick Start

설정이 있으면 아래와 같이 다른 옵션을 이용하지 않고도 번역이 가능합니다.

```shell
hugo-ai-translator
```
