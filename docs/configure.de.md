---
translated: true
---
# Konfiguration
In diesem Abschnitt wird die Beschreibung der Konfigurationsdatei des Hugo AI Translators bereitgestellt.

## Speicherort
Die Konfigurationsdatei des Hugo AI Translators befindet sich unter `~/.hugo-ai-translator/config.yaml`.

```shell
cat ~/.hugo-ai-translator/config.yaml
```

## Schema 
```yaml
openai:
    model: gpt-4o-mini
    api_key: {dein-openai-api-schlüssel}
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
- `model`: Gibt das Modell an, das über die OpenAI API verwendet werden soll. Weitere Informationen zum Modell finden Sie unter [Open AI Models](https://platform.openai.com/docs/models).
- `api_key`: Gibt den API-Schlüssel an, der zum Nutzen der OpenAI API erforderlich ist.

## `translator`
- `content_dir`: Wählt das Inhaltsverzeichnis von Hugo aus. Es dient als Wurzelverzeichnis, in dem die übersetzten Ergebnisse gespeichert werden.
- `source`
    - `source_language`: Gibt die Quellsprache des zu übersetzenden Markdowns an.
    - `ignore_rules`: Gibt Markdown-Dateien an, die nicht übersetzt werden sollen. Platzhalter wie *, ** können verwendet werden.
      - z.B. `ignore_rules: ["*.en.md", "*.ko.md"]`, `ignore_rules: ["some/path/**"]`
- `target`
  - `target_languages`: Gibt die Sprachen an, in die übersetzt werden soll. Es können mehrere Sprachen festgelegt werden. Unterstützte Sprachen finden Sie unter [Supported Languages](../README.md#supported-languages).
  - z.B. `target_languages: ["en", "ja", "fr", "de"]`

### `translator.target_path_rule`
Gibt den Pfad an, unter dem die übersetzten Ergebnisse gespeichert werden. Es können Platzhalter wie `{origin}`, `{fileName}`, `{language}` verwendet werden.
- `{origin}`: Bezieht sich auf den Verzeichnispfad der Quelldatei ab `translator.content_dir`. Im Fall von `~/dev/personal/YangTaeyoung.github.io/content/some/index.md` wird `~/dev/personal/YangTaeyoung.github.io/content` zu `content_dir`, und `some` wird zu `origin`.
- `{fileName}`: Bezieht sich auf den Dateinamen ohne Erweiterung.
- `{language}`: Bezieht sich auf den Code der zu übersetzenden Sprache.

#### Beispiel
Die Konfigurationswerte sind wie folgt:
```yaml
target:
    content_dir: /content
    target_path_rule: '{origin}/{fileName}.{language}.md'
    target_languages:
        - en
        - ja
```

Die Dateistruktur vor der Übersetzung ist wie folgt aufgebaut:
```
/content
    /some
        index.md
    index.md
```

Die Dateistruktur nach der Übersetzung ist wie folgt aufgebaut:
```
/content
    /some
        hamburger.md
        hamburger.en.md <-- übersetzte Datei
        hamburger.ja.md <-- übersetzte Datei
    index.md
    index.en.md <-- übersetzte Datei
    index.ja.md <-- übersetzte Datei
```
