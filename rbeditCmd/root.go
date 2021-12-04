package rbeditCmd

import (
	"context"

	"github.com/spf13/cobra"
)

func newRootCommand(ctx context.Context) (*cobra.Command, context.Context) {
	cmd := &cobra.Command{
		Use:   "rbedit",
		Short: "A bencode editor",
		Long: `
A dependency-free bencode editor`,
		// TODO: Use git tag.
		Version: "0.0",
		Args:    cobra.ExactArgs(0),
		Run:     func(cmd *cobra.Command, args []string) { printCommandUsage(cmd) },
	}

	return cmd, ctx
}

func Execute(ctx context.Context) {
	ctx = initContextFlagStateMap(ctx)

	rootCmd, ctx := newRootCommand(ctx)
	ctx = addCommand(ctx, rootCmd, newAnnounceCommand)
	ctx = addCommand(ctx, rootCmd, newAnnounceListCommand)
	ctx = addCommand(ctx, rootCmd, newGetCommand)
	ctx = addCommand(ctx, rootCmd, newPutCommand)

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		printErrorAndExit(err)
	}
}
