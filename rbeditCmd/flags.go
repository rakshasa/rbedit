package rbeditCmd

import (
	"context"

	"github.com/spf13/cobra"
)

// Add command flags:

func addInputFlags(ctx context.Context, cmd *cobra.Command) {
	cmd.Flags().StringP("input", "i", "", "Input a file by path")
}

func addOutputFlags(ctx context.Context, cmd *cobra.Command) {
	cmd.Flags().Bool("inplace", false, "Output inplace to input file")
}

func addAnyValueFlags(ctx context.Context, cmd *cobra.Command) {
	cmd.Flags().String("type", "", "Value type to write")
	cmd.Flags().String("value", "", "Value to write")
}
