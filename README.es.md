---
translated: true
---
![Versión de Go](https://img.shields.io/badge/Go-1.24-%23007d9c)
[![GoDoc](https://godoc.org/github.com/YangTaeyoung/hugo-ai-translator?status.svg)](https://pkg.go.dev/github.com/YangTaeyoung/hugo-ai-translator)
[![Acción de prueba de Go](https://github.com/YangTaeyoung/hugo-ai-translator/actions/workflows/test-ci.yaml/badge.svg)](https://github.com/YangTaeyoung/hugo-ai-translator/actions/workflows/test-ci.yaml)
![Lanzamiento de GitHub](https://img.shields.io/github/v/release/YangTaeyoung/hugo-ai-translator)
![Licencia de GitHub](https://img.shields.io/github/license/YangTaeyoung/hugo-ai-translator)
[![Informe de Go](https://goreportcard.com/badge/github.com/YangTaeyoung/hugo-ai-translator)](https://goreportcard.com/report/github.com/YangTaeyoung/hugo-ai-translator)

# ¡Hola Hugo AI Translator! 👋

Este documento es un traductor AI que traduce el contenido almacenado en un blog de Hugo.

Utiliza el modelo de [OpenAI](https://openai.com), y los resultados traducidos se almacenan de acuerdo con las reglas de nomenclatura especificadas por el usuario.

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
|----------|------|
| 한국어      | `ko` |
| English  | `en` |
| 日本語      | `ja` |
| 中文       | `cn` |
| Español  | `es` |
| Français | `fr` |
| Deutsch  | `de` |

# Instalación

Se puede instalar con un simple comando.

```shell
go install github.com/YangTaeyoung/hugo-ai-translator@v1.1.1
```

# Configurar

Puede configurar hugo-ai-translator para que funcione a través del siguiente comando.

```shell
hugo-ai-translator configure
```

Para obtener más detalles sobre la configuración, consulte el documento de [configuración](docs/configure.md).

# Uso

## Traducción Simple

Puede traducir todos los markdown en la carpeta actual con una regla simple.

### Inicio Rápido

```shell
cd path/to/markdown-directory

hugo-ai-translator simple --source-language en \
  --target-language ko \
  --target-languages all \
  --model gpt-4o-mini \
  --api-key {clave de api de open ai}
``` 

## Traducción Basada en Reglas

Se puede traducir aplicando reglas específicas.

Es necesario realizar un [Configure](docs/configure.md) primero, y la traducción se llevará a cabo según las reglas establecidas.

### Inicio Rápido

Si hay una configuración, la traducción se puede realizar sin usar otras opciones de la siguiente manera.

```shell
hugo-ai-translator
```