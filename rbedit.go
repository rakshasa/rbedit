package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/google/subcommands"
)

type listCmd struct {
}

func (*listCmd) Name() string     { return "list" }
func (*listCmd) Synopsis() string { return "List keys." }
func (*listCmd) Usage() string {
	return `
Usage:  print <PREFIX>

Print args to stdout
`
}

func (p *listCmd) SetFlags(f *flag.FlagSet) {
}

func (p *listCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// for _, arg := range f.Args() {
	// 	if p.capitalize {
	// 		arg = strings.ToUpper(arg)
	// 	}
	// 	fmt.Printf("%s ", arg)
	// }

	fmt.Printf("list:")

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
