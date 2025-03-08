---
translated: true
---
[![Go Test Action](https://github.com/YangTaeyoung/hugo-ai-translator/actions/workflows/test-ci.yaml/badge.svg)](https://github.com/YangTaeyoung/hugo-ai-translator/actions/workflows/test-ci.yaml)
![GitHub Release](https://img.shields.io/github/v/release/YangTaeyoung/hugo-ai-translator)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

# ¡Hola Hugo AI Translator! 👋

Este documento es un traductor AI que traduce el contenido almacenado en el blog de Hugo.

Se utiliza el modelo de [OpenAI](https://openai.com) y los resultados traducidos se guardan según las reglas de nomenclatura especificadas por el usuario.

# Traducción

Esta traducción fue realizada a través de `hugo-ai-translator`.

- [한국어](/README.md)
- [English](/README.en.md)
- [日本語](/README.ja.md)
- [中文](/README.cn.md)
- [Español](/README.es.md)
- [Français](/README.fr.md)
- [Deutsch](/README.de.md)


# Idiomas Soportados

| Idioma   | Código |
|----------|-------|
| 한국어      | `ko`  |
| English  | `en`  |
| 日本語      | `ja`  |
| 中文       | `cn`  |
| Español  | `es`  |
| Français | `fr`  |
| Deutsch  | `de`  |

# Instalación

Se puede instalar con un simple comando.

```shell
go install github.com/YangTaeyoung/hugo-ai-translator@v1.1.0
```

# Configuraciones

Se pueden realizar configuraciones para ejecutar hugo-ai-translator mediante el siguiente comando.

```shell
hugo-ai-translator configure
```

Para más detalles sobre la configuración, consulta el documento de [configuración](docs/configure.md).

# Uso

## Traducción Simple

Puedes traducir todos los markdown en la carpeta actual utilizando reglas simples.

### Inicio Rápido

```shell
cd path/to/markdown-directory

hugo-ai-translator simple --source-language en \
  --target-language ko \
  --target-languages all \
  --model gpt-4o-mini \
  --api-key {open ai api key}
```

## Traducción Basada en Reglas

Se puede traducir aplicando reglas específicas.

Es necesario realizar [Configurar](docs/configure.md) previamente y la traducción se realizará según las reglas establecidas.

### Inicio Rápido

Si hay configuraciones, se puede traducir sin usar otras opciones como a continuación:

```shell
hugo-ai-translator
```