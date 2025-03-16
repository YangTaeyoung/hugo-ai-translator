---
translated: true
---
![Go Version](https://img.shields.io/badge/Go-1.24-%23007d9c)
[![GoDoc](https://godoc.org/github.com/YangTaeyoung/hugo-ai-translator?status.svg)](https://pkg.go.dev/github.com/YangTaeyoung/hugo-ai-translator)
[![Go Test Action](https://github.com/YangTaeyoung/hugo-ai-translator/actions/workflows/test-ci.yaml/badge.svg)](https://github.com/YangTaeyoung/hugo-ai-translator/actions/workflows/test-ci.yaml)
![GitHub Release](https://img.shields.io/github/v/release/YangTaeyoung/hugo-ai-translator)
![GitHub License](https://img.shields.io/github/license/YangTaeyoung/hugo-ai-translator)
[![Go report](https://goreportcard.com/badge/github.com/YangTaeyoung/hugo-ai-translator)](https://goreportcard.com/report/github.com/YangTaeyoung/hugo-ai-translator)

# Hallo Hugo AI Translator! 👋

Dieses Dokument ist ein KI-Übersetzer, der Inhalte übersetzt, die in Hugo-Blogs gespeichert sind.

Es verwendet das Modell von [OpenAI](https://openai.com) und die übersetzten Ergebnisse werden gemäß der vom Benutzer festgelegten Namensregel gespeichert.

# Übersetzung

Diese Übersetzung wurde durch `hugo-ai-translator` erstellt.

- [Deutsch](/README.md)
- [English](/README.en.md)
- [日本語](/README.ja.md)
- [中文](/README.cn.md)
- [Español](/README.es.md)
- [Français](/README.fr.md)
- [Deutsch](/README.de.md)


# Unterstützte Sprachen

| Sprache   | Code |
|-----------|------|
| 한국어       | `ko` |
| English   | `en` |
| 日本語       | `ja` |
| 中文        | `cn` |
| Español   | `es` |
| Français  | `fr` |
| Deutsch   | `de` |

# Installation

Kann einfach mit einem Befehl installiert werden.

```shell
go install github.com/YangTaeyoung/hugo-ai-translator@v1.1.1
```

# Konfigurieren

Die folgenden Befehle können verwendet werden, um die Konfiguration von hugo-ai-translator auszuführen.

```shell
hugo-ai-translator configure
```

Für genauere Informationen zur Konfiguration siehe das [Konfigurationsdokument](docs/configure.md).

# Verwendung

## Einfache Übersetzung

Es können alle Markdown-Dateien im aktuellen Ordner mit einfachen Regeln übersetzt werden.

### Schneller Einstieg

```shell
cd path/to/markdown-directory

hugo-ai-translator simple --source-language en \
  --target-language ko \
  --target-languages all \
  --model gpt-4o-mini \
  --api-key {open ai api key}
``` 

## Regelbasierte Übersetzung

Es kann mit bestimmten Regeln übersetzt werden.

Die [Konfiguration](docs/configure.md) muss zuerst erfolgen, und die Übersetzung erfolgt gemäß den festgelegten Regeln.

### Schneller Einstieg

Wenn die Einstellungen vorhanden sind, ist es möglich, die Übersetzung auch ohne andere Optionen wie folgt durchzuführen.

```shell
hugo-ai-translator
```