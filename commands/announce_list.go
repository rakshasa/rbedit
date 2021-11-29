package commands

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
	"github.com/rakshasa/rbedit/objects"
)

var (
	announceListPath = []string{"announce-list"}
)

// AnnounceListCmd:

type AnnounceListCmd struct{}

func (*AnnounceListCmd) Name() string     { return "announce-list" }
func (*AnnounceListCmd) Synopsis() string { return "BitTorrent announce-list related commands" }
func (*AnnounceListCmd) Usage() string {
	return `
Usage:  rbedit announce-list COMMAND

BitTorrent announce-list related commands

`
}

func (c *AnnounceListCmd) SetFlags(f *flag.FlagSet) {
}

func (c *AnnounceListCmd) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	cmdr := subcommands.NewCommander(f, "announce-list")
	cmdr.Register(&AnnounceListAppendTrackerCmd{}, "")
	cmdr.Register(&AnnounceListGetCmd{}, "")
	cmdr.Register(&AnnounceListGetCategoryCmd{}, "")

	switch f.NArg() {
	case 0:
		cmdr.Explain(cmdr.Output)
		return subcommands.ExitSuccess

	default:
		return cmdr.Execute(ctx, args...)
	}
}

// AnnounceListAppendTrackerCmd:

type AnnounceListAppendTrackerCmd struct {
	CommandBase
}

func (*AnnounceListAppendTrackerCmd) Name() string { return "append-tracker" }
func (*AnnounceListAppendTrackerCmd) Synopsis() string {
	return "Append tracker to announce list category"
}
func (*AnnounceListAppendTrackerCmd) Usage() string {
	return `
Usage:  rbedit announce-list append-tracker CATEGORY

Append tracker to announce list category

`
}

func (c *AnnounceListAppendTrackerCmd) SetFlags(f *flag.FlagSet) {
	c.commonInputFlags(f)
	c.commonOutputFlags(f)

	f.Func("string", "String value", c.putString)
}

func (c *AnnounceListAppendTrackerCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	rootObj, announceList, category, statusErr := c.announceListWithCategory(f)
	if statusErr != nil {
		return printStatusError("announce-list append-tracker", statusErr)
	}

	if c.value == nil {
		return printStatusError("announce-list append-tracker", &exitUsageError{msg: "URI not provided"})
	}

	trackerURI, ok := objects.AsAbsoluteURI(c.value)
	if !ok {
		return printStatusError("announce-list append-tracker", &exitUsageError{msg: "value is not a valid absolute path URI"})
	}

	category.AppendURI(trackerURI)

	if _, statusErr := c.saveRootWithKeyPath(rootObj, announceList.ToListObject(), announceListPath); statusErr != nil {
		return printStatusError("announce-list append-tracker", statusErr)
	}

	return subcommands.ExitSuccess
}

// AnnounceListGetCmd:

type AnnounceListGetCmd struct {
	CommandBase
}

func (*AnnounceListGetCmd) Name() string     { return "get" }
func (*AnnounceListGetCmd) Synopsis() string { return "Get announce-list" }
func (*AnnounceListGetCmd) Usage() string {
	return `
Usage:  rbedit announce-list get

Get announce-list

`
}

func (c *AnnounceListGetCmd) SetFlags(f *flag.FlagSet) {
	c.commonInputFlags(f)
}

func (c *AnnounceListGetCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if f.NArg() != 0 {
		return printStatusError("announce-list get", &exitUsageError{msg: "command does not support arguments"})
	}

	_, obj, statusErr := c.loadRootWithKeyPath(announceListPath)
	if statusErr != nil {
		return printStatusError("announce-list get", statusErr)
	}

	_, err := objects.NewAnnounceList(obj)
	if err != nil {
		return printStatusError("announce-list get", &exitUsageError{msg: fmt.Sprintf("could not verify announce-list, %v", err)})
	}

	objects.PrintObject(obj)

	return subcommands.ExitSuccess
}

// AnnounceListGetCategoryCmd:

type AnnounceListGetCategoryCmd struct {
	CommandBase
}

func (*AnnounceListGetCategoryCmd) Name() string     { return "get-category" }
func (*AnnounceListGetCategoryCmd) Synopsis() string { return "Get an announce-list category" }
func (*AnnounceListGetCategoryCmd) Usage() string {
	return `
Usage:  rbedit announce-list get-category INDEX

Get announce-list category

`
}

func (c *AnnounceListGetCategoryCmd) SetFlags(f *flag.FlagSet) {
	c.commonInputFlags(f)
}

func (c *AnnounceListGetCategoryCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	_, _, category, statusErr := c.announceListWithCategory(f)
	if statusErr != nil {
		return printStatusError("announce-list get-category", statusErr)
	}

	objects.PrintList(category.ToListObject(), objects.WithValuesOnly(), objects.WithoutIndent())

	return subcommands.ExitSuccess
}
