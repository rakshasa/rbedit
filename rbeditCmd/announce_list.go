package rbeditCmd

import (
	"fmt"

	"github.com/rakshasa/rbedit/actions"
	"github.com/rakshasa/rbedit/inputs"
	"github.com/rakshasa/rbedit/objects"
	"github.com/rakshasa/rbedit/outputs"
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
		Run:   func(cmd *cobra.Command, args []string) { printCommandUsageAndExit(cmd) },
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

	input := inputs.NewSingleInput(inputs.NewDecodeBencode(), inputs.NewFileInput())
	output := outputs.NewSingleOutput(outputs.NewEncodeBencode(), outputs.NewFileOutput())

	if err := input.Execute(metadata, actions.NewGetAnnounceListAppendTrackerAction(output, categoryIdx, trackers)); err != nil {
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

	input := inputs.NewSingleInput(inputs.NewDecodeBencode(), inputs.NewFileInput())
	output := outputs.NewSingleOutput(outputs.NewEncodePrintAsListOfLists(), outputs.NewStdOutput())

	if err := input.Execute(metadata, actions.NewGetAnnounceListAction(output, announceListPath)); err != nil {
		printCommandErrorAndExit(cmd, err)
	}
}

// AnnounceListGetCategoryCmd:

func newAnnounceListGetCategoryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-category [OPTIONS] INDEX",
		Short: "Get announce-list category",
		// TODO: Create a function that verifies valid positive index.
		Args: cobra.ExactArgs(1),
		Run:  announceListGetCategoryCmdRun,
	}

	setupDefaultCommand(cmd)

	addInputFlags(cmd)

	return cmd
}

func announceListGetCategoryCmdRun(cmd *cobra.Command, args []string) {
	metadata, err := metadataFromCommand(cmd, WithInput())
	if err != nil {
		printCommandErrorAndExit(cmd, err)
	}

	input := inputs.NewSingleInput(inputs.NewDecodeBencode(), inputs.NewFileInput())
	output := outputs.NewSingleOutput(outputs.NewEncodePrintList(), outputs.NewStdOutput())

	batch := actions.NewBatch(output)
	batch.Append(actions.NewGetAnnounceListActionFunc(announceListPath))
	batch.Append(actions.NewGetObjectActionFunc(args))

	if err := input.Execute(metadata, batch.CreateFunction()); err != nil {
		printCommandErrorAndExit(cmd, err)
	}
}
