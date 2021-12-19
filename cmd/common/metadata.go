package common

import (
	"fmt"

	"github.com/rakshasa/rbedit/outputs"
	"github.com/rakshasa/rbedit/types"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

type metadataOpFunction func(*metadataOptions)
type metadataGetFlagFunction func(*flag.FlagSet, string) (interface{}, bool, error)
type metadataOutputTargetFlagFunction func(*flag.FlagSet, string) (outputs.OutputFunc, bool, error)

type metadataOptions struct {
	input    bool
	encodeFn outputs.EncodeFunc
	outputFn outputs.OutputFunc

	getValue map[string]metadataGetFlagFunction
}

func WithInput() metadataOpFunction {
	return func(opts *metadataOptions) {
		opts.input = true
	}
}

func WithDefaultOutput(encodeFn outputs.EncodeFunc, outputFn outputs.OutputFunc) metadataOpFunction {
	return func(opts *metadataOptions) {
		opts.encodeFn = encodeFn
		opts.outputFn = outputFn
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

func hasChangedFlags(cmd *cobra.Command) bool {
	var result bool

	cmd.Flags().Visit(func(f *flag.Flag) {
		result = true
	})

	return result
}

func metadataFromCommand(cmd *cobra.Command, options ...metadataOpFunction) (types.IOMetadata, outputs.Output, error) {
	opts := newMetadataOptions(options)

	metadata := types.IOMetadata{}

	if opts.input {
		value, err := cmd.Flags().GetString(inputFlagName)
		if err != nil || len(value) == 0 {
			return types.IOMetadata{}, nil, fmt.Errorf("missing valid input source")
		}

		metadata.InputFilename = value
	}

	if flag, err := cmd.Flags().GetBool(inplaceFlagName); err == nil && flag {
		if len(metadata.InputFilename) == 0 {
			return types.IOMetadata{}, nil, fmt.Errorf("inplace output requires a file input source")
		}

		opts.outputFn = outputs.NewInplaceFileOutput()
	}
	if filename, err := cmd.Flags().GetString(outputFlagName); err == nil && len(filename) != 0 {
		if opts.outputFn != nil {
			return types.IOMetadata{}, nil, fmt.Errorf("multiple output destinations")
		}

		opts.outputFn = outputs.NewFileOutput(filename)
	}

	if opts.outputFn == nil {
		return types.IOMetadata{}, nil, fmt.Errorf("missing valid output destination")
	}
	if opts.encodeFn == nil {
		return types.IOMetadata{}, nil, fmt.Errorf("missing valid output encoder")
	}

	output := outputs.NewSingleOutput(opts.encodeFn, opts.outputFn)

	if opts.getValue != nil {
		for name, fn := range opts.getValue {
			value, ok, err := fn(cmd.Flags(), name)
			if err != nil {
				return types.IOMetadata{}, nil, fmt.Errorf("could not parse value, %v", err)
			}
			if !ok {
				continue
			}

			if metadata.Value != nil {
				return types.IOMetadata{}, nil, fmt.Errorf("multiple values not supported")
			}

			metadata.Value = value
		}

		if metadata.Value == nil {
			return types.IOMetadata{}, nil, fmt.Errorf("no value provided")
		}
	}

	return metadata, output, nil
}
