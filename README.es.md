---
translated: true
---
# ¡Hola Traductor de AI de Hugo! 👋

Este documento es un traductor de AI que traduce contenido almacenado en blogs de Hugo.

Utiliza el modelo de [OpenAI](https://openai.com), y los resultados traducidos se guardan de acuerdo con las reglas de nomenclatura especificadas por el usuario.

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

| Idioma    | Código |
|-----------|-------|
| 한국어     | `ko`  |
| English   | `en`  |
| 日本語     | `ja`  |
| 中文      | `cn`  |
| Español   | `es`  |
| Français  | `fr`  |
| Deutsch   | `de`  |

# Instalación

Se puede instalar con un simple comando.

```shell
go install github.com/YangTaeyoung/hugo-ai-translator@v1.0.0
```

# Configuración

Se puede configurar hugo-ai-translator ejecutando el siguiente comando.

```shell
hugo-ai-translator configure
```

Por favor, consulta el documento de [configuración](docs/configure.es.md) para obtener más detalles sobre la configuración.

# Uso

## Traducción Simple

Se pueden traducir todos los markdown en la carpeta actual con una regla simple.

### Inicio Rápido

```shell
cd path/to/markdown-directory

hugo-ai-translator simple --source-language en \
  --target-language ko \
  --target-language ja \
  --model gpt-4 \
  --api-key {clave de API de open ai}
``` 

## Traducción Basada en Reglas

Se puede traducir aplicando reglas específicas.

Primero se debe ejecutar [Configurar](docs/configure.md), y la traducción se llevará a cabo de acuerdo con las reglas configuradas.

### Inicio Rápido

Si ya hay configuración, se puede traducir sin utilizar otras opciones de la siguiente manera.

```shell
hugo-ai-translator
```