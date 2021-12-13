package inputs

import (
	"bytes"
	"fmt"
	"os"

	bencode "github.com/rakshasa/bencode-go"
	"github.com/rakshasa/rbedit/types"
)

type DecodeFunc func([]byte) (interface{}, error)
type InputFunc func(types.IOMetadata) ([]byte, error)
type InputResultFunc func(interface{}, types.IOMetadata) error

type Input interface {
	// Executed once for every distinct root bencoded data object in
	// the input.
	Execute(metadata types.IOMetadata, fn InputResultFunc) error
}

// SingleInput:

type singleInput struct {
	decodeFn DecodeFunc
	inputFn  InputFunc
}

func NewSingleInput(decodeFn DecodeFunc, inputFn InputFunc) *singleInput {
	return &singleInput{
		decodeFn: decodeFn,
		inputFn:  inputFn,
	}
}

func (o *singleInput) Execute(metadata types.IOMetadata, resultFn InputResultFunc) error {
	data, err := o.inputFn(metadata)
	if err != nil {
		return err
	}

	object, err := o.decodeFn(data)
	if err != nil {
		return err
	}

	return resultFn(object, metadata)
}

// DecodeFunc:

func NewDecodeBencode() DecodeFunc {
	return func(data []byte) (interface{}, error) {
		object, err := bencode.Decode(bytes.NewReader(data))
		if err != nil {
			return nil, fmt.Errorf("failed to decode object from input: %v", err)
		}

		return object, nil
	}
}

// InputFunc:

func NewFileInput() InputFunc {
	return func(metadata types.IOMetadata) ([]byte, error) {
		data, err := os.ReadFile(metadata.InputFilename)
		if err != nil {
			if pathErr, ok := err.(*os.PathError); ok {
				err = pathErr.Err
			}

			return nil, fmt.Errorf("failed to read input, %v", err)
		}

		return data, nil
	}
}
