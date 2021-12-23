package outputs

import (
	"fmt"

	"github.com/rakshasa/rbedit/objects"
	"github.com/rakshasa/rbedit/types"
)

type EncodeFunc func(interface{}) ([]byte, error)
type OutputFunc func([]byte, types.IOMetadata) error

type Output interface {
	Execute(object interface{}, metadata types.IOMetadata) error
	ResultObject() interface{}
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

func (o *singleOutput) Execute(object interface{}, metadata types.IOMetadata) error {
	data, err := o.encodeFn(object)
	if err != nil {
		return err
	}

	if err := o.outputFn(data, metadata); err != nil {
		return err
	}

	return nil
}

func (o *singleOutput) ResultObject() interface{} {
	return nil
}

// ChainOutput:

type chainOutput struct {
	chainFn types.InputResultFunc
}

func NewChainOutput(chainFn types.InputResultFunc) *chainOutput {
	return &chainOutput{
		chainFn: chainFn,
	}
}

func (o *chainOutput) Execute(object interface{}, metadata types.IOMetadata) error {
	return o.chainFn(object, metadata)
}

func (o *chainOutput) ResultObject() interface{} {
	return nil
}

// EmptyOutput:

type emptyOutput struct {
}

func NewEmptyOutput() *emptyOutput {
	return &emptyOutput{}
}

func (o *emptyOutput) Execute(object interface{}, metadata types.IOMetadata) error {
	return nil
}

func (o *emptyOutput) ResultObject() interface{} {
	return nil
}

// Result Output:

type resultOutput struct {
	result interface{}
}

func NewResultOutput() *resultOutput {
	return &resultOutput{}
}

func (o *resultOutput) Execute(object interface{}, metadata types.IOMetadata) error {
	result, err := objects.CopyObject(object)
	if err != nil {
		return fmt.Errorf("failed to copy output result object: %v", err)
	}

	o.result = result
	return nil
}

func (o *resultOutput) ResultObject() interface{} {
	return o.result
}
