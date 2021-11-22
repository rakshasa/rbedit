package commands

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/google/subcommands"
)

func printStatusError(cmd string, err ExitStatusError) subcommands.ExitStatus {
	fmt.Fprintf(os.Stderr, "rbedit %s: %s\n", cmd, err.Error())
	return err.Status()
}

func printStatusErrorWithKey(cmd string, err ExitStatusError, keys []string) subcommands.ExitStatus {
	keyPath := newKeyPathAsPlain(keys)
	if len(keyPath) == 0 {
		keyPath = "root object"
	}

	fmt.Fprintf(os.Stderr, "rbedit %s: %s: %s\n", cmd, err.Error(), keyPath)
	return err.Status()
}

func stringArg(f *flag.FlagSet) (string, ExitStatusError) {
	if f.NArg() == 0 {
		return "", &exitUsageError{msg: "command requires an argument, got none"}
	}
	if f.NArg() > 1 {
		return "", &exitUsageError{msg: "command requires a single argument, got multiple"}
	}

	return f.Arg(0), nil
}

func newKeyPathAsPlain(path []string) string {
	return strings.Join(path, "::")
}

func splitKeyPath(path string) []string {
	if len(path) == 0 {
		return []string{}
	}

	return strings.Split(path, "::")
}
