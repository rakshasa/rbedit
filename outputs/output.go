package outputs

import (
	"fmt"

	"github.com/rakshasa/rbedit/objects"
	"github.com/rakshasa/rbedit/types"
)

// SingleOutput:

type singleOutput struct {
	encodeFn types.EncodeFunc
	outputFn types.OutputFunc
}

func NewSingleOutput(encodeFn types.EncodeFunc, outputFn types.OutputFunc) *singleOutput {
	return &singleOutput{
		encodeFn: encodeFn,
		outputFn: outputFn,
	}
}

func (o *singleOutput) Execute(metadata types.IOMetadata, object interface{}) error {
	metadata, data, err := o.encodeFn(metadata, object)
	if err != nil {
		return err
	}

	if err := o.outputFn(metadata, data); err != nil {
		return err
	}

	return nil
}

func (o *singleOutput) ResultObject() interface{} {
	return nil
}

// TorrentOutput:

type torrentOutput struct {
	encodeFn types.EncodeFunc
	outputFn types.OutputFunc
}

func NewTorrentOutput(encodeFn types.EncodeFunc, outputFn types.OutputFunc) *torrentOutput {
	return &torrentOutput{
		encodeFn: encodeFn,
		outputFn: outputFn,
	}
}

func (o *torrentOutput) Execute(metadata types.IOMetadata, object interface{}) error {
	metadata, data, err := o.encodeFn(metadata, object)
	if err != nil {
		return err
	}

	if err := o.outputFn(metadata, data); err != nil {
		return err
	}

	return nil
}

func (o *torrentOutput) ResultObject() interface{} {
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

func (o *chainOutput) Execute(metadata types.IOMetadata, object interface{}) error {
	return o.chainFn(metadata, object)
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

func (o *emptyOutput) Execute(metadata types.IOMetadata, object interface{}) error {
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

func (o *resultOutput) Execute(metadata types.IOMetadata, object interface{}) error {
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
