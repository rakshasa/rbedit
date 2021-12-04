package main

import (
	"context"

	cmd "github.com/rakshasa/rbedit/rbeditCmd"
)

func main() {
	cmd.Execute(context.Background())
}
