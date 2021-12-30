package actions

import (
	"github.com/rakshasa/rbedit/outputs"
	"github.com/rakshasa/rbedit/types"
)

type ActionFunc func(types.Output) types.InputResultFunc

type batch struct {
	actions []ActionFunc
}

func NewBatch() *batch {
	return &batch{
		actions: []ActionFunc{},
	}
}

func (b *batch) Append(actionFn ActionFunc) {
	b.actions = append(b.actions, actionFn)
}

func (b *batch) CreateFunction(output types.Output) types.InputResultFunc {
	if len(b.actions) == 0 {
		return func(metadata types.IOMetadata, rootObj interface{}) error {
			return nil
		}
	}

	for idx := len(b.actions) - 1; idx != 0; idx-- {
		action := b.actions[idx](output)

		output = outputs.NewChainOutput(func(metadata types.IOMetadata, object interface{}) error {
			return action(metadata, object)
		})
	}

	return b.actions[0](output)
}
