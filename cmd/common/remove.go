package common

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

	metadata, output, err := metadataFromCommand(cmd,
		WithInput(),
		WithDefaultOutput(outputs.NewEncodeBencode(), nil),
	)
	if err != nil {
		printCommandErrorAndExit(cmd, err)
	}

	input := inputs.NewSingleInput(inputs.NewDecodeBencode(), inputs.NewFileInput())

	if err := input.Execute(metadata, actions.NewRemoveAction(output, args)); err != nil {
		printCommandErrorAndExit(cmd, err)
	}
}
