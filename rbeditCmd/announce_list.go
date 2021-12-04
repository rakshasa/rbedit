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

// AnnounceListCmd:

func newAnnounceListCommand(ctx context.Context) (*cobra.Command, context.Context) {
	cmd := &cobra.Command{
		Use:   "announce-list COMMAND",
		Short: "BitTorrent announce-list related commands",
		Args:  cobra.ExactArgs(0),
		Run:   func(cmd *cobra.Command, args []string) { printCommandUsage(cmd) },
	}

	setupDefaultCommand(cmd, "rbedit-announce-list")

	ctx = addCommand(ctx, cmd, newAnnounceListAppendTrackerCommand)
	ctx = addCommand(ctx, cmd, newAnnounceListGetCommand)
	ctx = addCommand(ctx, cmd, newAnnounceListGetCategoryCommand)

	return cmd, ctx
}

// AnnounceListAppendTrackerCmd:

func newAnnounceListAppendTrackerCommand(ctx context.Context) (*cobra.Command, context.Context) {
	cmd := &cobra.Command{
		Use:   "append-tracker [OPTIONS] INDEX URI...",
		Short: "Append tracker(s) to a category index",
		Args:  cobra.MinimumNArgs(2),
		Run:   announceListAppendTrackerCmdRun,
	}

	setupDefaultCommand(cmd, "rbedit-announce-list-append-tracker-state")

	addInputFlags(ctx, cmd)
	addOutputFlags(ctx, cmd)

	return cmd, ctx
}

func announceListAppendTrackerCmdRun(cmd *cobra.Command, args []string) {
	categoryIdx, err := categoryIndexFromArgs(args[:1])
	if err != nil {
		printCommandErrorAndExit(cmd, err)
	}
	trackers := []string{}
	for idx, t := range args[1:] {
		if !objects.VerifyAbsoluteURI(t) {
			printCommandErrorAndExit(cmd, fmt.Errorf("failed to validate URI for tracker %d\n", idx))
		}
		trackers = append(trackers, t)
	}
	input := contextInputFromCommand(cmd)
	if input == nil {
		printCommandErrorAndExit(cmd, fmt.Errorf("no input source"))
	}
	output := contextOutputFromCommand(cmd)
	if output == nil {
		printCommandErrorAndExit(cmd, fmt.Errorf("no output target"))
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

		for _, t := range trackers {
			(*announceList.Categories()[categoryIdx]).AppendURI(t)
		}

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

// AnnounceListGetCmd:

func newAnnounceListGetCommand(ctx context.Context) (*cobra.Command, context.Context) {
	cmd := &cobra.Command{
		Use:   "get [OPTIONS]",
		Short: "Get announce-list url",
		Args:  cobra.ExactArgs(0),
		Run:   announceListGetCmdRun,
	}

	setupDefaultCommand(cmd, "rbedit-announce-list-get-state")

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
		Use:   "get-category [OPTIONS] INDEX",
		Short: "Get announce-list category",
		Args:  cobra.ExactArgs(1),
		Run:   announceListGetCategoryCmdRun,
	}

	setupDefaultCommand(cmd, "rbedit-announce-list-get-category-state")

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
