package rbeditCmd

import (
	"github.com/rakshasa/rbedit/actions"
	"github.com/rakshasa/rbedit/inputs"
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
	addOutputFlags(cmd)

	return cmd
}

func removeCmdRun(cmd *cobra.Command, args []string) {
	if len(args) == 0 && !hasChangedFlags(cmd) {
		printCommandUsageAndExit(cmd)
	}

	metadata, err := metadataFromCommand(cmd, WithInput(), WithOutput())
	if err != nil {
		printCommandErrorAndExit(cmd, err)
	}

	input := inputs.NewSingleInput(inputs.NewDecodeBencode(), inputs.NewFileInput())
	output := outputs.NewSingleOutput(outputs.NewEncodeBencode(), outputs.NewFileOutput())

	if err := input.Execute(metadata, actions.NewRemove(output, args)); err != nil {
		printCommandErrorAndExit(cmd, err)
	}
}
