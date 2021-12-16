package rbeditCmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/rakshasa/rbedit/types"
	"github.com/spf13/cobra"
)

func checkValidStateOrExit(cmd *cobra.Command, ok bool) {
	if !ok {
		fmt.Fprintf(os.Stderr, "%s: could not get command state from context\n", strings.Join(commandPathAsList(cmd), " "))
		os.Exit(1)
	}
}

func printCommandUsageAndExit(cmd *cobra.Command) {
	fmt.Printf("%s\n\n", cmd.Long)
	cmd.Usage()
	os.Exit(0)
}

func printErrorAndExit(err error) {
	fmt.Fprintf(os.Stderr, "%v", err)
	os.Exit(1)
}

func printCommandErrorAndExit(cmd *cobra.Command, err error) {
	cmdShort := strings.Join(commandPathAsList(cmd), " ")

	if keysErr, ok := err.(types.KeysError); ok {
		cmdShort += " (" + strings.Join(types.EscapeURIStringList(keysErr.Keys()), "/") + ")"
	}

	fmt.Fprintf(os.Stderr, "%s: %s\n", cmdShort, err.Error())

	switch e := err.(type) {
	case types.FileOutputError:
		fmt.Fprintf(os.Stderr, "\n"+
			"input:  %s\n"+
			"output: %s\n",
			e.Metadata().InputFilename,
			e.Filename(),
		)
	}

	os.Exit(1)
}
