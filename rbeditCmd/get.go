package rbeditCmd

import (
	"context"
	"fmt"

	"github.com/rakshasa/rbedit/objects"
	"github.com/spf13/cobra"
)

// GetCmd:

func newGetCommand(ctx context.Context) (*cobra.Command, context.Context) {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get object",
		Long: `
Get object`,
		Run: getCmdRun,
	}

	addStateKeyPrefixToCommand(cmd, "rbedit-get-state")
	addInputFlags(ctx, cmd)

	return cmd, ctx
}

func getCmdRun(cmd *cobra.Command, args []string) {
	keyPath := args

	input := contextInputFromCommand(cmd)
	if input == nil {
		printCommandErrorAndExit(cmd, fmt.Errorf("no input source"))
	}

	if err := input.execute(func(rootObj interface{}) error {
		obj, err := objects.LookupKeyPath(rootObj, keyPath)
		if err != nil {
			printCommandErrorAndExit(cmd, err)
		}

		objects.PrintObject(obj)
		return nil

	}); err != nil {
		printCommandErrorAndExit(cmd, err)
	}
}
