package rbeditCmd

import (
	"context"
	"fmt"

	"github.com/rakshasa/rbedit/objects"
	"github.com/spf13/cobra"
)

var (
	announcePath = []string{"announce"}
)

// AnnounceCmd:

func newAnnounceCommand(ctx context.Context) (*cobra.Command, context.Context) {
	cmd := &cobra.Command{
		Use:   "announce",
		Short: "BitTorrent announce related commands",
		Long: `
BitTorrent announce related commands`,
		Args: cobra.ExactArgs(0),
		Run:  func(cmd *cobra.Command, args []string) { printCommandUsage(cmd) },
	}

	ctx = addCommand(ctx, cmd, newAnnounceGetCommand)
	ctx = addCommand(ctx, cmd, newAnnouncePutCommand)

	return cmd, ctx
}

// AnnounceGetCmd:

func newAnnounceGetCommand(ctx context.Context) (*cobra.Command, context.Context) {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get announce url",
		Long: `
Get announce url`,
		Args: cobra.ExactArgs(0),
		Run:  announceGetCmdRun,
	}

	addStateKeyPrefixToCommand(cmd, "rbedit-announce-get-state")
	addInputFlags(ctx, cmd)

	return cmd, ctx
}

func announceGetCmdRun(cmd *cobra.Command, args []string) {
	input := contextInputFromCommand(cmd)
	if input == nil {
		printCommandErrorAndExit(cmd, fmt.Errorf("no input source"))
	}

	if err := input.execute(func(rootObj interface{}) error {
		obj, err := objects.LookupKeyPath(rootObj, announcePath)
		if err != nil {
			printCommandErrorAndExit(cmd, err)
		}

		if _, ok := objects.AsAbsoluteURI(obj); !ok {
			printCommandErrorAndExit(cmd, fmt.Errorf("announce not a valid URI string"))
		}

		objects.PrintObject(obj)
		return nil

	}); err != nil {
		printCommandErrorAndExit(cmd, err)
	}
}

// AnnouncesPutCmd:

// TODO: Pass uri as arg.

func newAnnouncePutCommand(ctx context.Context) (*cobra.Command, context.Context) {
	cmd := &cobra.Command{
		Use:   "put",
		Short: "Put announce url",
		Long: `
Put announce url`,
		Args: cobra.ExactArgs(0),
		Run:  announcePutCmdRun,
	}

	addStateKeyPrefixToCommand(cmd, "rbedit-announce-put-state")
	addInputFlags(ctx, cmd)
	addOutputFlags(ctx, cmd)
	addURIFlags(ctx, cmd)

	return cmd, ctx
}

func announcePutCmdRun(cmd *cobra.Command, args []string) {
	input := contextInputFromCommand(cmd)
	if input == nil {
		printCommandErrorAndExit(cmd, fmt.Errorf("no input source"))
	}
	output := contextOutputFromCommand(cmd)
	if output == nil {
		printCommandErrorAndExit(cmd, fmt.Errorf("no output target"))
	}
	uri := contextURIFromCommand(cmd)
	if output == nil {
		printCommandErrorAndExit(cmd, fmt.Errorf("no tracker URI"))
	}

	if err := input.execute(func(rootObj interface{}) error {
		rootObj, err := objects.SetObject(rootObj, uri.String(), announcePath)
		if err != nil {
			printCommandErrorAndExit(cmd, err)
		}

		if err := output.execute(rootObj, input.filePath); err != nil {
			printCommandErrorAndExit(cmd, err)
		}

		return nil

	}); err != nil {
		printCommandErrorAndExit(cmd, err)
	}
}
