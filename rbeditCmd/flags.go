package rbeditCmd

import (
	"github.com/spf13/cobra"
)

// Add command flags:

func addInputFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("input", "i", "", "Input a file by path")
}

func addOutputFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("inplace", false, "Output inplace to input file")
}

func addAnyValueFlags(cmd *cobra.Command) {
	cmd.Flags().String("type", "", "Value type to write")
	cmd.Flags().String("value", "", "Value to write")
}
