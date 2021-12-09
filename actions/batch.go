package actions

import (
	"github.com/rakshasa/rbedit/inputs"
	"github.com/rakshasa/rbedit/outputs"
)

type ActionFunc func(outputs.Output) inputs.InputResultFunc

type batch struct {
	actions []ActionFunc
	output  outputs.Output
}

func NewBatch(output outputs.Output) *batch {
	return &batch{
		actions: []ActionFunc{},
		output:  output,
	}
}

func (b *batch) Append(actionFn ActionFunc) {
	b.actions = append(b.actions, actionFn)
}

func (b *batch) CreateFunction() inputs.InputResultFunc {
	if len(b.actions) == 0 {
		return func(rootObj interface{}, metadata inputs.IOMetadata) error {
			return nil
		}
	}

	outputHead := b.output

	for idx := len(b.actions) - 1; idx != 0; idx-- {
		action := b.actions[idx](outputHead)

		outputHead = outputs.NewChainOutput(func(object interface{}, metadata inputs.IOMetadata) error {
			return action(object, metadata)
		})
	}

	return b.actions[0](outputHead)
}
