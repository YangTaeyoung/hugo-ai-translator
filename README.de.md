---
translated: true
---
[![Go Test Action](https://github.com/YangTaeyoung/hugo-ai-translator/actions/workflows/test-ci.yaml/badge.svg)](https://github.com/YangTaeyoung/hugo-ai-translator/actions/workflows/test-ci.yaml)  
![GitHub Release](https://img.shields.io/github/v/release/YangTaeyoung/hugo-ai-translator)  
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)  

# Hallo Hugo AI Übersetzer! 👋  

Dieses Dokument ist ein KI-Übersetzer, der Inhalte übersetzt, die in Hugo-Blogs gespeichert sind.  

Es verwendet das Modell von [OpenAI](https://openai.com) und die übersetzten Ergebnisse werden gemäß der vom Benutzer angegebenen Namensregel gespeichert.  

# Übersetzung  

Diese Übersetzung wurde über `hugo-ai-translator` erstellt.  

- [Koreanisch](/README.md)  
- [English](/README.en.md)  
- [日本語](/README.ja.md)  
- [中文](/README.cn.md)  
- [Español](/README.es.md)  
- [Français](/README.fr.md)  
- [Deutsch](/README.de.md)  


# Unterstützte Sprachen  

| Sprache   | Code |  
|-----------|------|  
| 한국어        | `ko` |  
| English   | `en` |  
| 日本語        | `ja` |  
| 中文         | `cn` |  
| Español   | `es` |  
| Français  | `fr` |  
| Deutsch   | `de` |  

# Installation  

Sie können es mit einem einfachen Befehl installieren.  

```shell  
go install github.com/YangTaeyoung/hugo-ai-translator@v1.1.0  
```  

# Konfigurieren  

Sie können die Konfiguration für den hugo-ai-translator mit dem folgenden Befehl vornehmen.  

```shell  
hugo-ai-translator configure  
```  

Für detailliertere Informationen zur Konfiguration werfen Sie bitte einen Blick in das [Konfigurationsdokument](docs/configure.md).  

# Nutzung  

## Einfache Übersetzung  

Sie können alle Markdown-Dateien im aktuellen Ordner mit einfachen Regeln übersetzen.  

### Schnelleinstieg  

```shell  
cd path/to/markdown-directory  
  
hugo-ai-translator simple --source-language en \  
  --target-language ko \  
  --target-languages all \  
  --model gpt-4o-mini \  
  --api-key {open ai api key}  
```  

## Regelbasierte Übersetzung  

Sie können Übersetzungen mithilfe spezifischer Regeln durchführen.  

Es ist erforderlich, [Configure](docs/configure.md) vorher durchzuführen, und die Übersetzung erfolgt gemäß den festgelegten Regeln.  

### Schnelleinstieg  

Wenn die Konfiguration vorhanden ist, können Sie die Übersetzung auch ohne Verwendung anderer Optionen wie folgt durchführen.  

```shell  
hugo-ai-translator  
```  
