---
translated: true
---
# Hello Hugo AI Translator! 👋

Dieses Dokument ist ein KI-Übersetzer, der Inhalte speichert, die in Hugo Blogs gespeichert sind.

Es nutzt das Modell von [OpenAI](https://openai.com), und die übersetzten Ergebnisse werden gemäß der von den Benutzern angegebenen Namensregel gespeichert.

# Übersetzung

Diese Übersetzung wurde durch `hugo-ai-translator` erstellt.

- [한국어](/README.md)
- [English](/README.en.md)
- [日本語](/README.ja.md)
- [中文](/README.cn.md)
- [Español](/README.es.md)
- [Français](/README.fr.md)
- [Deutsch](/README.de.md)


# Unterstützte Sprachen

| Sprache   | Code |
|-----------|------|
| 한국어      | `ko` |
| English   | `en` |
| 日本語      | `ja` |
| 中文       | `cn` |
| Español   | `es` |
| Français  | `fr` |
| Deutsch   | `de` |

# Installation

Einfach mit einem Befehl zu installieren.

```shell
go install github.com/YangTaeyoung/hugo-ai-translator@v1.0.0
```

# Konfiguration

Die folgende Kommando kann verwendet werden, um die Konfiguration für den hugo-ai-translator festzulegen.

```shell
hugo-ai-translator configure
```

Weitere Informationen zur Konfiguration finden Sie im Dokument [Konfiguration](docs/configure.de.md).

# Nutzung

## Einfache Übersetzung

Alle Markdowns im aktuellen Ordner können durch einfache Regeln übersetzt werden.

### Schnellstart

```shell
cd path/to/markdown-directory

hugo-ai-translator simple --source-language en \
  --target-language ko \
  --target-language ja \
  --model gpt-4 \
  --api-key {open ai api key}
``` 

## Regelbasierte Übersetzung

Übersetzungen können unter Anwendung spezifischer Regeln durchgeführt werden.

[Konfigurieren](docs/configure.md) muss vorher durchgeführt werden, und die Übersetzung erfolgt gemäß der festgelegten Regeln.

### Schnellstart

Mit verfügbaren Einstellungen kann die Übersetzung auch ohne zusätzliche Optionen wie folgt durchgeführt werden.

```shell
hugo-ai-translator
```