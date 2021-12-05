package rbeditCmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	stateKeyPrefixKey = "state-key-prefix"
)

type valueState struct {
	category string
	value    string
}

func stateKeyPrefixFromCommand(cmd *cobra.Command) string {
	if cmd.Annotations == nil {
		printCommandErrorAndExit(cmd, fmt.Errorf("command annotation map has not been initialized"))
	}

	stateKeyPrefix, ok := cmd.Annotations[stateKeyPrefixKey]
	if !ok {
		printCommandErrorAndExit(cmd, fmt.Errorf("command is missing state key prefix in annotations map"))
	}

	return stateKeyPrefix
}
