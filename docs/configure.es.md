---
translated: true
---
# Configuración
Esta sección proporciona una descripción del archivo de configuración del Hugo AI Translator.

## Ubicación
El archivo de configuración de Hugo AI Translator se ubica en `~/.hugo-ai-translator/config.yaml`.

```shell
cat ~/.hugo-ai-translator/config.yaml
```

## Esquema  
```yaml
openai:
    model: gpt-4o-mini
    api_key: {your-openai-api-key}
translator:
    content_dir: ~/dev/personal/YangTaeyoung.github.io/content
    source:
        source_language: ko
        ignore_rules: []
    target:
        target_languages:
            - en
            - ja
            - fr
            - de
        target_path_rule: '{origin}/{fileName}.{language}.md'
```

## `openai`
- `model`: Especifica el modelo que se utilizará en la API de OpenAI. Para más información sobre el modelo, consulta [Modelos de Open AI](https://platform.openai.com/docs/models).
- `api_key`: Especifica la clave API para utilizar la API de OpenAI.

## `translator`
- `content_dir`: Selecciona el directorio de contenido de Hugo. Actúa como la ruta raíz donde se guardarán los resultados traducidos.
- `source`
    - `source_language`: Especifica el idioma original del markdown a traducir.
    - `ignore_rules`: Especifica los archivos markdown que no serán traducidos. Puedes usar comodines como *, **.
      - ej) `ignore_rules: ["*.en.md", "*.ko.md"]`, `ignore_rules: ["some/path/**"]`
- `target`
  - `target_languages`: Especifica los idiomas a traducir. Puedes especificar varios idiomas. Consulta los idiomas admitidos en [Idiomas Soportados](../README.md#supported-languages).
  - ej) `target_languages: ["en", "ja", "fr", "de"]`

### `translator.target_path_rule`
Especifica la ruta donde se guardarán los resultados traducidos. Puedes usar palabras reservadas como `{origin}`, `{fileName}`, `{language}`.
- `{origin}`: Se refiere a la ruta del directorio del archivo original desde `translator.content_dir`. En el caso de `~/dev/personal/YangTaeyoung.github.io/content/some/index.md`, `~/dev/personal/YangTaeyoung.github.io/content` es `content_dir` y `some` es `origin`.
- `{fileName}`: Se refiere al nombre del archivo sin la extensión.
- `{language}`: Se refiere al código del idioma al que se traducirá.

#### Ejemplo
Los valores de configuración son los siguientes:
```yaml
target:
    content_dir: /content
    target_path_rule: '{origin}/{fileName}.{language}.md'
    target_languages:
        - en
        - ja
```

La estructura de archivos antes de la traducción es:
```
/content
    /some
        index.md
    index.md
```

La estructura de archivos después de la traducción es:
```
/content
    /some
        hamburger.md
        hamburger.en.md <-- Archivo traducido
        hamburger.ja.md <-- Archivo traducido
    index.md
    index.en.md <-- Archivo traducido
    index.ja.md <-- Archivo traducido
```