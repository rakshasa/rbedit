package common

import (
	"github.com/rakshasa/rbedit/actions"
	"github.com/rakshasa/rbedit/inputs"
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

	return cmd
}

func getCmdRun(cmd *cobra.Command, args []string) {
	if len(args) == 0 && !hasChangedFlags(cmd) {
		printCommandUsageAndExit(cmd)
	}

	metadata, output, err := metadataFromCommand(cmd,
		WithInput(),
		WithDefaultOutput(outputs.NewEncodePrint(), outputs.NewStdOutput()),
	)
	if err != nil {
		printCommandErrorAndExit(cmd, err)
	}

	input := inputs.NewSingleInput(inputs.NewDecodeBencode(), inputs.NewFileInput())

	if err := input.Execute(metadata, actions.NewGetObjectAction(output, args)); err != nil {
		printCommandErrorAndExit(cmd, err)
	}
}
