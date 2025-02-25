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
				Name: "configure",
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
		},
		Action: TranslateAction,
	}
}
