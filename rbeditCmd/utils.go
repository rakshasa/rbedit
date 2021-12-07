package rbeditCmd

import (
	"fmt"
	"os"
	"strings"

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
	fmt.Fprintf(os.Stderr, "%s: %s\n", strings.Join(commandPathAsList(cmd), " "), err.Error())
	os.Exit(1)
}
