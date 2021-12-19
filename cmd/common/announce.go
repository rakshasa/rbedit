package common

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
	metadata, output, err := metadataFromCommand(cmd,
		WithInput(),
		WithDefaultOutput(outputs.NewEncodePrint(), outputs.NewStdOutput()),
	)
	if err != nil {
		printCommandErrorAndExit(cmd, err)
	}

	input := inputs.NewSingleInput(inputs.NewDecodeBencode(), inputs.NewFileInput())

	batch := actions.NewBatch()
	batch.Append(actions.NewGetObject(announcePath))
	batch.Append(actions.NewVerifyResultIsURI())

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
	metadata, output, err := metadataFromCommand(cmd,
		WithInput(),
		WithDefaultOutput(outputs.NewEncodeBencode(), nil),
	)
	if err != nil {
		printCommandErrorAndExit(cmd, err)
	}
	metadata.Value = args[0]

	input := inputs.NewSingleInput(inputs.NewDecodeBencode(), inputs.NewFileInput())

	batch := actions.NewBatch()
	batch.Append(actions.NewVerifyValueIsURI())
	batch.Append(actions.NewPut(announcePath))

	if err := input.Execute(metadata, batch.CreateFunction(output)); err != nil {
		printCommandErrorAndExit(cmd, err)
	}
}
