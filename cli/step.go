package cli

import (
	"fmt"
	"slices"
	"strings"

	"github.com/YangTaeyoung/hugo-ai-translator/config"
	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
)

func openAIStep(cfg *config.Config) error {
	var (
		p   promptui.Prompt
		s   promptui.Select
		err error
	)
	fmt.Println("# OpenAI Setting")

	p = promptui.Prompt{
		Label: "Enter your OpenAI API key",
		Mask:  '*',
		Validate: func(s string) error {
			if strings.Trim(s, " ") == "" {
				return ErrEmptyInput
			}

			return nil
		},
	}

	cfg.OpenAI.ApiKey, err = p.Run()
	if err != nil {
		return err
	}

	s = promptui.Select{
		Label: "Select OpenAI model",
		Items: ChatModels,
	}

	_, cfg.OpenAI.Model, err = s.Run()
	if err != nil {
		return errors.Wrap(err, "failed to get openai model")
	}

	return nil
}

func contentDirStep(cfg *config.Config) error {
	var (
		p   promptui.Prompt
		err error
	)

	fmt.Println("# Content Directory Setting")

	p = promptui.Prompt{
		Label: "Enter the path to the content directory",
		Validate: func(s string) error {
			if strings.Trim(s, " ") == "" {
				return ErrEmptyInput
			}

			return nil
		},
	}

	cfg.Translator.ContentDir, err = p.Run()
	if err != nil {
		return err
	}

	return nil
}

func languageChoiceStep(cfg *config.Config) error {
	var (
		p promptui.Prompt
	)

	fmt.Println("# Language Chioce Setting")

	fmt.Println("Supported languages:")
	for code, name := range config.LanguageCodeToLanguage {
		fmt.Printf("%s (%s)\n", code, name)
	}
	fmt.Println()

	p = promptui.Prompt{
		Label: "Enter source language code",
		Validate: func(s string) error {
			if _, ok := config.LanguageCodeToLanguage[config.LanguageCode(s)]; !ok {
				return fmt.Errorf("unsupported language code: %s", s)
			}

			return nil
		},
	}

	sourceLanguage, err := p.Run()
	if err != nil {
		return err
	}

	cfg.Translator.Source.SourceLanguage = config.LanguageCode(sourceLanguage)

	p = promptui.Prompt{
		Label: "Enter target language code (please separate with comma if you want to use multiple languages) ex. ko, en",
		Validate: func(s string) error {
			s = strings.ReplaceAll(s, " ", "")
			for _, code := range strings.Split(s, ",") {
				if _, ok := config.LanguageCodeToLanguage[config.LanguageCode(code)]; !ok {
					return fmt.Errorf("unsupported language code: %s", code)
				}
			}

			return nil
		},
	}
	languageStr, err := p.Run()
	if err != nil {
		return err
	}

	for _, code := range strings.Split(strings.ReplaceAll(languageStr, " ", ""), ",") {
		cfg.Translator.Target.TargetLanguages = append(cfg.Translator.Target.TargetLanguages,
			config.LanguageCode(code),
		)
	}

	return nil
}

func ignoreRuleStep(cfg *config.Config) error {
	var p promptui.Prompt

	fmt.Println("# Ignore Rule Setting")

	p = promptui.Prompt{
		Label:   "Do you want to set the ignore rules for the translator? (y/n) default: n",
		Default: "n",
	}

	answer, err := p.Run()
	if err != nil {
		return err
	}

	if answer == "y" {
		for {
			p = promptui.Prompt{
				Label: `Enter ignore rules. (ex: *.md, ignoreDir/**) (type "exit" or blank to finish)`,
			}

			answer, err = p.Run()
			if err != nil {
				return err
			}

			if slices.Contains([]string{"exit", ""}, answer) {
				break
			}

			cfg.Translator.Source.IgnoreRules = append(cfg.Translator.Source.IgnoreRules, answer)
		}
	}

	return nil
}

// TODO: 문서화 작성 후 More Detail 링크 추가
const targetPathRuleDescription = `
!!!VERY IMPORTANT!!!
The target path rule is a rule that determines the path of the translated file.
You can use the following variables in the rule:
- {language}: the language code of the target language
- {origin}: the origin path of the source file directory path (Note: not include filename and extension)
- {filename}: the name of the source file

# Example
content directory: ~/hugo_root/content
target languages: [en, fr]
target path rule: "{origin}/{filename}.{language}.md"
current file system:
	~/hugo_root
	└── content
		└── _index.md
		└── some-posting.md
		└── some-posting2
			└── index.md


the translated file saved in:
- ~/hugo_root/content/_index.en.md
- ~/hugo_root/content/_index.fr.md
- ~/hugo_root/content/some-posting.md.en.md
- ~/hugo_root/content/some-posting.md.fr.md
- ~/hugo_root/content/some-posting2/index.en.md
- ~/hugo_root/content/some-posting2/index.fr.md
`

func targetPathRuleStep(cfg *config.Config) error {
	var (
		p   promptui.Prompt
		err error
	)

	fmt.Println("# Target Path Rule Setting")
	fmt.Println(targetPathRuleDescription)

	p = promptui.Prompt{
		Label: "Enter the target path rule for the translator",
		Validate: func(s string) error {
			if strings.Trim(s, " ") == "" {
				return ErrEmptyInput
			}

			return nil
		},
	}

	cfg.Translator.Target.TargetPathRule, err = p.Run()
	if err != nil {
		return err
	}

	return nil
}
