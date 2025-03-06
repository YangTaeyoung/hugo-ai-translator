---
translated: true
---
# Configuration
Cette section décrit le fichier de configuration de Hugo AI Translator.

## Location
Le fichier de configuration de Hugo AI Translator se trouve à `~/.hugo-ai-translator/config.yaml`.

```shell
cat ~/.hugo-ai-translator/config.yaml
```

## Schema 
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
- `model` : Spécifie le modèle à utiliser avec l'API OpenAI. Pour plus d'informations sur les modèles, veuillez consulter [Open AI Models](https://platform.openai.com/docs/models).
- `api_key` : Spécifie la clé API nécessaire pour utiliser l'API OpenAI.

## `translator`
- `content_dir` : Sélectionnez le répertoire de contenu de Hugo. Il sert de chemin racine où les résultats traduits seront enregistrés.
- `source`
    - `source_language` : Indique la langue d'origine du markdown à traduire.
    - `ignore_rules` : Spécifie les fichiers markdown à ignorer dans la traduction. Des caractères génériques comme *, ** peuvent être utilisés.
      - ex) `ignore_rules: ["*.en.md", "*.ko.md"]`, `ignore_rules: ["some/path/**"]`
- `target`
  - `target_languages` : Indique les langues vers lesquelles traduire. Plusieurs langues peuvent être spécifiées. Veuillez consulter [Supported Languages](../README.md#supported-languages) pour les langues prises en charge.
  - ex) `target_languages: ["en", "ja", "fr", "de"]`

### `translator.target_path_rule`
Spécifie le chemin où le résultat traduit sera enregistré. Les mots réservés `{origin}`, `{fileName}`, `{language}` peuvent être utilisés.
- `{origin}` : Cela fait référence au chemin du répertoire du fichier source à partir de `translator.content_dir`. Pour `~/dev/personal/YangTaeyoung.github.io/content/some/index.md`, `~/dev/personal/YangTaeyoung.github.io/content` est `content_dir` et `some` est `origin`.
- `{fileName}` : Fait référence au nom de fichier sans extension.
- `{language}` : Fait référence au code de la langue dans laquelle la traduction sera faite.

#### Exemple
Les valeurs de configuration sont comme suit.
```yaml
target:
    content_dir: /content
    target_path_rule: '{origin}/{fileName}.{language}.md'
    target_languages:
        - en
        - ja
```

La structure des fichiers avant traduction est comme suit :
```
/content
    /some
        index.md
    index.md
```

La structure des fichiers après traduction est comme suit :
```
/content
    /some
        hamburger.md
        hamburger.en.md <-- Fichier traduit
        hamburger.ja.md <-- Fichier traduit
    index.md
    index.en.md <-- Fichier traduit
    index.ja.md <-- Fichier traduit
```