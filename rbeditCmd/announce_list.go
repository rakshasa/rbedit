package rbeditCmd

import (
	"fmt"

	"github.com/rakshasa/rbedit/actions"
	"github.com/rakshasa/rbedit/inputs"
	"github.com/rakshasa/rbedit/outputs"
	"github.com/rakshasa/rbedit/types"
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
	cmd.AddCommand(newAnnounceListClearCommand())
	cmd.AddCommand(newAnnounceListClearCategoryCommand())
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
		if !types.VerifyAbsoluteURI(t) {
			printCommandErrorAndExit(cmd, fmt.Errorf("failed to validate URI for tracker %d\n", idx))
		}
		trackers = append(trackers, t)
	}

	metadata, output, err := metadataFromCommand(cmd,
		WithInput(),
		WithDefaultOutput(outputs.NewEncodeBencode(), nil),
	)
	if err != nil {
		printCommandErrorAndExit(cmd, err)
	}

	input := inputs.NewSingleInput(inputs.NewDecodeBencode(), inputs.NewFileInput())

	if err := input.Execute(metadata, actions.NewGetAnnounceListAppendTrackerAction(output, categoryIdx, trackers)); err != nil {
		printCommandErrorAndExit(cmd, err)
	}
}

// AnnounceListClearCmd:

func newAnnounceListClearCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clear",
		Short: "Clear announce list",
		Args:  cobra.ExactArgs(0),

		Run: func(cmd *cobra.Command, args []string) {
			metadata, output, err := metadataFromCommand(cmd,
				WithInput(),
				WithDefaultOutput(outputs.NewEncodeBencode(), nil),
			)
			if err != nil {
				printCommandErrorAndExit(cmd, err)
			}
			metadata.Value = new([]interface{})

			input := inputs.NewSingleInput(inputs.NewDecodeBencode(), inputs.NewFileInput())

			if err := input.Execute(metadata, actions.NewPutAction(output, announceListPath)); err != nil {
				printCommandErrorAndExit(cmd, err)
			}
		},
	}

	setupDefaultCommand(cmd)

	addInputFlags(cmd)
	addOutputFlags(cmd)

	return cmd
}

// AnnounceListClearCategoryCmd:

func newAnnounceListClearCategoryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clear-category",
		Short: "Clear an announce list category",
		Args:  cobra.ExactArgs(1),

		Run: func(cmd *cobra.Command, args []string) {
			metadata, output, err := metadataFromCommand(cmd,
				WithInput(),
				WithDefaultOutput(outputs.NewEncodeBencode(), nil),
			)
			if err != nil {
				printCommandErrorAndExit(cmd, err)
			}

			input := inputs.NewSingleInput(inputs.NewDecodeBencode(), inputs.NewFileInput())

			batch := actions.NewBatch()
			batch.Append(actions.NewReplaceWithBatchResult(announceListPath,
				actions.NewGetObject(announceListPath),
				actions.NewReplaceIndexWithBatchResult(args[0],
					actions.NewListValue(make([]interface{}, 0, 0)),
				),
			))

			if err := input.Execute(metadata, batch.CreateFunction(output)); err != nil {
				printCommandErrorAndExit(cmd, err)
			}
		},
	}

	setupDefaultCommand(cmd)

	addInputFlags(cmd)
	addOutputFlags(cmd)

	return cmd
}

// AnnounceListGetCmd:

func newAnnounceListGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [OPTIONS]",
		Short: "Get announce-list url",
		Args:  cobra.ExactArgs(0),

		Run: func(cmd *cobra.Command, args []string) {
			metadata, output, err := metadataFromCommand(cmd,
				WithInput(),
				WithDefaultOutput(outputs.NewEncodePrintAsListOfLists(), outputs.NewStdOutput()),
			)
			if err != nil {
				printCommandErrorAndExit(cmd, err)
			}

			input := inputs.NewSingleInput(inputs.NewDecodeBencode(), inputs.NewFileInput())

			batch := actions.NewBatch()
			batch.Append(actions.NewGetObject(announceListPath))
			batch.Append(actions.NewVerifyAnnounceList())

			if err := input.Execute(metadata, batch.CreateFunction(output)); err != nil {
				printCommandErrorAndExit(cmd, err)
			}
		},
	}

	setupDefaultCommand(cmd)
	addInputFlags(cmd)

	return cmd
}

// AnnounceListGetCategoryCmd:

func newAnnounceListGetCategoryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-category [OPTIONS] INDEX",
		Short: "Get announce-list category",
		// TODO: Create a function that verifies valid positive index.
		Args: cobra.ExactArgs(1),

		Run: func(cmd *cobra.Command, args []string) {
			metadata, output, err := metadataFromCommand(cmd,
				WithInput(),
				WithDefaultOutput(outputs.NewEncodePrintList(), outputs.NewStdOutput()),
			)
			if err != nil {
				printCommandErrorAndExit(cmd, err)
			}

			input := inputs.NewSingleInput(inputs.NewDecodeBencode(), inputs.NewFileInput())

			batch := actions.NewBatch()
			batch.Append(actions.NewGetObject(announceListPath))
			batch.Append(actions.NewVerifyResultIsList())
			batch.Append(actions.NewGetListIndex(args[0]))
			batch.Append(actions.NewVerifyAnnounceListCategory())

			if err := input.Execute(metadata, batch.CreateFunction(output)); err != nil {
				printCommandErrorAndExit(cmd, err)
			}
		},
	}

	setupDefaultCommand(cmd)
	addInputFlags(cmd)

	return cmd
}
