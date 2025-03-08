---
translated: true
---
[![Go Test Action](https://github.com/YangTaeyoung/hugo-ai-translator/actions/workflows/test-ci.yaml/badge.svg)](https://github.com/YangTaeyoung/hugo-ai-translator/actions/workflows/test-ci.yaml)
![GitHub Release](https://img.shields.io/github/v/release/YangTaeyoung/hugo-ai-translator)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

# Bonjour Hugo AI Translator! 👋

Ce document est un traducteur AI qui traduit le contenu stocké dans les blogs Hugo.

Utilisant le modèle d'[OpenAI](https://openai.com), les résultats de la traduction sont stockés selon les règles de nomination spécifiées par l'utilisateur.

# Traduction

Cette traduction a été faite via `hugo-ai-translator`.

- [한국어](/README.md)
- [English](/README.en.md)
- [日本語](/README.ja.md)
- [中文](/README.cn.md)
- [Español](/README.es.md)
- [Français](/README.fr.md)
- [Deutsch](/README.de.md)


# Langues Supportées

| Langue   | Code |
|----------|------|
| 한국어      | `ko` |
| English  | `en` |
| 日本語      | `ja` |
| 中文       | `cn` |
| Español  | `es` |
| Français | `fr` |
| Deutsch  | `de` |

# Installation

Vous pouvez l'installer avec une simple commande.

```shell
go install github.com/YangTaeyoung/hugo-ai-translator@v1.1.0
```

# Configuration

Vous pouvez configurer le hugo-ai-translator en utilisant la commande suivante.

```shell
hugo-ai-translator configure
```

Pour plus de détails sur la configuration, veuillez consulter le document [Configuration](docs/configure.md).

# Utilisation

## Traduction Simple

Vous pouvez traduire tous les markdowns dans le dossier actuel avec des règles simples.

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

Vous pouvez traduire en appliquant des règles spécifiques.

La [configuration](docs/configure.md) doit être effectuée au préalable, et la traduction sera effectuée selon les règles définies.

### Démarrage Rapide

Avec la configuration en place, vous pouvez traduire sans utiliser d'autre option comme suit.

```shell
hugo-ai-translator
```