package outputs

import (
	"fmt"
	"os"

	"github.com/rakshasa/rbedit/types"
)

func NewInplaceFileOutput() OutputFunc {
	return func(data []byte, metadata types.IOMetadata) error {
		if err := os.WriteFile(metadata.InputFilename, data, 0666); err != nil {
			if pathErr, ok := err.(*os.PathError); ok {
				err = pathErr.Err
			}

			return fmt.Errorf("failed to write to inplace file output, %v", err)
		}

		return nil
	}
}

func NewStdOutput() OutputFunc {
	return func(data []byte, metadata types.IOMetadata) error {
		fmt.Printf("%s\n", data)
		return nil
	}
}
