package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"
	"github.com/rakshasa/rbedit/commands"
)

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&commands.MapCmd{}, "")

	flag.Parse()
	ctx := context.Background()

	exitCode := int(subcommands.Execute(ctx))

	os.Exit(exitCode)
}
