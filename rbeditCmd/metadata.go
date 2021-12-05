package rbeditCmd

import (
	"fmt"

	"github.com/rakshasa/rbedit/objects"
	"github.com/spf13/cobra"
)

type metadataOpFunction func(*metadataOptions)

type metadataOptions struct {
	input        bool
	output       bool
	requireValue bool
	stringValue  bool
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

func WithURIValue() metadataOpFunction {
	return func(opts *metadataOptions) {
		opts.requireValue = true
		opts.stringValue = true
	}
}

func metadataFromCommand(cmd *cobra.Command, options ...metadataOpFunction) (objects.IOMetadata, error) {
	opts := newMetadataOptions(options)

	metadata := objects.IOMetadata{}

	if opts.input {
		value, err := cmd.Flags().GetString("input")
		if err != nil {
			return objects.IOMetadata{}, fmt.Errorf("no valid input source")
		}

		metadata.InputFilename = value
	}

	if opts.output {
		value, err := cmd.Flags().GetBool("inplace")
		if err != nil {
			return objects.IOMetadata{}, fmt.Errorf("no valid output destination")
		}

		metadata.Inplace = value
	}

	if opts.stringValue {
		if value, err := cmd.Flags().GetString("string"); err == nil {
			metadata.Value = value
		}
	}

	if opts.requireValue && metadata.Value == nil {
		return objects.IOMetadata{}, fmt.Errorf("no value provided")
	}

	return metadata, nil
}
