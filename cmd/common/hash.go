package common

import (
	"github.com/rakshasa/rbedit/actions"
	"github.com/rakshasa/rbedit/data/encodings"
	"github.com/rakshasa/rbedit/data/outputs"
	"github.com/spf13/cobra"
)

var (
	infoHashKeys = []string{"info"}
)

// Hash Command:

func newHashCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hash COMMAND",
		Short: "SHA256 hashing commands",
		Args:  cobra.ExactArgs(0),
		Run:   func(cmd *cobra.Command, args []string) { printCommandUsageAndExit(cmd) },
	}

	setupDefaultCommand(cmd)

	cmd.AddCommand(newHashInfoCommand())

	return cmd
}

// HashInfoCommand:

func newHashInfoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info [OPTIONS]",
		Short: "Get info hash",
		Args:  cobra.ExactArgs(0),

		Run: func(cmd *cobra.Command, args []string) {
			metadata, input, output, err := metadataFromCommand(cmd,
				WithDefaultInput(),
				WithDefaultOutput(encodings.NewEncodeAsHexString(), outputs.NewStandardOutput()),
			)
			if err != nil {
				printCommandErrorAndExit(cmd, err)
			}

			batch := actions.NewBatch()
			// TODO: Allow both Bytes and Hash alone.
			batch.Append(actions.NewTemplateExecute("{{ .Input.Torrent.Hash.Bytes }}"))

			if err := input.Execute(metadata, batch.CreateFunction(output)); err != nil {
				printCommandErrorAndExit(cmd, err)
			}
		},
	}

	setupDefaultCommand(cmd)
	addInputFlags(cmd)
	addDataOutputFlags(cmd)

	return cmd
}
