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
type metadataOutputTargetFlagFunction func(*flag.FlagSet, string) (outputs.OutputFunc, bool, error)

type metadataOptions struct {
	defaultDecodeFn inputs.DecodeFunc
	defaultEncodeFn outputs.EncodeFunc
	defaultOutputFn outputs.OutputFunc

	getValue map[string]metadataGetFlagFunction
}

func WithDefaultInput(defaultDecodeFn inputs.DecodeFunc) metadataOpFunction {
	return func(opts *metadataOptions) {
		opts.defaultDecodeFn = defaultDecodeFn
	}
}

func WithDefaultOutput(defaultEncodeFn outputs.EncodeFunc, defaultOutputFn outputs.OutputFunc) metadataOpFunction {
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

func metadataFromCommand(cmd *cobra.Command, options ...metadataOpFunction) (types.IOMetadata, types.Input, outputs.Output, error) {
	opts := newMetadataOptions(options)

	metadata := types.IOMetadata{}

	inputFn, err := metadataSetInputSourceAndType(&metadata, cmd.Flags())
	if err != nil {
		return types.IOMetadata{}, nil, nil, err
	}
	outputFn, err := metadataSetOutputDestinatinoAndType(&metadata, cmd.Flags(), opts)
	if err != nil {
		return types.IOMetadata{}, nil, nil, err
	}

	if opts.defaultDecodeFn == nil {
		return types.IOMetadata{}, nil, nil, fmt.Errorf("missing valid input decoder")
	}
	if opts.defaultEncodeFn == nil {
		return types.IOMetadata{}, nil, nil, fmt.Errorf("missing valid output encoder")
	}

	input := inputs.NewSingleInput(opts.defaultDecodeFn, inputFn)
	output := outputs.NewSingleOutput(opts.defaultEncodeFn, outputFn)

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

func metadataSetInputSourceAndType(metadata *types.IOMetadata, flagSet *flag.FlagSet) (inputs.InputFunc, error) {
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

	if len(inputSource) == 0 {
		return nil, fmt.Errorf("missing valid input source")
	}

	switch inputType {
	case types.BatchInputTypeName:
		// inputFn = inputs.NewBatchInput(inputSource)
		return nil, fmt.Errorf("unimplemented input source type: %s", inputType)

	case types.FileInputTypeName, "":
		// TODO: Set this in inputFn
		metadata.InputFilename = inputSource
		return inputs.NewFileInput(), nil

	default:
		return nil, fmt.Errorf("unknown input source type: %s", inputType)
	}
}

func metadataSetOutputDestinatinoAndType(metadata *types.IOMetadata, flagSet *flag.FlagSet, opts *metadataOptions) (outputs.OutputFunc, error) {
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

	switch outputType {
	case types.InplaceOutputTypeName:
		return outputs.NewInplaceFileOutput(), nil
	}

	if opts.defaultOutputFn == nil {
		return nil, fmt.Errorf("missing valid output destination")
	}

	switch outputType {
	case "":
		return opts.defaultOutputFn, nil

	case types.FileInputTypeName:
		if len(outputValue) == 0 {
			return nil, fmt.Errorf("missing valid output destination")
		}

		// TODO: Do not allow batch ops, set flag in metadata.
		return outputs.NewFileOutput(outputValue), nil

	default:
		return nil, fmt.Errorf("unknown output destination type: %s", outputType)
	}
}
