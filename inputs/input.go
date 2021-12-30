package inputs

import (
	"fmt"

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

	metadata, d, err := o.inputFn(metadata)
	if err != nil {
		return fmt.Errorf("expected single input source, got error on getting EOF input: %v", err)
	}
	if d != nil {
		return fmt.Errorf("expected single input source")
	}

	metadata, object, err := o.decodeFn(metadata, data)
	if err != nil {
		return err
	}

	return resultFn(metadata, object)
}

// SequentialBatchInput:

type sequentialBatchInput struct {
	decodeFn types.DecodeFunc
	inputFn  types.InputFunc
}

func NewSequentialBatchInput(decodeFn types.DecodeFunc, inputFn types.InputFunc) *sequentialBatchInput {
	return &sequentialBatchInput{
		decodeFn: decodeFn,
		inputFn:  inputFn,
	}
}

func (o *sequentialBatchInput) Execute(metadata types.IOMetadata, resultFn types.InputResultFunc) error {
	for {
		metadata, data, err := o.inputFn(metadata)
		if err != nil {
			return err
		}
		if data == nil {
			break
		}

		metadata, object, err := o.decodeFn(metadata, data)
		if err != nil {
			return err
		}

		if err := resultFn(metadata, object); err != nil {
			return err
		}
	}

	return nil
}
