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

func (c *MapKeysCmd) stringArg(f *flag.FlagSet) (string, ExitStatusError) {
	if f.NArg() == 0 {
		return "", &exitUsageError{msg: "command requires an argument, got none"}
	}
	if f.NArg() > 1 {
		return "", &exitUsageError{msg: "command requires a single argument, got multiple"}
	}

	return f.Arg(0), nil
}

type MapKeysCmd struct {
	filename string
}

func (*MapKeysCmd) Name() string     { return "keys" }
func (*MapKeysCmd) Synopsis() string { return "Map keys" }
func (*MapKeysCmd) Usage() string {
	return `Usage:  map keys <PREFIX>

Map keys in a hash map

`
}

func (c *MapKeysCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.filename, "file", "", "Input file")
}

func (c *MapKeysCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if len(c.filename) == 0 {
		return printStatusError("map keys", &exitUsageError{msg: "command requires a bencoded file to read"})
	}

	obj, err := objects.DecodeBencodeFile(c.filename)
	if err != nil {
		return printStatusError("map keys", &exitUsageError{msg: err.Error()})
	}

	if err := objects.PrintMapObjectKeysAsPlain(obj); err != nil {
		return printStatusError("map keys", &exitUsageError{msg: err.Error()})
	}

	return subcommands.ExitSuccess
}
