package main

import (
	"context"
	"os"

	"github.com/rakshasa/rbedit/cmd/common"
)

func main() {
	ctx := context.Background()

	if err := common.NewRootCommand().ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
