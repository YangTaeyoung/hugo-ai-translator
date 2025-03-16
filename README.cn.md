---
translated: true
---
![Go版本](https://img.shields.io/badge/Go-1.24-%23007d9c)
[![GoDoc](https://godoc.org/github.com/YangTaeyoung/hugo-ai-translator?status.svg)](https://pkg.go.dev/github.com/YangTaeyoung/hugo-ai-translator)
[![Go测试行动](https://github.com/YangTaeyoung/hugo-ai-translator/actions/workflows/test-ci.yaml/badge.svg)](https://github.com/YangTaeyoung/hugo-ai-translator/actions/workflows/test-ci.yaml)
![GitHub发布](https://img.shields.io/github/v/release/YangTaeyoung/hugo-ai-translator)
![GitHub许可证](https://img.shields.io/github/license/YangTaeyoung/hugo-ai-translator)
[![Go报告](https://goreportcard.com/badge/github.com/YangTaeyoung/hugo-ai-translator)](https://goreportcard.com/report/github.com/YangTaeyoung/hugo-ai-translator)

# 嗨，Hugo AI翻译器！ 👋

该文档是用于翻译存储在Hugo博客中的内容的AI翻译器。

使用[OpenAI](https://openai.com)的模型，翻译结果将根据用户指定的命名规则进行保存。

# 翻译

该翻译版是通过 `hugo-ai-translator`进行翻译的。

- [한국어](/README.md)
- [English](/README.en.md)
- [日本語](/README.ja.md)
- [中文](/README.cn.md)
- [Español](/README.es.md)
- [Français](/README.fr.md)
- [Deutsch](/README.de.md)


# 支持的语言

| 语言   | 代码 |
|--------|------|
| 한국어    | `ko` |
| English | `en` |
| 日本語    | `ja` |
| 中文     | `cn` |
| Español | `es` |
| Français| `fr` |
| Deutsch | `de` |

# 安装

可以通过简单的命令进行安装。

```shell
go install github.com/YangTaeyoung/hugo-ai-translator@v1.1.1
```

# 配置

可以通过以下命令设置以运行hugо-ai-translator。

```shell
hugo-ai-translator configure
```

有关配置的更多详细信息，请参见[设置](docs/configure.md)文档。

# 使用

## 简单翻译

可以用简单的规则翻译当前文件夹内的所有Markdown。

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

可以应用特定规则进行翻译。

需要先[配置](docs/configure.md)，翻译将根据已设置的规则进行进行。

### 快速开始

如果有设置，则可以如下而无需使用其他选项进行翻译。

```shell
hugo-ai-translator
```  
