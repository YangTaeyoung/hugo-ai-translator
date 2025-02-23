package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/YangTaeyoung/hugo-ai-translator/cli"
)

func main() {
	ctx := context.Background()

	slog.SetLogLoggerLevel(slog.LevelDebug)

	if err := cli.NewCommand().Run(ctx, os.Args); err != nil {
		log.Fatal(err)
	}
}
