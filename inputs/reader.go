package inputs

import (
	"bufio"
	"fmt"
	"os"

	"github.com/rakshasa/rbedit/types"
)

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
