package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/google/subcommands"
)

func printUsageError(cmd, msg string) subcommands.ExitStatus {
	fmt.Printf("rbedit %s: %s", cmd, msg)
	return subcommands.ExitUsageError
}

func (c *listCmd) stringArg(f *flag.FlagSet) (string, subcommands.ExitStatus) {
	if f.NArg() == 0 {
		return "", printUsageError(c.Name(), "command requires an argument, got none")
	}
	if f.NArg() > 1 {
		return "", printUsageError(c.Name(), "command requires a single argument, got multiple")
	}

	return f.Arg(0), subcommands.ExitSuccess
}

type listCmd struct {
	file string
}

func (*listCmd) Name() string     { return "list" }
func (*listCmd) Synopsis() string { return "List keys" }
func (*listCmd) Usage() string {
	return `Usage:  list <PREFIX>

List keys in a hash map

`
}

func (c *listCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.file, "file", "", "Input file")
}

func (c *listCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	filename, status := c.stringArg(f)
	if status != subcommands.ExitSuccess {
		return status
	}

	fmt.Printf("list: %s\n", filename)

	return subcommands.ExitSuccess
}

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&listCmd{}, "")

	flag.Parse()
	ctx := context.Background()

	exitCode := int(subcommands.Execute(ctx))

	os.Exit(exitCode)
}
