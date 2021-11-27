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

// TorrentCmd:

type TorrentCmd struct{}

func (*TorrentCmd) Name() string     { return "torrent" }
func (*TorrentCmd) Synopsis() string { return "Torrent-specific commands" }
func (*TorrentCmd) Usage() string {
	return `Usage:  rbedit torrent COMMAND

Torrent-specific commands

`
}

func (c *TorrentCmd) SetFlags(f *flag.FlagSet) {
}

func (c *TorrentCmd) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	cmdr := subcommands.NewCommander(f, "torrent")
	cmdr.Register(&TorrentAnnounceCmd{}, "")

	switch f.NArg() {
	case 0:
		cmdr.Explain(cmdr.Output)
		return subcommands.ExitSuccess

	default:
		return cmdr.Execute(ctx, args...)
	}
}

// TorrentsAnnounceCmd:

type TorrentAnnounceCmd struct {
	CommandBase
}

func (*TorrentAnnounceCmd) Name() string     { return "announce" }
func (*TorrentAnnounceCmd) Synopsis() string { return "Torrent announce" }
func (*TorrentAnnounceCmd) Usage() string {
	return `Usage:  rbedit torrent announce COMMAND

Manipulate torrent announces

`
}

func (c *TorrentAnnounceCmd) SetFlags(f *flag.FlagSet) {
	c.commonInputFlags(f)
}

func (c *TorrentAnnounceCmd) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	cmdr := subcommands.NewCommander(f, "announce")
	cmdr.Register(&TorrentAnnounceGetCmd{}, "")
	cmdr.Register(&TorrentAnnouncePutCmd{}, "")

	switch f.NArg() {
	case 0:
		cmdr.Explain(cmdr.Output)
		return subcommands.ExitSuccess

	default:
		return cmdr.Execute(ctx, args...)
	}
}

// TorrentsAnnounceGetCmd:

type TorrentAnnounceGetCmd struct {
	CommandBase
}

func (*TorrentAnnounceGetCmd) Name() string     { return "get" }
func (*TorrentAnnounceGetCmd) Synopsis() string { return "Get torrent announce url" }
func (*TorrentAnnounceGetCmd) Usage() string {
	return `Usage:  rbedit torrent announce get

Get torrent announce url

`
}

func (c *TorrentAnnounceGetCmd) SetFlags(f *flag.FlagSet) {
	c.commonInputFlags(f)
}

func (c *TorrentAnnounceGetCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if f.NArg() != 0 {
		return printStatusError("torrent announce get", &exitUsageError{msg: "command does not support arguments"})
	}

	_, obj, statusErr := c.loadRootWithKeyPath(announcePath)
	if statusErr != nil {
		return printStatusError("torrent announce get", statusErr)
	}

	_, ok := objects.AsAbsoluteURI(obj)
	if !ok {
		return printStatusError("torrent announce get", &exitUsageError{msg: "announce not a valid URI string"})
	}

	objects.PrintObject(obj)

	return subcommands.ExitSuccess
}

// TorrentsAnnouncePutCmd:

type TorrentAnnouncePutCmd struct {
	CommandBase
}

func (*TorrentAnnouncePutCmd) Name() string     { return "put" }
func (*TorrentAnnouncePutCmd) Synopsis() string { return "Put torrent announce url" }
func (*TorrentAnnouncePutCmd) Usage() string {
	return `Usage:  rbedit torrent announce put

Put torrent announce url

`
}

func (c *TorrentAnnouncePutCmd) SetFlags(f *flag.FlagSet) {
	c.commonInputFlags(f)
	c.commonOutputFlags(f)

	f.Func("string", "String value", c.putString)
}

func (c *TorrentAnnouncePutCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if f.NArg() != 0 {
		return printStatusError("torrent announce put", &exitUsageError{msg: "command does not support arguments"})
	}

	_, ok := objects.AsAbsoluteURI(c.value)
	if !ok {
		return printStatusError("torrent announce put", &exitUsageError{msg: "value is not a valid URI string"})
	}

	rootObj, statusErr := c.loadRoot()
	if statusErr != nil {
		return printStatusError("torrent announce put", statusErr)
	}

	_, statusErr = c.saveRootWithKeyPath(rootObj, c.value, announcePath)
	if statusErr != nil {
		return printStatusError("torrent announce put", statusErr)
	}

	return subcommands.ExitSuccess
}
