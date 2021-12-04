package objects

import (
	"bytes"
	"fmt"
	"os"

	bencode "github.com/rakshasa/bencode-go"
)

type EncodeFunc func(interface{}) ([]byte, error)
type OutputFunc func([]byte, OutputMetadata) error

type Output interface {
	Execute(object interface{}, metadata OutputMetadata) error
}

type OutputMetadata struct {
	InputFilename string
	Inplace       bool
}

// SingleOutput:

type singleOutput struct {
	encodeFn func(interface{}) ([]byte, error)
	outputFn func([]byte, OutputMetadata) error
}

func NewSingleOutput(encodeFn EncodeFunc, outputFn OutputFunc) *singleOutput {
	return &singleOutput{
		encodeFn: encodeFn,
		outputFn: outputFn,
	}
}

func (o *singleOutput) Execute(object interface{}, metadata OutputMetadata) error {
	data, err := o.encodeFn(object)
	if err != nil {
		return fmt.Errorf("failed to encode object for output: %v", err)
	}

	if err := o.outputFn(data, metadata); err != nil {
		return fmt.Errorf("failed to write encoded object to output: %v", err)
	}

	return nil
}

// EncodeFunc:

func NewEncodeBencode() EncodeFunc {
	return func(object interface{}) ([]byte, error) {
		var buf bytes.Buffer

		if err := bencode.Marshal(&buf, object); err != nil {
			return nil, fmt.Errorf("failed to encode data: %v", err)
		}

		return buf.Bytes(), nil
	}
}

// OutputFunc:

func NewFileOutput() OutputFunc {
	return func(data []byte, metadata OutputMetadata) error {
		if !metadata.Inplace {
			return fmt.Errorf("output to file only supports inplace write")
		}
		if len(metadata.InputFilename) == 0 {
			return fmt.Errorf("output to file requires a valid input filename")
		}
		path := metadata.InputFilename

		if err := os.WriteFile(path, data, 0666); err != nil {
			return fmt.Errorf("failed to write to output: %v", err)
		}

		return nil
	}
}
