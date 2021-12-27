package main

import (
	"context"
	"fmt"
	"os"

	"github.com/rakshasa/rbedit/cmd/common"
)

func main() {
	ctx := context.Background()

	if err := common.NewRootCommand().ExecuteContext(ctx); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
