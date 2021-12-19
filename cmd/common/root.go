package common

import (
	"github.com/spf13/cobra"
)

// TODO: Add custom PositionalArgs.

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rbedit [OPTIONS] COMMAND",
		Short: "A dependency-free bencode editor",
		// TODO: Use git tag or commit.
		Version: "0.0",
		Args:    cobra.ExactArgs(0),
		Run:     func(cmd *cobra.Command, args []string) { printCommandUsageAndExit(cmd) },
	}

	setupDefaultCommand(cmd)

	cmd.AddCommand(newAnnounceCommand())
	cmd.AddCommand(newAnnounceListCommand())
	cmd.AddCommand(newGetCommand())
	cmd.AddCommand(newHashCommand())
	cmd.AddCommand(newPutCommand())
	cmd.AddCommand(newRemoveCommand())

	return cmd
}
