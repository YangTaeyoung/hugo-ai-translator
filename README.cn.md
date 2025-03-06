---
translated: true
---
# Hello Hugo AI Translator! 👋

该文档是一个AI翻译器，用于翻译存储在Hugo博客中的内容。

使用[OpenAI](https://openai.com)的模型，翻译结果会根据用户指定的命名规则保存。

# Translation

该翻译本是通过 `hugo-ai-translator` 翻译的。

- [韩语](/README.md)
- [English](/README.en.md)
- [日本語](/README.ja.md)
- [中文](/README.cn.md)
- [Español](/README.es.md)
- [Français](/README.fr.md)
- [Deutsch](/README.de.md)

# Supported Languages

| 语言   | 代码 |
|--------|------|
| 韩语   | `ko` |
| English | `en` |
| 日本语 | `ja` |
| 中文   | `cn` |
| 西班牙语 | `es` |
| 法语   | `fr` |
| 德语   | `de` |

# Installation

可以通过简单的命令进行安装。

```shell
go install github.com/YangTaeyoung/hugo-ai-translator@v1.0.0
```

# Configure

可以通过以下命令进行hugo-ai-translator的配置。

```shell
hugo-ai-translator configure
```

有关配置的更多信息，请参见[设置](docs/configure.cn.md)文档。

# Usage

## Simple Translation

可以用简单的规则翻译当前文件夹内的所有Markdown文件。

### Quick Start

```shell
cd path/to/markdown-directory

hugo-ai-translator simple --source-language en \
  --target-language ko \
  --target-language ja \
  --model gpt-4 \
  --api-key {open ai api key}
```  

## Rule Base Translation

可以应用特定规则进行翻译。

需要先进行[Configure](docs/configure.md)，翻译将根据设置的规则进行。

### Quick Start

如果有设置，则可以像下面这样翻译，无需使用其他选项。

```shell
hugo-ai-translator
```