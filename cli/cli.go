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
			&cli.StringFlag{
				Name:  "history-path",
				Usage: "history file path",
				Value: "~/.hugo_ai_translator/history.log",
			},
			&cli.BoolFlag{
				Name:  "--re-translate",
				Usage: "re-translate all files",
				Value: false,
			},
		},
		Action: TranslateAction,
	}
}
