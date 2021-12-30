package common

import (
	"github.com/rakshasa/rbedit/actions"
	"github.com/rakshasa/rbedit/outputs"
	"github.com/spf13/cobra"
)

// GetCmd:

func newGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [OPTIONS] [KEY-PATH]...",
		Short: "Get object",
		Run:   getCmdRun,
	}

	setupDefaultCommand(cmd)
	addInputFlags(cmd)
	addDataOutputFlags(cmd)

	return cmd
}

func getCmdRun(cmd *cobra.Command, args []string) {
	if len(args) == 0 && !hasChangedFlags(cmd) {
		printCommandUsageAndExit(cmd)
	}

	metadata, input, output, err := metadataFromCommand(cmd,
		WithDefaultInput(),
		WithDefaultOutput(outputs.NewEncodePrint(), outputs.NewStdOutput()),
	)
	if err != nil {
		printCommandErrorAndExit(cmd, err)
	}
	if err := input.Execute(metadata, actions.NewGetObjectAction(output, args)); err != nil {
		printCommandErrorAndExit(cmd, err)
	}
}
