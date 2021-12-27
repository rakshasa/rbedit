package inputs

import (
	"bufio"
	"bytes"
	"fmt"
	"os"

	bencode "github.com/rakshasa/bencode-go"
	"github.com/rakshasa/rbedit/types"
)

// SingleInput:

type singleInput struct {
	decodeFn types.DecodeFunc
	inputFn  types.InputFunc
}

func NewSingleInput(decodeFn types.DecodeFunc, inputFn types.InputFunc) *singleInput {
	return &singleInput{
		decodeFn: decodeFn,
		inputFn:  inputFn,
	}
}

func (o *singleInput) Execute(metadata types.IOMetadata, resultFn types.InputResultFunc) error {
	metadata, data, err := o.inputFn(metadata)
	if err != nil {
		return err
	}

	_, d, err := o.inputFn(metadata)
	if err != nil {
		return fmt.Errorf("expected single input source, got error on getting EOF input: %v", err)
	}
	if d != nil {
		return fmt.Errorf("expected single input source")
	}

	object, err := o.decodeFn(data)
	if err != nil {
		return err
	}

	return resultFn(object, metadata)
}

// DecodeFunc:

func NewDecodeBencode() types.DecodeFunc {
	return func(data []byte) (interface{}, error) {
		object, err := bencode.Decode(bytes.NewReader(data))
		if err != nil {
			return nil, fmt.Errorf("failed to decode object from input: %v", err)
		}

		return object, nil
	}
}

// InputFunc:

func readFileInput(metadata types.IOMetadata, filename string) (types.IOMetadata, []byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		if pathErr, ok := err.(*os.PathError); ok {
			err = pathErr.Err
		}

		return types.IOMetadata{}, nil, fmt.Errorf("failed to read input, %v", err)
	}

	metadata.InputFilename = filename
	return metadata, data, nil
}

func NewFileInput(filename string) types.InputFunc {
	var isDone bool

	return func(metadata types.IOMetadata) (types.IOMetadata, []byte, error) {
		if isDone {
			return metadata, nil, nil
		}
		isDone = true

		return readFileInput(metadata, filename)
	}
}

func NewBatchFileInput(batchFilename string) types.InputFunc {
	file, err := os.Open(batchFilename)
	if err != nil {
		if pathErr, ok := err.(*os.PathError); ok {
			err = pathErr.Err
		}

		return func(metadata types.IOMetadata) (types.IOMetadata, []byte, error) {
			return metadata, nil, fmt.Errorf("failed to read input, %v", err)
		}
	}

	scanner := bufio.NewScanner(file)

	return func(metadata types.IOMetadata) (types.IOMetadata, []byte, error) {
		if !scanner.Scan() {
			if scanner.Err() != nil {
				return types.IOMetadata{}, nil, fmt.Errorf("failed to read input, %v", err)
			}

			return metadata, nil, nil
		}

		return readFileInput(metadata, scanner.Text())
	}
}

func NewBatchFilenameInput(batchFilename string) types.InputFunc {
	file, err := os.Open(batchFilename)
	if err != nil {
		if pathErr, ok := err.(*os.PathError); ok {
			err = pathErr.Err
		}

		return func(metadata types.IOMetadata) (types.IOMetadata, []byte, error) {
			return metadata, nil, fmt.Errorf("failed to read input, %v", err)
		}
	}

	scanner := bufio.NewScanner(file)

	return func(metadata types.IOMetadata) (types.IOMetadata, []byte, error) {
		if !scanner.Scan() {
			if scanner.Err() != nil {
				return types.IOMetadata{}, nil, fmt.Errorf("failed to read input, %v", err)
			}

			return metadata, nil, nil
		}

		return metadata, []byte(scanner.Text()), nil
	}
}
