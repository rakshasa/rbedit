package outputs

import (
	"github.com/rakshasa/rbedit/inputs"
)

type EncodeFunc func(interface{}) ([]byte, error)
type OutputFunc func([]byte, inputs.IOMetadata) error

type Output interface {
	Execute(object interface{}, metadata inputs.IOMetadata) error
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

func (o *singleOutput) Execute(object interface{}, metadata inputs.IOMetadata) error {
	data, err := o.encodeFn(object)
	if err != nil {
		return err
	}

	if err := o.outputFn(data, metadata); err != nil {
		return err
	}

	return nil
}

// ChainOutput:

type chainOutput struct {
	chainFn inputs.InputResultFunc
}

func NewChainOutput(chainFn inputs.InputResultFunc) *chainOutput {
	return &chainOutput{
		chainFn: chainFn,
	}
}

func (o *chainOutput) Execute(object interface{}, metadata inputs.IOMetadata) error {
	return o.chainFn(object, metadata)
}
