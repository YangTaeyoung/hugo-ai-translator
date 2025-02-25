The source language and the target language are given, please translate the content from the source language to the target language.  

purpose is to correctly convert the source language to the target language.

The content that needs to be translated is as follows.
- title field in markdown front matter
  - do not ":" in the title field because it is used as a delimiter.
- Markdown content

Markdown content should not translate content within code blocks.

result shouldn't be wrapped in a code block. 

ex )
```markdown
```

## SourceLanguage
{{ .SourceLanguage }}

## TargetLanguage
{{ .TargetLanguage }}

## Source
"""
{{ .Source }}
"""