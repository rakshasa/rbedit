package common

import (
	"github.com/rakshasa/rbedit/actions"
	"github.com/rakshasa/rbedit/outputs"
	"github.com/spf13/cobra"
)

// PutCmd:

func newPutCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "put [OPTIONS] KEY-PATH...",
		Short: "Put object",
		Run:   putCmdRun,
	}

	setupDefaultCommand(cmd)
	addInputFlags(cmd)
	addFileOutputFlags(cmd)
	addAnyValueFlags(cmd)

	return cmd
}

func putCmdRun(cmd *cobra.Command, args []string) {
	if len(args) == 0 && !hasChangedFlags(cmd) {
		printCommandUsageAndExit(cmd)
	}

	metadata, input, output, err := metadataFromCommand(cmd,
		WithDefaultInput(),
		WithDefaultOutput(outputs.NewEncodeTorrentBencode(), nil),
		WithAnyValue(),
	)
	if err != nil {
		printCommandErrorAndExit(cmd, err)
	}

	if err := input.Execute(metadata, actions.NewPutAction(output, args)); err != nil {
		printCommandErrorAndExit(cmd, err)
	}
}
