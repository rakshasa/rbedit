package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/google/subcommands"
	"github.com/rakshasa/rbedit/commands"
)

const (
	helperCategory     = "helper"
	bencodeCategory    = "bencode"
	bittorrentCategory = "bittorrent"
)

func main() {
	subcommands.Register(subcommands.HelpCommand(), helperCategory)
	subcommands.Register(subcommands.FlagsCommand(), helperCategory)
	subcommands.Register(subcommands.CommandsCommand(), helperCategory)
	subcommands.Register(&commands.GetCmd{}, bencodeCategory)
	subcommands.Register(&commands.PutCmd{}, bencodeCategory)
	subcommands.Register(&commands.AnnounceCmd{}, bittorrentCategory)
	subcommands.Register(&commands.AnnounceListCmd{}, bittorrentCategory)

	// TODO: Disable scientific notation and float, unless passed a flag.

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
