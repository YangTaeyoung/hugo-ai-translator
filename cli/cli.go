package cli

import (
	"github.com/urfave/cli/v3"
)

func NewCommand() *cli.Command {
	return &cli.Command{
		Name: "hugo-ai-translator",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config",
				Usage: "config file path",
				Value: "~/.hugo_ai_translator/config.yaml",
			},
			&cli.BoolFlag{
				Name:  "re-translate",
				Usage: "re-translate all files",
				Value: false,
			},
			&cli.BoolFlag{
				Name:   "debug",
				Usage:  "debug mode",
				Value:  false,
				Action: DebugModeAction,
			},
		},
		Commands: []*cli.Command{
			{
				Name:        "configure",
				Description: "configure hugo-ai-translator",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "config",
						Usage:   "config file path",
						Aliases: []string{"c"},
						Value:   "~/.hugo_ai_translator/config.yaml",
					},
					&cli.BoolFlag{
						Name:  "dry-run",
						Usage: "dry run",
						Value: false,
					},
					&cli.BoolFlag{
						Name:   "debug",
						Usage:  "debug mode",
						Value:  false,
						Action: DebugModeAction,
					},
				},
				Action: ConfigureAction,
			},
			{
				Name:        "simple",
				Description: "translate all markdown files in current directory\n",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "config",
						Usage:   "config file path",
						Aliases: []string{"c"},
						Value:   "~/.hugo_ai_translator/config.yaml",
					},
					&cli.StringFlag{
						Name:    "source-language",
						Usage:   "source language code (if don't set, it follows the config file's source language)",
						Aliases: []string{"s"},
					},
					&cli.StringSliceFlag{
						Name:    "target-languages",
						Usage:   "destination language codes (if don't set, it follows the config file's target languages)",
						Aliases: []string{"t"},
					},
					&cli.StringFlag{
						Name:    "model",
						Usage:   "OpenAI model name. you can check available models in https://platform.openai.com/docs/models#current-model-aliases \n(if don't set, follow the config file's model)",
						Aliases: []string{"m"},
					},
					&cli.StringFlag{
						Name:    "api-key",
						Usage:   "OpenAI API Key (if don't set, it follows the config file's api key)",
						Aliases: []string{"k"},
					},
				},
				Action: SimpleTranslateAction,
			},
		},
		Action: TranslateAction,
	}
}
