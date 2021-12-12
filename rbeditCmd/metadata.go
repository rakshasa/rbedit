package rbeditCmd

import (
	"fmt"
	"strings"

	"github.com/rakshasa/rbedit/inputs"
	"github.com/rakshasa/rbedit/types"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

type metadataOpFunction func(*metadataOptions)
type metadataGetFlagFunction func(*flag.FlagSet, string) (interface{}, bool, error)

type metadataOptions struct {
	input    bool
	output   bool
	keyStack [][]string
	getValue map[string]metadataGetFlagFunction
}

func WithInput() metadataOpFunction {
	return func(opts *metadataOptions) {
		opts.input = true
	}
}

func WithOutput() metadataOpFunction {
	return func(opts *metadataOptions) {
		opts.output = true
	}
}

func WithAnyValue() metadataOpFunction {
	return func(opts *metadataOptions) {
		if opts.getValue != nil {
			printErrorAndExit(fmt.Errorf("opts.getValue already initialized"))
		}

		opts.getValue = map[string]metadataGetFlagFunction{
			(bencodeValueFlagName): getBencodeValueFromFlag,
			(integerValueFlagName): getIntegerValueFromFlag,
			(jsonValueFlagName):    getJSONValueFromFlag,
			(stringValueFlagName):  getStringValueFromFlag,
		}
	}
}

func newMetadataOptions(opOptions []metadataOpFunction) *metadataOptions {
	opts := &metadataOptions{}
	for _, opt := range opOptions {
		opt(opts)
	}

	return opts
}

func (m *metadataOptions) PushKeys(keys []string) {
	m.keyStack = append(m.keyStack, keys)
}

func (m *metadataOptions) PopKeys() {
	if len(m.keyStack) == 0 {
		printErrorAndExit(fmt.Errorf("could not pop empty key stack"))
	}

	m.keyStack = m.keyStack[:len(m.keyStack)-1]
}

func (m *metadataOptions) KeysAsString() string {
	var escapedKeys []string

	for _, keys := range m.keyStack {
		for _, key := range keys {
			escapedKeys = append(escapedKeys, types.EscapeURIString(key))
		}
	}

	return strings.Join(escapedKeys, "/")
}

func hasChangedFlags(cmd *cobra.Command) bool {
	var result bool

	cmd.Flags().Visit(func(f *flag.Flag) {
		result = true
	})

	return result
}

func metadataFromCommand(cmd *cobra.Command, options ...metadataOpFunction) (inputs.IOMetadata, error) {
	opts := newMetadataOptions(options)

	metadata := inputs.IOMetadata{}

	if opts.input {
		value, err := cmd.Flags().GetString(inputFlagName)
		if err != nil || len(value) == 0 {
			return inputs.IOMetadata{}, fmt.Errorf("no valid input source")
		}

		metadata.InputFilename = value
	}

	if opts.output {
		value, err := cmd.Flags().GetBool(inplaceFlagName)
		if err != nil || !value {
			return inputs.IOMetadata{}, fmt.Errorf("no valid output destination")
		}

		metadata.Inplace = value
	}

	if opts.getValue != nil {
		for name, fn := range opts.getValue {
			value, ok, err := fn(cmd.Flags(), name)
			if err != nil {
				return inputs.IOMetadata{}, fmt.Errorf("could not parse value, %v", err)
			}
			if !ok {
				continue
			}

			if metadata.Value != nil {
				return inputs.IOMetadata{}, fmt.Errorf("multiple values not supported")
			}

			metadata.Value = value
		}

		if metadata.Value == nil {
			return inputs.IOMetadata{}, fmt.Errorf("no value provided")
		}
	}

	return metadata, nil
}
