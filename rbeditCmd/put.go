package rbeditCmd

import (
	"github.com/rakshasa/rbedit/objects"
	"github.com/spf13/cobra"
)

// PutCmd:

func newPutCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "put [OPTIONS] KEY-PATH...",
		Short: "Put object",
		Args:  cobra.MinimumNArgs(1),
		Run:   putCmdRun,
	}

	setupDefaultCommand(cmd)
	addInputFlags(cmd)
	addOutputFlags(cmd)
	addAnyValueFlags(cmd)

	return cmd
}

func putCmdRun(cmd *cobra.Command, args []string) {
	keyPath := args

	metadata, err := metadataFromCommand(cmd, WithInput(), WithOutput(), WithAnyValue())
	if err != nil {
		printCommandErrorAndExit(cmd, err)
	}

	input := objects.NewSingleInput(objects.NewDecodeBencode(), objects.NewFileInput())
	output := objects.NewSingleOutput(objects.NewEncodeBencode(), objects.NewFileOutput())

	if err := input.Execute(metadata, func(rootObj interface{}, metadata objects.IOMetadata) error {
		rootObj, err := objects.SetObject(rootObj, metadata.Value, keyPath)
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
