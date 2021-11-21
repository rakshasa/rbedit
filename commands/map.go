package commands

import (
	"context"
	"flag"

	"github.com/google/subcommands"
	"github.com/rakshasa/rbedit/objects"
)

// MapCmd:

type MapCmd struct{}

func (*MapCmd) Name() string     { return "map" }
func (*MapCmd) Synopsis() string { return "Map commands" }
func (*MapCmd) Usage() string {
	return `Usage:  map COMMAND

Map commands

`
}

func (c *MapCmd) SetFlags(f *flag.FlagSet) {
}

func (c *MapCmd) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	cmdr := subcommands.NewCommander(f, "map")
	cmdr.Register(&MapKeysCmd{}, "")

	switch f.NArg() {
	case 0:
		cmdr.Explain(cmdr.Output)
		return subcommands.ExitSuccess

	default:
		return cmdr.Execute(ctx, args...)
	}
}

// MapsKeyCmd:

type MapKeysCmd struct {
	loader objects.Loader
}

func (*MapKeysCmd) Name() string     { return "keys" }
func (*MapKeysCmd) Synopsis() string { return "Map keys" }
func (*MapKeysCmd) Usage() string {
	return `Usage:  map keys <PREFIX>

Map keys in a hash map

`
}

func (c *MapKeysCmd) loadFile(path string) error {
	loader, err := objects.NewFileLoader(path)
	if err != nil {
		return err
	}

	c.loader = loader
	return nil
}

func (c *MapKeysCmd) SetFlags(f *flag.FlagSet) {
	const (
		fileUsage = "Input file"
	)

	f.Func("file", fileUsage, c.loadFile)
	f.Func("f", fileUsage+"(shorthand)", c.loadFile)
}

func (c *MapKeysCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// if f.NArg() > 1 {
	// 	return printStatusError("map keys", &exitUsageError{msg: "command requires none or a single key path argument"})
	// }
	keys := f.Args()

	rootObj, err := c.loader.WaitResult()
	if err != nil {
		return printStatusError("map keys", &exitFailureError{msg: err.Error()})
	}

	obj, err := objects.LookupKeyPath(rootObj, keys)
	if err != nil {
		return printStatusErrorWithKey("map keys", &exitFailureError{msg: err.Error()}, keys)
	}

	if err := objects.PrintMapObjectKeysAsPlain(obj); err != nil {
		return printStatusErrorWithKey("map keys", &exitFailureError{msg: err.Error()}, keys)
	}

	return subcommands.ExitSuccess
}
