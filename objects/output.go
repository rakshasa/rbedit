package objects

import (
	"bytes"
	"fmt"
	"os"

	bencode "github.com/rakshasa/bencode-go"
)

type EncodeFunc func(interface{}) ([]byte, error)
type OutputFunc func([]byte, IOMetadata) error

type Output interface {
	Execute(object interface{}, metadata IOMetadata) error
}

// SingleOutput:

type singleOutput struct {
	encodeFn EncodeFunc
	outputFn OutputFunc
}

func NewSingleOutput(encodeFn EncodeFunc, outputFn OutputFunc) *singleOutput {
	return &singleOutput{
		encodeFn: encodeFn,
		outputFn: outputFn,
	}
}

func (o *singleOutput) Execute(object interface{}, metadata IOMetadata) error {
	data, err := o.encodeFn(object)
	if err != nil {
		return err
	}

	if err := o.outputFn(data, metadata); err != nil {
		return err
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
	return func(data []byte, metadata IOMetadata) error {
		if !metadata.Inplace {
			return fmt.Errorf("output to file only supports inplace write")
		}
		path := metadata.InputFilename

		if err := os.WriteFile(path, data, 0666); err != nil {
			if pathErr, ok := err.(*os.PathError); ok {
				err = pathErr.Err
			}

			return fmt.Errorf("failed to write to output, %v", err)
		}

		return nil
	}
}
