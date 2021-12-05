package rbeditCmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
)

// TODO: Add custom PositionalArgs.

func newRootCommand(ctx context.Context) (*cobra.Command, context.Context) {
	cmd := &cobra.Command{
		Use:   "rbedit [OPTIONS] COMMAND",
		Short: "A dependency-free bencode editor",
		// TODO: Use git tag or commit.
		Version: "0.0",
		Args:    cobra.ExactArgs(0),
		Run:     func(cmd *cobra.Command, args []string) { printCommandUsage(cmd) },
	}

	setupDefaultCommand(cmd, "rbedit-root-state")

	ctx = addCommand(ctx, cmd, newAnnounceCommand)
	ctx = addCommand(ctx, cmd, newAnnounceListCommand)
	ctx = addCommand(ctx, cmd, newGetCommand)
	ctx = addCommand(ctx, cmd, newPutCommand)

	return cmd, ctx
}

func Execute(ctx context.Context) {
	rootCmd, ctx := newRootCommand(ctx)

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
