package common

import (
	"fmt"

	"github.com/rakshasa/rbedit/inputs"
	"github.com/rakshasa/rbedit/outputs"
	"github.com/rakshasa/rbedit/types"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

type metadataOpFunction func(*metadataOptions)
type metadataGetFlagFunction func(*flag.FlagSet, string) (interface{}, bool, error)
type metadataOutputTargetFlagFunction func(*flag.FlagSet, string) (types.OutputFunc, bool, error)

type metadataOptions struct {
	defaultDecodeFn types.DecodeFunc
	defaultEncodeFn types.EncodeFunc
	defaultOutputFn types.OutputFunc

	getValue map[string]metadataGetFlagFunction
}

func WithDefaultInput(defaultDecodeFn types.DecodeFunc) metadataOpFunction {
	return func(opts *metadataOptions) {
		opts.defaultDecodeFn = defaultDecodeFn
	}
}

func WithDefaultOutput(defaultEncodeFn types.EncodeFunc, defaultOutputFn types.OutputFunc) metadataOpFunction {
	return func(opts *metadataOptions) {
		opts.defaultEncodeFn = defaultEncodeFn
		opts.defaultOutputFn = defaultOutputFn
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

func metadataFromCommand(cmd *cobra.Command, options ...metadataOpFunction) (types.IOMetadata, types.Input, types.Output, error) {
	opts := newMetadataOptions(options)

	metadata := types.IOMetadata{}

	input, err := inputSourceAndTypeFromFlagSet(cmd.Flags(), opts)
	if err != nil {
		return types.IOMetadata{}, nil, nil, err
	}
	output, err := outputDestinatinoAndTypeFromFlagSet(cmd.Flags(), opts)
	if err != nil {
		return types.IOMetadata{}, nil, nil, err
	}

	if opts.getValue != nil {
		for name, fn := range opts.getValue {
			value, ok, err := fn(cmd.Flags(), name)
			if err != nil {
				return types.IOMetadata{}, nil, nil, fmt.Errorf("could not parse value, %v", err)
			}
			if !ok {
				continue
			}

			if metadata.Value != nil {
				return types.IOMetadata{}, nil, nil, fmt.Errorf("multiple values not supported")
			}

			metadata.Value = value
		}

		if metadata.Value == nil {
			return types.IOMetadata{}, nil, nil, fmt.Errorf("no value provided")
		}
	}

	return metadata, input, output, nil
}

func inputSourceAndTypeFromFlagSet(flagSet *flag.FlagSet, opts *metadataOptions) (types.Input, error) {
	var inputSource string
	var inputType string

	if flagSet.Changed(inputFlagName) {
		v, err := flagSet.GetString(inputFlagName)
		if err != nil {
			return nil, err
		}

		inputSource = v
	}
	if flagSet.Changed(inputBatchFlagName) {
		flag, err := flagSet.GetBool(inputBatchFlagName)
		if err != nil {
			return nil, err
		}

		if flag {
			inputType = types.BatchInputTypeName
		}
	}

	if opts.defaultDecodeFn == nil {
		return nil, fmt.Errorf("missing valid input decoder")
	}
	if len(inputSource) == 0 {
		return nil, fmt.Errorf("missing valid input source")
	}

	switch inputType {
	case types.BatchInputTypeName:
		// return inputs.NewParallelBatchInput(opts.defaultDecodeFn, inputs.NewBatchFilenameInput(inputSource)), nil
		return inputs.NewSequentialBatchInput(opts.defaultDecodeFn, inputs.NewBatchFileInput(inputSource)), nil

	case types.FileInputTypeName, "":
		return inputs.NewSingleInput(opts.defaultDecodeFn, inputs.NewFileInput(inputSource)), nil

	default:
		return nil, fmt.Errorf("unknown input source type: %s", inputType)
	}
}

func outputDestinatinoAndTypeFromFlagSet(flagSet *flag.FlagSet, opts *metadataOptions) (types.Output, error) {
	var outputValue string
	var outputType string

	if flagSet.Changed(outputFlagName) {
		v, err := flagSet.GetString(outputFlagName)
		if err != nil {
			return nil, err
		}

		outputValue = v
	}
	if flagSet.Changed(outputInplaceFlagName) {
		flag, err := flagSet.GetBool(outputInplaceFlagName)
		if err != nil {
			return nil, err
		}

		if flag {
			outputType = types.InplaceOutputTypeName
		}
	}

	if opts.defaultEncodeFn == nil {
		return nil, fmt.Errorf("missing valid output encoder")
	}

	switch outputType {
	case types.InplaceOutputTypeName:
		return outputs.NewSingleOutput(opts.defaultEncodeFn, outputs.NewInplaceFileOutput()), nil
	}

	if opts.defaultOutputFn == nil {
		return nil, fmt.Errorf("missing valid output destination")
	}

	switch outputType {
	case "":
		return outputs.NewSingleOutput(opts.defaultEncodeFn, opts.defaultOutputFn), nil

	case types.FileInputTypeName:
		if len(outputValue) == 0 {
			return nil, fmt.Errorf("missing valid output destination")
		}

		return outputs.NewSingleOutput(opts.defaultEncodeFn, outputs.NewFileOutput(outputValue)), nil

	default:
		return nil, fmt.Errorf("unknown output destination type: %s", outputType)
	}
}
