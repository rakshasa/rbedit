package actions

import (
	"github.com/rakshasa/rbedit/types"
)

func NewListValueAction(output types.Output, value []interface{}) types.InputResultFunc {
	return func(metadata types.IOMetadata, rootObj interface{}) error {
		if err := output.Execute(metadata, value); err != nil {
			return err
		}

		return nil
	}
}

func NewListValue(value []interface{}) ActionFunc {
	return func(output types.Output) types.InputResultFunc {
		return NewListValueAction(output, value)
	}
}

func NewStringValueAction(output types.Output, value string) types.InputResultFunc {
	return func(metadata types.IOMetadata, rootObj interface{}) error {
		if err := output.Execute(metadata, value); err != nil {
			return err
		}

		return nil
	}
}

func NewStringValue(value string) ActionFunc {
	return func(output types.Output) types.InputResultFunc {
		return NewStringValueAction(output, value)
	}
}
