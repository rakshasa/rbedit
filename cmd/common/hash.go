package common

import (
	"github.com/rakshasa/rbedit/actions"
	"github.com/rakshasa/rbedit/inputs"
	"github.com/rakshasa/rbedit/outputs"
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
			metadata, output, err := metadataFromCommand(cmd,
				WithInput(),
				WithDefaultOutput(outputs.NewEncodeAsHexString(), outputs.NewStdOutput()),
			)
			if err != nil {
				printCommandErrorAndExit(cmd, err)
			}

			input := inputs.NewSingleInput(inputs.NewDecodeBencode(), inputs.NewFileInput())

			batch := actions.NewBatch()
			batch.Append(actions.NewCalculateInfoHash())
			batch.Append(actions.NewCachedInfoHash())

			// batch.Append(actions.NewSHA1(infoHashKeys, types.ObjectResultTarget))

			if err := input.Execute(metadata, batch.CreateFunction(output)); err != nil {
				printCommandErrorAndExit(cmd, err)
			}
		},
	}

	setupDefaultCommand(cmd)
	addInputFlags(cmd)

	return cmd
}
