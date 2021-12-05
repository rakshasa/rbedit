package rbeditCmd

import (
	"github.com/rakshasa/rbedit/objects"
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
	keyPath := args

	metadata, err := metadataFromCommand(cmd, WithInput())
	if err != nil {
		printCommandErrorAndExit(cmd, err)
	}

	input := objects.NewSingleInput(objects.NewDecodeBencode(), objects.NewFileInput())

	if err := input.Execute(metadata, func(rootObj interface{}, metadata objects.IOMetadata) error {
		obj, err := objects.LookupKeyPath(rootObj, keyPath)
		if err != nil {
			printCommandErrorAndExit(cmd, err)
		}

		objects.PrintObject(obj)
		return nil

	}); err != nil {
		printCommandErrorAndExit(cmd, err)
	}
}
