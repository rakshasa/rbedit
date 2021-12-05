package rbeditCmd

import (
	"fmt"

	"github.com/rakshasa/rbedit/objects"
	"github.com/spf13/cobra"
)

var (
	announceListPath = []string{"announce-list"}
)

// AnnounceListCmd:

func newAnnounceListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "announce-list COMMAND",
		Short: "BitTorrent announce-list related commands",
		Args:  cobra.ExactArgs(0),
		Run:   func(cmd *cobra.Command, args []string) { printCommandUsage(cmd) },
	}

	setupDefaultCommand(cmd)

	cmd.AddCommand(newAnnounceListAppendTrackerCommand())
	cmd.AddCommand(newAnnounceListGetCommand())
	cmd.AddCommand(newAnnounceListGetCategoryCommand())

	return cmd
}

// AnnounceListAppendTrackerCmd:

func newAnnounceListAppendTrackerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "append-tracker [OPTIONS] INDEX URI...",
		Short: "Append tracker(s) to a category index",
		Args:  cobra.MinimumNArgs(2),
		Run:   announceListAppendTrackerCmdRun,
	}

	setupDefaultCommand(cmd)

	addInputFlags(cmd)
	addOutputFlags(cmd)

	return cmd
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

	metadata, err := metadataFromCommand(cmd, WithInput(), WithOutput())
	if err != nil {
		printCommandErrorAndExit(cmd, err)
	}

	input := objects.NewSingleInput(objects.NewDecodeBencode(), objects.NewFileInput())
	output := objects.NewSingleOutput(objects.NewEncodeBencode(), objects.NewFileOutput())

	if err := input.Execute(metadata, func(rootObj interface{}, metadata objects.IOMetadata) error {
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
		if err := output.Execute(rootObj, metadata); err != nil {
			printCommandErrorAndExit(cmd, err)
		}

		return nil

	}); err != nil {
		printCommandErrorAndExit(cmd, err)
	}
}

// AnnounceListGetCmd:

func newAnnounceListGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [OPTIONS]",
		Short: "Get announce-list url",
		Args:  cobra.ExactArgs(0),
		Run:   announceListGetCmdRun,
	}

	setupDefaultCommand(cmd)

	addInputFlags(cmd)

	return cmd
}

func announceListGetCmdRun(cmd *cobra.Command, args []string) {
	metadata, err := metadataFromCommand(cmd, WithInput())
	if err != nil {
		printCommandErrorAndExit(cmd, err)
	}

	input := objects.NewSingleInput(objects.NewDecodeBencode(), objects.NewFileInput())

	if err := input.Execute(metadata, func(rootObj interface{}, metadata objects.IOMetadata) error {
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

func newAnnounceListGetCategoryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-category [OPTIONS] INDEX",
		Short: "Get announce-list category",
		Args:  cobra.ExactArgs(1),
		Run:   announceListGetCategoryCmdRun,
	}

	setupDefaultCommand(cmd)

	addInputFlags(cmd)

	return cmd
}

func announceListGetCategoryCmdRun(cmd *cobra.Command, args []string) {
	categoryIdx, err := categoryIndexFromArgs(args)
	if err != nil {
		printCommandErrorAndExit(cmd, err)
	}

	metadata, err := metadataFromCommand(cmd, WithInput())
	if err != nil {
		printCommandErrorAndExit(cmd, err)
	}

	input := objects.NewSingleInput(objects.NewDecodeBencode(), objects.NewFileInput())

	if err := input.Execute(metadata, func(rootObj interface{}, metadata objects.IOMetadata) error {
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
