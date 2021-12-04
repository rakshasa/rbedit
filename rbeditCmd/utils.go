package rbeditCmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type addCommandFn func(context.Context) (*cobra.Command, context.Context)

type contextState interface {
	initialize()
}

func addCommand(ctx context.Context, rootCmd *cobra.Command, newCmdFn addCommandFn) context.Context {
	cmd, ctx := newCmdFn(ctx)
	rootCmd.AddCommand(cmd)
	return ctx
}

func checkValidStateOrExit(cmd *cobra.Command, ok bool) {
	if !ok {
		fmt.Fprintf(os.Stderr, "%s: could not get command state from context\n", strings.Join(commandPathAsList(cmd), " "))
		os.Exit(1)
	}
}

func commandPathAsList(cmd *cobra.Command) []string {
	path := []string{}
	for c := cmd; c.HasParent(); c = c.Parent() {
		path = append([]string{c.Name()}, path...)
	}

	return path
}

func printCommandUsage(cmd *cobra.Command) {
	fmt.Printf("%s\n\n", cmd.Long)
	cmd.Usage()
}

func printErrorAndExit(err error) {
	fmt.Fprintf(os.Stderr, "%v", err)
	os.Exit(1)
}

func printCommandErrorAndExit(cmd *cobra.Command, err error) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", strings.Join(commandPathAsList(cmd), " "), err.Error())
	os.Exit(1)
}
