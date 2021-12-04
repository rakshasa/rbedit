package rbeditCmd

import (
	"fmt"

	"github.com/rakshasa/rbedit/objects"
	"github.com/spf13/cobra"
)

const (
	stateKeyPrefixKey = "state-key-prefix"
)

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

// Input State:

type inputState struct {
	filePath string
}

func inputStateFromInterface(value interface{}) *inputState {
	if value == nil {
		return &inputState{}
	}

	state, ok := value.(*inputState)
	if !ok {
		printErrorAndExit(fmt.Errorf("value is not an input state type"))
	}

	return state
}

func (s *inputState) execute(fn func(interface{}) error) error {
	var err error
	var input objects.Input

	if len(s.filePath) == 0 {
		return fmt.Errorf("no bencode input")
	}

	input, err = objects.NewFileInput(s.filePath)
	if err != nil {
		return err
	}

	return input.Execute(fn)
}

// Output State:

type outputState struct {
	inplace bool
}

func outputStateFromInterface(value interface{}) *outputState {
	if value == nil {
		return &outputState{}
	}

	state, ok := value.(*outputState)
	if !ok {
		printErrorAndExit(fmt.Errorf("value is not an output state type"))
	}

	return state
}

// Replace 'filePath' with a special objects.Input interface.
func (s *outputState) execute(rootObj interface{}, filePath string) error {
	output := objects.NewSingleOutput(objects.NewEncodeBencode(), objects.NewFileOutput())

	metadata := objects.OutputMetadata{
		InputFilename: filePath,
		Inplace:       s.inplace,
	}

	return output.Execute(rootObj, metadata)
}

// Value State:

type valueState struct {
	category string
	value    string
}

func valueStateFromInterface(value interface{}) *valueState {
	if value == nil {
		return &valueState{}
	}

	state, ok := value.(*valueState)
	if !ok {
		printErrorAndExit(fmt.Errorf("value is not an value state type"))
	}

	return state
}
