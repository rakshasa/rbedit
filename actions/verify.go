package actions

import (
	"fmt"

	"github.com/rakshasa/rbedit/objects"
	"github.com/rakshasa/rbedit/outputs"
	"github.com/rakshasa/rbedit/types"
)

func NewVerifyResultIsURIAction(output types.Output) types.InputResultFunc {
	return func(result interface{}, metadata types.IOMetadata) error {
		if _, ok := objects.AsAbsoluteURI(result); !ok {
			return fmt.Errorf("not a valid absolute path URI")
		}
		if err := output.Execute(result, metadata); err != nil {
			return err
		}

		return nil
	}
}

func NewVerifyResultIsURI() ActionFunc {
	return func(output types.Output) types.InputResultFunc {
		return NewVerifyResultIsURIAction(output)
	}
}

func NewVerifyValueIsURIAction(output types.Output) types.InputResultFunc {
	return func(result interface{}, metadata types.IOMetadata) error {
		if _, ok := objects.AsAbsoluteURI(metadata.Value); !ok {
			return fmt.Errorf("not a valid absolute path URI")
		}
		if err := output.Execute(result, metadata); err != nil {
			return err
		}

		return nil
	}
}

func NewVerifyValueIsURI() ActionFunc {
	return func(output types.Output) types.InputResultFunc {
		return NewVerifyValueIsURIAction(output)
	}
}

func NewVerifyResultIsListAction(output types.Output) types.InputResultFunc {
	return func(result interface{}, metadata types.IOMetadata) error {
		if _, ok := objects.AsList(result); !ok {
			return fmt.Errorf("could not verify: not a list")
		}
		if err := output.Execute(result, metadata); err != nil {
			return err
		}

		return nil
	}
}

func NewVerifyResultIsList() ActionFunc {
	return func(output types.Output) types.InputResultFunc {
		return NewVerifyResultIsListAction(output)
	}
}

func NewVerifyResultIsListContentAction(output types.Output, verifyFn ActionFunc) types.InputResultFunc {
	return func(result interface{}, metadata types.IOMetadata) error {
		objectList, ok := objects.AsList(result)
		if !ok {
			return fmt.Errorf("could not verify list content: not a list")
		}

		for _, childObj := range objectList {
			if err := verifyFn(outputs.NewEmptyOutput())(childObj, metadata); err != nil {
				return err
			}
		}

		if err := output.Execute(result, metadata); err != nil {
			return err
		}

		return nil
	}
}

func NewVerifyResultIsListContent(verifyFn ActionFunc) ActionFunc {
	return func(output types.Output) types.InputResultFunc {
		return NewVerifyResultIsListContentAction(output, verifyFn)
	}
}
