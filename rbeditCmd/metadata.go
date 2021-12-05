package rbeditCmd

import (
	"fmt"

	"github.com/rakshasa/rbedit/objects"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

type metadataOpFunction func(*metadataOptions)
type metadataGetFlagFunction func(*flag.FlagSet, string) (interface{}, bool, error)

type metadataOptions struct {
	input    bool
	output   bool
	getValue map[string]metadataGetFlagFunction
}

func newMetadataOptions(opOptions []metadataOpFunction) *metadataOptions {
	opts := &metadataOptions{}

	for _, opt := range opOptions {
		opt(opts)
	}

	return opts
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
			(stringValueFlagName):  getStringValueFromFlag,
		}
	}
}

func metadataFromCommand(cmd *cobra.Command, options ...metadataOpFunction) (objects.IOMetadata, error) {
	opts := newMetadataOptions(options)

	metadata := objects.IOMetadata{}

	if opts.input {
		value, err := cmd.Flags().GetString(inputFlagName)
		if err != nil {
			return objects.IOMetadata{}, fmt.Errorf("no valid input source")
		}

		metadata.InputFilename = value
	}

	if opts.output {
		value, err := cmd.Flags().GetBool(inplaceFlagName)
		if err != nil {
			return objects.IOMetadata{}, fmt.Errorf("no valid output destination")
		}

		metadata.Inplace = value
	}

	if opts.getValue != nil {
		for name, fn := range opts.getValue {
			value, ok, err := fn(cmd.Flags(), name)
			if err != nil {
				return objects.IOMetadata{}, fmt.Errorf("could not parse value, %v", err)
			}
			if !ok {
				continue
			}

			if metadata.Value != nil {
				return objects.IOMetadata{}, fmt.Errorf("multiple values not supported")
			}

			metadata.Value = value
		}

		if metadata.Value == nil {
			return objects.IOMetadata{}, fmt.Errorf("no value provided")
		}
	}

	return metadata, nil
}
