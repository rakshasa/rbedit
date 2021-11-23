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
	CommandBase
}

func (*MapKeysCmd) Name() string     { return "keys" }
func (*MapKeysCmd) Synopsis() string { return "Map keys" }
func (*MapKeysCmd) Usage() string {
	return `Usage:  map keys <PREFIX>

Map keys in a hash map

`
}

func (c *MapKeysCmd) SetFlags(f *flag.FlagSet) {
	c.commonInputFlags(f)
}

func (c *MapKeysCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	keys := f.Args()

	_, obj, statusErr := c.loadRootWithKeyPath(keys)
	if statusErr != nil {
		return printStatusErrorWithKey("map keys", statusErr, keys)
	}

	if err := objects.PrintMapObject(obj, objects.WithKeysOnly()); err != nil {
		return printStatusErrorWithKey("map keys", &exitFailureError{msg: err.Error()}, keys)
	}

	return subcommands.ExitSuccess
}
