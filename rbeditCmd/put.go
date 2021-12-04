package rbeditCmd

import (
	"context"
	"fmt"

	"github.com/rakshasa/rbedit/objects"
	"github.com/spf13/cobra"
)

// PutCmd:

func newPutCommand(ctx context.Context) (*cobra.Command, context.Context) {
	cmd := &cobra.Command{
		Use:   "put [OPTIONS] KEY-PATH...",
		Short: "Put object",
		Args:  cobra.MinimumNArgs(1),
		Run:   putCmdRun,
	}

	setupDefaultCommand(cmd, "rbedit-put-state")

	addInputFlags(ctx, cmd)
	addOutputFlags(ctx, cmd)
	addAnyValueFlags(ctx, cmd)

	return cmd, ctx
}

func putCmdRun(cmd *cobra.Command, args []string) {
	keyPath := args

	input := contextInputFromCommand(cmd)
	if input == nil {
		printCommandErrorAndExit(cmd, fmt.Errorf("no input source"))
	}
	output := contextOutputFromCommand(cmd)
	if output == nil {
		printCommandErrorAndExit(cmd, fmt.Errorf("no output target"))
	}
	value := contextAnyValueFromCommand(cmd)
	if output == nil {
		printCommandErrorAndExit(cmd, fmt.Errorf("no value"))
	}

	if err := input.execute(func(rootObj interface{}) error {
		// TODO: Fix value to return any object.
		rootObj, err := objects.SetObject(rootObj, value.value, keyPath)
		if err != nil {
			printCommandErrorAndExit(cmd, err)
		}

		if err := output.execute(rootObj, input.filePath); err != nil {
			printCommandErrorAndExit(cmd, err)
		}

		return nil

	}); err != nil {
		printCommandErrorAndExit(cmd, err)
	}
}
