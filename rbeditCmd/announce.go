package rbeditCmd

import (
	"github.com/rakshasa/rbedit/actions"
	"github.com/rakshasa/rbedit/inputs"
	"github.com/rakshasa/rbedit/outputs"
	"github.com/spf13/cobra"
)

var (
	announcePath = []string{"announce"}
)

// AnnounceCmd:

func newAnnounceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "announce [OPTIONS] COMMAND",
		Short: "BitTorrent announce related commands",
		Args:  cobra.ExactArgs(0),
		Run:   func(cmd *cobra.Command, args []string) { printCommandUsageAndExit(cmd) },
	}

	setupDefaultCommand(cmd)

	cmd.AddCommand(newAnnounceGetCommand())
	cmd.AddCommand(newAnnouncePutCommand())

	return cmd
}

// AnnounceGetCmd:

func newAnnounceGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [OPTIONS]",
		Short: "Get announce url",
		Args:  cobra.ExactArgs(0),
		Run:   announceGetCmdRun,
	}

	setupDefaultCommand(cmd)

	addInputFlags(cmd)

	return cmd
}

func announceGetCmdRun(cmd *cobra.Command, args []string) {
	metadata, err := metadataFromCommand(cmd, WithInput())
	if err != nil {
		printCommandErrorAndExit(cmd, err)
	}

	input := inputs.NewSingleInput(inputs.NewDecodeBencode(), inputs.NewFileInput())
	output := outputs.NewSingleOutput(outputs.NewEncodePrint(), outputs.NewStdOutput())

	batch := actions.NewBatch()
	batch.Append(actions.NewGetObjectFunction(announcePath))
	batch.Append(actions.NewVerifyResultIsURIFunction())

	if err := input.Execute(metadata, batch.CreateFunction(output)); err != nil {
		printCommandErrorAndExit(cmd, err)
	}
}

// AnnouncesPutCmd:

func newAnnouncePutCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "put [OPTIONS] URI",
		Short: "Set tracker announce URI",
		Args:  cobra.ExactArgs(1),
		Run:   announcePutCmdRun,
	}

	setupDefaultCommand(cmd)

	addInputFlags(cmd)
	addOutputFlags(cmd)

	return cmd
}

func announcePutCmdRun(cmd *cobra.Command, args []string) {
	metadata, err := metadataFromCommand(cmd, WithInput(), WithOutput())
	if err != nil {
		printCommandErrorAndExit(cmd, err)
	}
	metadata.Value = args[0]

	input := inputs.NewSingleInput(inputs.NewDecodeBencode(), inputs.NewFileInput())
	output := outputs.NewSingleOutput(outputs.NewEncodeBencode(), outputs.NewFileOutput())

	batch := actions.NewBatch()
	batch.Append(actions.NewVerifyValueIsURIFunction())
	batch.Append(actions.NewPutFunction(announcePath))

	if err := input.Execute(metadata, batch.CreateFunction(output)); err != nil {
		printCommandErrorAndExit(cmd, err)
	}
}
