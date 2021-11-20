package commands

import (
	"fmt"
	"os"

	"github.com/google/subcommands"
)

type ExitStatusError interface {
	Error() string
	Status() subcommands.ExitStatus
}

type exitUsageError struct {
	msg string
}

func (e *exitUsageError) Error() string {
	return e.msg
}

func (e *exitUsageError) Status() subcommands.ExitStatus {
	return subcommands.ExitUsageError
}

func printStatusError(cmd string, err ExitStatusError) subcommands.ExitStatus {
	fmt.Fprintf(os.Stderr, "rbedit %s: %s\n", cmd, err.Error())
	return err.Status()
}
