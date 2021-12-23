package actions

import (
	"github.com/rakshasa/rbedit/outputs"
	"github.com/rakshasa/rbedit/types"
)

func NewListValueAction(output outputs.Output, value []interface{}) types.InputResultFunc {
	return func(rootObj interface{}, metadata types.IOMetadata) error {
		if err := output.Execute(value, metadata); err != nil {
			return err
		}

		return nil
	}
}

func NewListValue(value []interface{}) ActionFunc {
	return func(output outputs.Output) types.InputResultFunc {
		return NewListValueAction(output, value)
	}
}

func NewStringValueAction(output outputs.Output, value string) types.InputResultFunc {
	return func(rootObj interface{}, metadata types.IOMetadata) error {
		if err := output.Execute(value, metadata); err != nil {
			return err
		}

		return nil
	}
}

func NewStringValue(value string) ActionFunc {
	return func(output outputs.Output) types.InputResultFunc {
		return NewStringValueAction(output, value)
	}
}
