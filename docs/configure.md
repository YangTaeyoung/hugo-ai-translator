# Configuration
해당 섹션에서는 Hugo AI Translator의 설정 파일에 대한 설명을 제공합니다.

## Location
Hugo AI Translator의 설정 파일은 `~/.hugo-ai-translator/config.yaml`에 위치해있습니다.

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
- `model`: OpenAI API에서 사용할 모델을 지정합니다 모델에 대한 정보는 [Open AI Models](https://platform.openai.com/docs/models)를 참고해주세요 
- `api_key`: OpenAI API를 사용하기 위한 API 키를 지정합니다 

## `translator`
- `content_dir`: Hugo의 컨텐츠 디렉토리를 선택합니다. 번역된 결과가 저장될 루트 경로로의 역할을 수행합니다.
- `source`
    - `source_language`: 번역할 마크다운의 원본 언어를 지정합니다.
    - `ignore_rules`: 번역하지 않을 마크다운 파일을 지정합니다. *, ** 등의 와일드카드를 사용할 수 있습니다.
      - ex) `ignore_rules: ["*.en.md", "*.ko.md"]`, `ignore_rules: ["some/path/**"]`
- `target`
  - `target_languages`: 번역할 언어를 지정합니다. 여러 언어를 지정할 수 있습니다. 지원 언어는 [Supported Languages](../README.md#supported-languages)를 참고해주세요.
  - ex) `target_languages: ["en", "ja", "fr", "de"]`

### `translator.target_path_rule`
번역된 결과가 저장될 경로를 지정합니다. `{origin}`, `{fileName}`, `{language}`의 예약어를 활용할 수 있습니다.
- `{origin}`:`translator.content_dir`부터의 원본 파일의 디렉토리 경로를 의미합니다. `~/dev/personal/YangTaeyoung.github.io/content/some/index.md`의 경우, `~/dev/personal/YangTaeyoung.github.io/content`가 `content_dir`, `some`이 `origin`이 됩니다. 
- `{fileName}`: 확장자를 제외한 파일 이름을 의미합니다.
- `{language}`: 번역될 언어의 코드를 의미합니다.

#### 예시
설정값은 다음과 같습니다.
```yaml
target:
    content_dir: /content
    target_path_rule: '{origin}/{fileName}.{language}.md'
    target_languages:
        - en
        - ja
```

번역 전 파일 스키마는 다음과 같이 구성되어 있습니다.

```
/content
    /some
        index.md
    index.md
```

번역 후 파일 스키마는 다음과 같이 구성됩니다.
```
/content
    /some
        hamburger.md
        hamburger.en.md <-- 번역된 파일
        hamburger.ja.md <-- 번역된 파일
    index.md
    index.en.md <-- 번역된 파일
    index.ja.md <-- 번역된 파일
```
