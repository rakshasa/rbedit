package commands

import (
	"github.com/google/subcommands"
)

type ExitStatusError interface {
	Error() string
	Status() subcommands.ExitStatus
}

type exitFailureError struct {
	msg string
}

func (e *exitFailureError) Error() string {
	return e.msg
}

func (e *exitFailureError) Status() subcommands.ExitStatus {
	return subcommands.ExitFailure
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
