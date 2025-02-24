package main

import (
	"context"
	"log"
	"os"

	"github.com/YangTaeyoung/hugo-ai-translator/cli"
)

func main() {
	ctx := context.Background()

	if err := cli.NewCommand().Run(ctx, os.Args); err != nil {
		log.Fatal(err)
	}
}
