package commands

import (
	"context"
	"flag"

	"github.com/google/subcommands"
	"github.com/rakshasa/rbedit/objects"
)

// GetCmd:

type GetCmd struct {
	CommandBase
}

func (*GetCmd) Name() string     { return "get" }
func (*GetCmd) FullName() string { return "get" }
func (*GetCmd) Synopsis() string { return "Get commands" }
func (*GetCmd) Usage() string {
	return `Usage:  get KEY/INDEX [KEY/INDEX ...]

Get commands

`
}

func (c *GetCmd) SetFlags(f *flag.FlagSet) {
	c.commonInputFlags(f)
}

func (c *GetCmd) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	if f.NArg() == 0 && !c.includeRoot {
		f.Usage()
		return subcommands.ExitSuccess
	}

	keys := f.Args()

	obj, err := c.lookupKeyPath(keys)
	if err != nil {
		return printStatusErrorWithKey("get", &exitFailureError{msg: err.Error()}, keys)
	}

	objects.PrintObject(obj)

	return subcommands.ExitSuccess
}
