package common

import (
	"github.com/rakshasa/rbedit/actions"
	"github.com/rakshasa/rbedit/outputs"
	"github.com/spf13/cobra"
)

// RemoveCmd:

func newRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove [OPTIONS] KEY-PATH...",
		Short: "Remove object",
		Args:  cobra.MinimumNArgs(1),
		Run:   removeCmdRun,
	}

	setupDefaultCommand(cmd)
	addInputFlags(cmd)
	addFileOutputFlags(cmd)

	return cmd
}

func removeCmdRun(cmd *cobra.Command, args []string) {
	if len(args) == 0 && !hasChangedFlags(cmd) {
		printCommandUsageAndExit(cmd)
	}

	metadata, input, output, err := metadataFromCommand(cmd,
		WithDefaultInput(),
		WithDefaultOutput(outputs.NewEncodeTorrentBencode(), nil),
	)
	if err != nil {
		printCommandErrorAndExit(cmd, err)
	}

	if err := input.Execute(metadata, actions.NewRemoveAction(output, args)); err != nil {
		printCommandErrorAndExit(cmd, err)
	}
}
