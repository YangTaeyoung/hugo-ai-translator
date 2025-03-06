---
translated: true
---
# Bonjour Hugo AI Translator! 👋

Ce document est un traducteur AI pour traduire du contenu stocké dans des blogs Hugo.

Il utilise le modèle de [OpenAI](https://openai.com), et les résultats traduits sont enregistrés selon les règles de nommage spécifiées par l'utilisateur.

# Traduction

Cette traduction a été réalisée via `hugo-ai-translator`.

- [한국어](/README.md)
- [English](/README.en.md)
- [日本語](/README.ja.md)
- [中文](/README.cn.md)
- [Español](/README.es.md)
- [Français](/README.fr.md)
- [Deutsch](/README.de.md)


# Langues prises en charge

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

Il peut être installé avec une simple commande.

```shell
go install github.com/YangTaeyoung/hugo-ai-translator@v1.0.1
```

# Configuration

Vous pouvez configurer hugo-ai-translator avec la commande suivante.

```shell
hugo-ai-translator configure
```

Pour plus de détails sur la configuration, veuillez consulter le document [configuration](docs/configure.fr.md).

# Utilisation

## Traduction simple

Il est possible de traduire tous les fichiers markdown dans le dossier actuel avec une règle simple.

### Démarrage rapide

```shell
cd path/to/markdown-directory

hugo-ai-translator simple --source-language en \
  --target-language ko \
  --target-language ja \
  --model gpt-4 \
  --api-key {clé api open ai}
```

## Traduction basées sur des règles

Vous pouvez traduire en appliquant des règles spécifiques.

La [configuration](docs/configure.md) doit être réalisée au préalable, et la traduction se déroulera selon les règles définies.

### Démarrage rapide

Si la configuration est faite, vous pouvez traduire sans utiliser d'autres options comme suit.

```shell
hugo-ai-translator
```