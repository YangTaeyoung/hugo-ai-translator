---
translated: true
---
![Version Go](https://img.shields.io/badge/Go-1.24-%23007d9c)
[![GoDoc](https://godoc.org/github.com/YangTaeyoung/hugo-ai-translator?status.svg)](https://pkg.go.dev/github.com/YangTaeyoung/hugo-ai-translator)
[![Action de Test Go](https://github.com/YangTaeyoung/hugo-ai-translator/actions/workflows/test-ci.yaml/badge.svg)](https://github.com/YangTaeyoung/hugo-ai-translator/actions/workflows/test-ci.yaml)
![Publication GitHub](https://img.shields.io/github/v/release/YangTaeyoung/hugo-ai-translator)
![Licence GitHub](https://img.shields.io/github/license/YangTaeyoung/hugo-ai-translator)
[![Rapport Go](https://goreportcard.com/badge/github.com/YangTaeyoung/hugo-ai-translator)](https://goreportcard.com/report/github.com/YangTaeyoung/hugo-ai-translator)

# Bonjour Hugo AI Translator! 👋

Ce document est un traducteur AI qui traduit le contenu enregistré dans les blogs Hugo.

Il utilise le modèle de [OpenAI](https://openai.com) et les résultats traduits sont enregistrés selon les règles de nommage spécifiées par l'utilisateur.

# Traduction

Cette traduction a été réalisée par `hugo-ai-translator`.

- [한국어](/README.md)
- [English](/README.en.md)
- [日本語](/README.ja.md)
- [中文](/README.cn.md)
- [Español](/README.es.md)
- [Français](/README.fr.md)
- [Deutsch](/README.de.md)


# Langues Supportées

| Langue     | Code |
|------------|------|
| 한국어        | `ko` |
| English    | `en` |
| 日本語        | `ja` |
| 中文         | `cn` |
| Español    | `es` |
| Français   | `fr` |
| Deutsch    | `de` |

# Installation

Vous pouvez l'installer avec de simples commandes.

```shell
 go install github.com/YangTaeyoung/hugo-ai-translator@v1.1.1
```

# Configuration

Vous pouvez configurer le hugo-ai-translator avec la commande suivante.

```shell
hugo-ai-translator configure
```

Pour plus de détails sur la configuration, veuillez consulter le document [Configuration](docs/configure.md).

# Utilisation

## Traduction Simple

Vous pouvez traduire tous les markdown dans le dossier courant avec des règles simples.

### Démarrage Rapide

```shell
cd path/to/markdown-directory

hugo-ai-translator simple --source-language en \
  --target-language ko \
  --target-languages all \
  --model gpt-4o-mini \
  --api-key {open ai api key}
``` 

## Traduction Basée sur des Règles

Vous pouvez appliquer des règles spécifiques pour traduire.

La [Configuration](docs/configure.md) doit être effectuée au préalable, et la traduction se fera selon les règles définies.

### Démarrage Rapide

Si la configuration est faite, vous pouvez traduire sans utiliser d'autres options comme suit.

```shell
hugo-ai-translator
```