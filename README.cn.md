---
translated: true
---
[![Go Test Action](https://github.com/YangTaeyoung/hugo-ai-translator/actions/workflows/test-ci.yaml/badge.svg)](https://github.com/YangTaeyoung/hugo-ai-translator/actions/workflows/test-ci.yaml)
![GitHub Release](https://img.shields.io/github/v/release/YangTaeyoung/hugo-ai-translator)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

# 你好 Hugo AI 翻译器! 👋

该文档是用来翻译存储在 Hugo 博客中的内容的 AI 翻译器。

使用 [OpenAI](https://openai.com) 的模型，翻译结果将根据用户指定的命名规则保存。

# 翻译

该翻译版本是通过 `hugo-ai-translator` 完成的。

- [한국어](/README.md)
- [English](/README.en.md)
- [日本語](/README.ja.md)
- [中文](/README.cn.md)
- [Español](/README.es.md)
- [Français](/README.fr.md)
- [Deutsch](/README.de.md)


# 支持的语言

| 语言   | 代码  |
|--------|-------|
| 한국어   | `ko`  |
| English | `en`  |
| 日本語   | `ja`  |
| 中文    | `cn`  |
| Español | `es`  |
| Français| `fr`  |
| Deutsch | `de`  |

# 安装

可以通过简单的命令进行安装。

```shell
go install github.com/YangTaeyoung/hugo-ai-translator@v1.1.0
```

# 配置

可以通过以下命令为 hugo-ai-translator 进行设置。

```shell
hugo-ai-translator configure
```

有关设置的更多详细信息，请参阅 [设置](docs/configure.md) 文档。

# 使用

## 简单翻译

可以按照简单的规则翻译当前文件夹内的所有 Markdown。

### 快速开始

```shell
cd path/to/markdown-directory

hugo-ai-translator simple --source-language en \
  --target-language ko \
  --target-languages all \
  --model gpt-4o-mini \
  --api-key {open ai api key}
``` 

## 规则基础翻译

可以应用特定的规则进行翻译。

需要先进行 [Configure](docs/configure.md)，翻译将根据设置的规则进行。

### 快速开始

如果有设置，则可以直接进行翻译，而无需使用其他选项。

```shell
hugo-ai-translator
``` 
