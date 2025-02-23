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
		},
		Action: TranslateAction,
	}
}
