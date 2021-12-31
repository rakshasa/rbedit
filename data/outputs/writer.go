package outputs

import (
	"fmt"
	"os"

	"github.com/rakshasa/rbedit/data/templates"
	"github.com/rakshasa/rbedit/types"
)

type fileOutputError struct {
	err      string
	filename string
	metadata types.IOMetadata
}

func (e *fileOutputError) Error() string              { return e.err }
func (e *fileOutputError) Filename() string           { return e.filename }
func (e *fileOutputError) Metadata() types.IOMetadata { return e.metadata }

func NewInplaceFileOutput() types.OutputFunc {
	return func(metadata types.IOMetadata, data []byte) error {
		if err := os.WriteFile(metadata.InputFilename, data, 0666); err != nil {
			if pathErr, ok := err.(*os.PathError); ok {
				err = pathErr.Err
			}

			return &fileOutputError{
				err:      fmt.Sprintf("failed to write to inplace file output, %v", err),
				filename: metadata.InputFilename,
				metadata: metadata,
			}
		}

		return nil
	}
}

func NewFileOutput(filename string) types.OutputFunc {
	return func(metadata types.IOMetadata, data []byte) error {
		if err := os.WriteFile(filename, data, 0666); err != nil {
			if pathErr, ok := err.(*os.PathError); ok {
				err = pathErr.Err
			}

			return &fileOutputError{
				err:      fmt.Sprintf("failed to write to file output, %v", err),
				filename: filename,
				metadata: metadata,
			}
		}

		return nil
	}
}

func NewFileOutputWithTemplateFilename(filenameTemplate string) types.OutputFunc {
	return func(metadata types.IOMetadata, data []byte) error {
		filename, err := templates.ExecuteTemplate(metadata, filenameTemplate)
		if err != nil {
			return &fileOutputError{
				err:      fmt.Sprintf("invalid output filename %v", err),
				filename: filenameTemplate,
				metadata: metadata,
			}
		}

		if err := os.WriteFile(filename, data, 0666); err != nil {
			if pathErr, ok := err.(*os.PathError); ok {
				err = pathErr.Err
			}

			return &fileOutputError{
				err:      fmt.Sprintf("failed to write to file output, %v", err),
				filename: filename,
				metadata: metadata,
			}
		}

		return nil
	}
}

func NewStandardOutput() types.OutputFunc {
	return func(metadata types.IOMetadata, data []byte) error {
		fmt.Printf("%s\n", data)
		return nil
	}
}
