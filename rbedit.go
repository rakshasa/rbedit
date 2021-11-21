package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/google/subcommands"
	"github.com/rakshasa/rbedit/commands"
)

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&commands.MapCmd{}, "")

	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "rbedit: %s\n", err.Error())

		exitErr, ok := err.(commands.ExitStatusError)
		if ok {
			os.Exit(int(exitErr.Status()))
		}

		os.Exit(int(subcommands.ExitFailure))
	}

	ctx := context.Background()

	exitCode := int(subcommands.Execute(ctx))

	os.Exit(exitCode)
}
