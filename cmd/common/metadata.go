package common

import (
	"fmt"
	"strings"

	"github.com/rakshasa/rbedit/data/encodings"
	"github.com/rakshasa/rbedit/data/inputs"
	"github.com/rakshasa/rbedit/data/outputs"
	"github.com/rakshasa/rbedit/types"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

type metadataOpFunction func(*metadataOptions)
type metadataGetFlagFunction func(*flag.FlagSet, string) (interface{}, bool, error)
type metadataOutputTargetFlagFunction func(*flag.FlagSet, string) (types.OutputFunc, bool, error)

type metadataOptions struct {
	defaultEncodeFn types.EncodeFunc
	defaultOutputFn types.OutputFunc

	getValue map[string]metadataGetFlagFunction
}

type inputType string
type outputType string

const (
	batchInputTypeName inputType = "batch-input"
	fileInputTypeName  inputType = "file-input"

	fileOutputTypeName                     outputType = "file-output"
	fileOutputWithInplaceTypeName          outputType = "inplace-output"
	fileOutputWithTemplateFilenameTypeName outputType = "template-output"
	printTemplateTypeName                  outputType = "print-template"
)

func WithDefaultInput() metadataOpFunction {
	return func(opts *metadataOptions) {
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

func metadataFromCommand(cmd *cobra.Command, options ...metadataOpFunction) (types.IOMetadata, types.Input, types.Output, error) {
	opts := newMetadataOptions(options)

	metadata := types.IOMetadata{}

	if opts.defaultEncodeFn == nil {
		return types.IOMetadata{}, nil, nil, fmt.Errorf("missing valid default encoding")
	}

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
	var inputType inputType
	var decodeFn types.DecodeFunc

	inputFlags := []string{}
	inputIsTorrent := true

	if v, ok := getChangedString(flagSet, inputFlagName); ok {
		inputSource = v
		inputFlags = append(inputFlags, inputFlagName)
	}
	if getChangedTrue(flagSet, inputBatchFlagName) {
		inputType = batchInputTypeName
	}
	if getChangedTrue(flagSet, inputNotTorrentFlagName) {
		inputIsTorrent = false
	}

	if len(inputFlags) > 1 {
		return nil, fmt.Errorf("multiple input sources: --%s", strings.Join(inputFlags, ", --"))
	}
	if len(inputSource) == 0 {
		return nil, fmt.Errorf("missing valid input source")
	}

	if inputIsTorrent {
		decodeFn = encodings.NewDecodeTorrentBencode()
	} else {
		decodeFn = encodings.NewDecodeGenericBencode()
	}

	switch inputType {
	case batchInputTypeName:
		return inputs.NewSequentialBatchInput(decodeFn, inputs.NewBatchFileInput(inputSource)), nil

	case fileInputTypeName, "":
		return inputs.NewSingleInput(decodeFn, inputs.NewFileInput(inputSource)), nil

	default:
		return nil, fmt.Errorf("unknown input source type: %s", inputType)
	}
}

func outputDestinatinoAndTypeFromFlagSet(flagSet *flag.FlagSet, opts *metadataOptions) (types.Output, error) {
	var outputValue string
	var encodeValue string
	var outputType outputType

	outputFlags := []string{}

	if v, ok := getChangedString(flagSet, outputFlagName); ok {
		outputValue = v
		outputFlags = append(outputFlags, outputFlagName)
	}
	if flag := getChangedTrue(flagSet, outputInplaceFlagName); flag {
		outputType = fileOutputWithInplaceTypeName
		outputFlags = append(outputFlags, outputInplaceFlagName)
	}
	if v, ok := getChangedString(flagSet, outputTemplateFlagName); ok {
		outputValue = v
		outputType = fileOutputWithTemplateFilenameTypeName
		outputFlags = append(outputFlags, outputTemplateFlagName)
	}
	if v, ok := getChangedString(flagSet, printTemplateFlagName); ok {
		encodeValue = v
		outputType = printTemplateTypeName
		outputFlags = append(outputFlags, printTemplateFlagName)
	}

	if len(outputFlags) > 1 {
		return nil, fmt.Errorf("multiple output targets: --%s", strings.Join(outputFlags, ", --"))
	}

	switch outputType {
	case fileOutputTypeName:
		if len(outputValue) == 0 {
			return nil, fmt.Errorf("missing valid output destination")
		}

		return outputs.NewSingleOutput(opts.defaultEncodeFn, outputs.NewFileOutput(outputValue)), nil

	case fileOutputWithInplaceTypeName:
		return outputs.NewSingleOutput(opts.defaultEncodeFn, outputs.NewInplaceFileOutput()), nil

	case fileOutputWithTemplateFilenameTypeName:
		if len(outputValue) == 0 {
			return nil, fmt.Errorf("missing valid output destination")
		}

		return outputs.NewSingleOutput(opts.defaultEncodeFn, outputs.NewFileOutputWithTemplateFilename(outputValue)), nil

	case printTemplateTypeName:
		if len(encodeValue) == 0 {
			return nil, fmt.Errorf("missing valid print template")
		}

		return outputs.NewSingleOutput(encodings.NewEncodePrintTemplate(encodeValue), outputs.NewStandardOutput()), nil

	case "":
		if opts.defaultOutputFn == nil {
			return nil, fmt.Errorf("unexpected missing default output")
		}

		return outputs.NewSingleOutput(opts.defaultEncodeFn, opts.defaultOutputFn), nil

	default:
		return nil, fmt.Errorf("unknown output destination type: %s", outputType)
	}
}
