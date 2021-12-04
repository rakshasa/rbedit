package rbeditCmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

const (
	contextFlagsKey = "rbedit-flags"
)

type contextFlagMap map[string]interface{}

func initContextFlagStateMap(ctx context.Context) context.Context {
	if s := ctx.Value(contextFlagsKey); s != nil {
		printErrorAndExit(fmt.Errorf("context flag state map already initialized"))
	}

	return context.WithValue(ctx, contextFlagsKey, &contextFlagMap{})
}

func addContextFlagsToMap(ctx context.Context, flagKey string, fn func(*contextFlagMap) interface{}) {
	flagStates, ok := ctx.Value(contextFlagsKey).(*contextFlagMap)
	if !ok {
		printErrorAndExit(fmt.Errorf("context flag state map not initialized"))
	}

	value := ctx.Value(flagKey)
	if value != nil {
		printErrorAndExit(fmt.Errorf("context flag state already exists: %s", flagKey))
	}

	(*flagStates)[flagKey] = fn(flagStates)
}

// Add command flags:

// TODO: Simplify by using pflag.GetFoo().

func addInputFlags(ctx context.Context, cmd *cobra.Command) {
	addContextFlagsToMap(ctx, stateKeyPrefixFromCommand(cmd)+"-input", func(flagStates *contextFlagMap) interface{} {
		state := &inputState{}

		cmd.Flags().StringVarP(&state.filePath, "file", "f", "", "Input a single file by path")
		return state
	})
}

func addOutputFlags(ctx context.Context, cmd *cobra.Command) {
	addContextFlagsToMap(ctx, stateKeyPrefixFromCommand(cmd)+"-output", func(flagStates *contextFlagMap) interface{} {
		state := &outputState{}

		cmd.Flags().BoolVar(&state.inplace, "inplace", false, "Output inplace to input file")
		return state
	})
}

func addURIFlags(ctx context.Context, cmd *cobra.Command) {
	addContextFlagsToMap(ctx, stateKeyPrefixFromCommand(cmd)+"-uri", func(flagStates *contextFlagMap) interface{} {
		state := new(uriValue)

		cmd.Flags().Var(state, "uri", "Absolute path URI value")
		return state
	})
}

func addAnyValueFlags(ctx context.Context, cmd *cobra.Command) {
	addContextFlagsToMap(ctx, stateKeyPrefixFromCommand(cmd)+"-any-value", func(flagStates *contextFlagMap) interface{} {
		state := &valueState{}

		categoryPtr := (*stringValue)(&state.category)
		valuePtr := (*stringValue)(&state.value)

		cmd.Flags().Var(categoryPtr, "category", "Value type to write")
		cmd.Flags().Var(valuePtr, "value", "Value to write")
		return state
	})
}

// Get context states:

func contextStateFromMap(ctx context.Context, flagKey string) interface{} {
	flagStates, ok := ctx.Value(contextFlagsKey).(*contextFlagMap)
	if !ok {
		printErrorAndExit(fmt.Errorf("context flag state map not initialized"))
	}

	return (*flagStates)[flagKey]
}

func contextInputFromCommand(cmd *cobra.Command) *inputState {
	value := contextStateFromMap(cmd.Context(), stateKeyPrefixFromCommand(cmd)+"-input")
	if value == nil {
		return nil
	}

	state, ok := value.(*inputState)
	if !ok {
		printErrorAndExit(fmt.Errorf("stored context state is not an input state type"))
	}

	return state
}

func contextOutputFromCommand(cmd *cobra.Command) *outputState {
	value := contextStateFromMap(cmd.Context(), stateKeyPrefixFromCommand(cmd)+"-output")
	if value == nil {
		return nil
	}

	state, ok := value.(*outputState)
	if !ok {
		printErrorAndExit(fmt.Errorf("stored context state is not an output state type"))
	}

	return state
}

func contextURIFromCommand(cmd *cobra.Command) *uriValue {
	value := contextStateFromMap(cmd.Context(), stateKeyPrefixFromCommand(cmd)+"-uri")
	if value == nil {
		return nil
	}

	state, ok := value.(*uriValue)
	if !ok {
		printErrorAndExit(fmt.Errorf("stored context state is not a uri value type"))
	}

	return state
}

func contextAnyValueFromCommand(cmd *cobra.Command) *valueState {
	value := contextStateFromMap(cmd.Context(), stateKeyPrefixFromCommand(cmd)+"-any-value")
	if value == nil {
		return nil
	}

	state, ok := value.(*valueState)
	if !ok {
		printErrorAndExit(fmt.Errorf("stored context state is not a value state type"))
	}

	if len(state.value) == 0 {
		return nil
	}

	return state
}
