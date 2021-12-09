package outputs

import (
	"fmt"
	"os"

	"github.com/rakshasa/rbedit/inputs"
)

func NewFileOutput() OutputFunc {
	return func(data []byte, metadata inputs.IOMetadata) error {
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

func NewStdOutput() OutputFunc {
	return func(data []byte, metadata inputs.IOMetadata) error {
		fmt.Printf("%s\n", data)
		return nil
	}
}
