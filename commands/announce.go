package commands

import (
	"context"
	"flag"

	"github.com/google/subcommands"
	"github.com/rakshasa/rbedit/objects"
)

var (
	announcePath = []string{"announce"}
)

// AnnounceCmd:

type AnnounceCmd struct{}

func (*AnnounceCmd) Name() string     { return "announce" }
func (*AnnounceCmd) Synopsis() string { return "BitTorrent announce related commands" }
func (*AnnounceCmd) Usage() string {
	return `
Usage:  rbedit announce COMMAND

BitTorrent announce related commands

`
}

func (c *AnnounceCmd) SetFlags(f *flag.FlagSet) {
}

func (c *AnnounceCmd) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	cmdr := subcommands.NewCommander(f, "announce")
	cmdr.Register(&AnnounceGetCmd{}, "")
	cmdr.Register(&AnnouncePutCmd{}, "")

	switch f.NArg() {
	case 0:
		cmdr.Explain(cmdr.Output)
		return subcommands.ExitSuccess

	default:
		return cmdr.Execute(ctx, args...)
	}
}

// AnnounceGetCmd:

type AnnounceGetCmd struct {
	CommandBase
}

func (*AnnounceGetCmd) Name() string     { return "get" }
func (*AnnounceGetCmd) Synopsis() string { return "Get announce url" }
func (*AnnounceGetCmd) Usage() string {
	return `
Usage:  rbedit announce get

Get announce url

`
}

func (c *AnnounceGetCmd) SetFlags(f *flag.FlagSet) {
	c.commonInputFlags(f)
}

func (c *AnnounceGetCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if f.NArg() != 0 {
		return printStatusError("announce get", &exitUsageError{msg: "command does not support arguments"})
	}

	_, obj, statusErr := c.loadRootWithKeyPath(announcePath)
	if statusErr != nil {
		return printStatusError("announce get", statusErr)
	}

	_, ok := objects.AsAbsoluteURI(obj)
	if !ok {
		return printStatusError("announce get", &exitUsageError{msg: "announce not a valid URI string"})
	}

	objects.PrintObject(obj)

	return subcommands.ExitSuccess
}

// AnnouncesPutCmd:

type AnnouncePutCmd struct {
	CommandBase
}

func (*AnnouncePutCmd) Name() string     { return "put" }
func (*AnnouncePutCmd) Synopsis() string { return "Put announce url" }
func (*AnnouncePutCmd) Usage() string {
	return `
Usage:  rbedit announce put

Put announce url

`
}

func (c *AnnouncePutCmd) SetFlags(f *flag.FlagSet) {
	c.commonInputFlags(f)
	c.commonOutputFlags(f)

	f.Func("string", "String value", c.putString)
}

func (c *AnnouncePutCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if f.NArg() != 0 {
		return printStatusError("announce put", &exitUsageError{msg: "command does not support arguments"})
	}
	if c.value == nil {
		return printStatusError("announce put", &exitUsageError{msg: "new announce URI not provided"})
	}

	_, ok := objects.AsAbsoluteURI(c.value)
	if !ok {
		return printStatusError("announce put", &exitUsageError{msg: "announce field requires an absolute path URI"})
	}

	rootObj, statusErr := c.loadRoot()
	if statusErr != nil {
		return printStatusError("announce put", statusErr)
	}

	_, statusErr = c.saveRootWithKeyPath(rootObj, c.value, announcePath)
	if statusErr != nil {
		return printStatusError("announce put", statusErr)
	}

	return subcommands.ExitSuccess
}
