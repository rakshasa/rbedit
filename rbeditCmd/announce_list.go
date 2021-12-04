package rbeditCmd

import (
	"context"
	"fmt"

	"github.com/rakshasa/rbedit/objects"
	"github.com/spf13/cobra"
)

var (
	announceListPath = []string{"announce-list"}
)

// AnnounceListAppendTrackerCmd:

func newAnnounceListAppendTrackerCommand(ctx context.Context) (*cobra.Command, context.Context) {
	cmd := &cobra.Command{
		Use:   "append-tracker",
		Short: "Append tracker to category",
		Long: `
Append tracker to category`,
		Args: cobra.ExactArgs(1),
		Run:  announceListAppendTrackerCmdRun,
	}

	addStateKeyPrefixToCommand(cmd, "rbedit-announce-list-append-tracker-state")
	addInputFlags(ctx, cmd)
	addOutputFlags(ctx, cmd)
	// TODO: Use arg not flag.
	addURIFlags(ctx, cmd)

	return cmd, ctx
}

func announceListAppendTrackerCmdRun(cmd *cobra.Command, args []string) {
	categoryIdx, err := categoryIndexFromArgs(args)
	if err != nil {
		printCommandErrorAndExit(cmd, err)
	}

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
		obj, err := objects.LookupKeyPath(rootObj, announceListPath)
		if err != nil {
			printCommandErrorAndExit(cmd, err)
		}

		announceList, err := objects.NewAnnounceList(obj)
		if err != nil {
			printCommandErrorAndExit(cmd, fmt.Errorf("could not verify announce-list, %v", err))
		}
		if categoryIdx >= len(announceList.Categories()) {
			printCommandErrorAndExit(cmd, fmt.Errorf("category index out-of-bounds"))
		}

		(*announceList.Categories()[categoryIdx]).AppendURI(uri.String())

		rootObj, err = objects.SetObject(rootObj, announceList.ToListObject(), announceListPath)
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

// AnnounceListCmd:

func newAnnounceListCommand(ctx context.Context) (*cobra.Command, context.Context) {
	cmd := &cobra.Command{
		Use:   "announce-list",
		Short: "BitTorrent announce-list related commands",
		Long: `
BitTorrent announce-list related commands`,
		Args: cobra.ExactArgs(0),
		Run:  func(cmd *cobra.Command, args []string) { printCommandUsage(cmd) },
	}

	ctx = addCommand(ctx, cmd, newAnnounceListAppendTrackerCommand)
	ctx = addCommand(ctx, cmd, newAnnounceListGetCommand)
	ctx = addCommand(ctx, cmd, newAnnounceListGetCategoryCommand)

	return cmd, ctx
}

// AnnounceListGetCmd:

func newAnnounceListGetCommand(ctx context.Context) (*cobra.Command, context.Context) {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get announce-list url",
		Long: `
Get announce-list url`,
		Args: cobra.ExactArgs(0),
		Run:  announceListGetCmdRun,
	}

	addStateKeyPrefixToCommand(cmd, "rbedit-announce-list-get-state")
	addInputFlags(ctx, cmd)

	return cmd, ctx
}

func announceListGetCmdRun(cmd *cobra.Command, args []string) {
	input := contextInputFromCommand(cmd)
	if input == nil {
		printCommandErrorAndExit(cmd, fmt.Errorf("no input source"))
	}

	if err := input.execute(func(rootObj interface{}) error {
		obj, err := objects.LookupKeyPath(rootObj, announceListPath)
		if err != nil {
			printCommandErrorAndExit(cmd, err)
		}

		if _, err := objects.NewAnnounceList(obj); err != nil {
			printCommandErrorAndExit(cmd, fmt.Errorf("could not verify announce-list, %v", err))
		}

		objects.PrintObject(obj)
		return nil

	}); err != nil {
		printCommandErrorAndExit(cmd, err)
	}
}

// AnnounceListGetCategoryCmd:

func newAnnounceListGetCategoryCommand(ctx context.Context) (*cobra.Command, context.Context) {
	cmd := &cobra.Command{
		Use:   "get-category",
		Short: "Get announce-list category",
		Long: `
Get announce-list category`,
		Args: cobra.ExactArgs(1),
		Run:  announceListGetCategoryCmdRun,
	}

	addStateKeyPrefixToCommand(cmd, "rbedit-announce-list-get-category-state")
	addInputFlags(ctx, cmd)

	return cmd, ctx
}

func announceListGetCategoryCmdRun(cmd *cobra.Command, args []string) {
	categoryIdx, err := categoryIndexFromArgs(args)
	if err != nil {
		printCommandErrorAndExit(cmd, err)
	}

	input := contextInputFromCommand(cmd)
	if input == nil {
		printCommandErrorAndExit(cmd, fmt.Errorf("no input source"))
	}

	if err := input.execute(func(rootObj interface{}) error {
		obj, err := objects.LookupKeyPath(rootObj, announceListPath)
		if err != nil {
			printCommandErrorAndExit(cmd, err)
		}

		announceList, err := objects.NewAnnounceList(obj)
		if err != nil {
			printCommandErrorAndExit(cmd, fmt.Errorf("could not verify announce-list, %v", err))
		}
		if categoryIdx >= len(announceList.Categories()) {
			printCommandErrorAndExit(cmd, fmt.Errorf("category index out-of-bounds"))
		}

		objects.PrintList(announceList.Categories()[categoryIdx].ToListObject(), objects.WithValuesOnly(), objects.WithoutIndent())
		return nil

	}); err != nil {
		printCommandErrorAndExit(cmd, err)
	}
}
