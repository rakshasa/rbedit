package commands

import (
	"context"
	"flag"

	"github.com/google/subcommands"
	"github.com/rakshasa/rbedit/objects"
)

// PutCmd:

type PutCmd struct {
	CommandBase
}

func (*PutCmd) Name() string     { return "put" }
func (*PutCmd) FullName() string { return "put" }
func (*PutCmd) Synopsis() string { return "Put commands" }
func (*PutCmd) Usage() string {
	return `Usage:  put KEY/INDEX [KEY/INDEX ...]

Put commands

`
}

func (c *PutCmd) SetFlags(f *flag.FlagSet) {
	c.commonInputFlags(f)
	c.commonOutputFlags(f)

	f.Func("string", "String value", c.putString)
}

func (c *PutCmd) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	if f.NArg() == 0 && !c.includeRoot {
		f.Usage()
		return subcommands.ExitSuccess
	}
	if f.NArg() == 0 {
		return printStatusError("put", &exitUsageError{msg: "cannot use root object as target"})
	}
	if c.value == nil {
		return printStatusError("put", &exitUsageError{msg: "no value provided"})
	}

	keys := f.Args()

	rootObj, err := objects.WaitLoaderResult(c.loader)
	if err != nil {
		return printStatusErrorWithKey("put", &exitFailureError{msg: err.Error()}, keys)
	}

	if _, statusErr := c.saveRootWithKeyPath(rootObj, c.value, keys); statusErr != nil {
		return printStatusErrorWithKey("put", statusErr, keys)
	}

	return subcommands.ExitSuccess
}
