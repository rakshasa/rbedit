package main

import (
	"context"
	"fmt"
	"os"

	"github.com/rakshasa/rbedit/cmd/common"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var (
	markdownCmd *cobra.Command
)

func init() {
	markdownCmd = &cobra.Command{
		Use:   "rbedit-markdown [OPTIONS] DEST-DIR",
		Short: "Markdown documentation generator for rbEdit",
		Args:  cobra.RangeArgs(0, 1),
		Run: func(cmd *cobra.Command, args []string) {
			if cmd.Flags().NArg() == 0 {
				fmt.Printf("%s\n\n", cmd.Short)
				cmd.Usage()
				os.Exit(0)
			}

			destPath := cmd.Flags().Arg(0)

			fmt.Printf("generating markdown documentation...\n")
			fmt.Printf("destination: %s\n\n", destPath)

			rootCmd := common.NewRootCommand()
			common.AddRootCommandDocs(rootCmd)

			if err := doc.GenMarkdownTree(rootCmd, destPath); err != nil {
				os.Exit(1)
			}

			fmt.Printf("completed\n")
		},
	}
}

func main() {
	ctx := context.Background()

	if err := markdownCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
