package actions

import (
	"fmt"

	"github.com/rakshasa/rbedit/inputs"
	"github.com/rakshasa/rbedit/objects"
	"github.com/rakshasa/rbedit/outputs"
)

func NewVerifyResultIsURIAction(output outputs.Output) inputs.InputResultFunc {
	return func(result interface{}, metadata inputs.IOMetadata) error {
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
	return func(output outputs.Output) inputs.InputResultFunc {
		return NewVerifyResultIsURIAction(output)
	}
}

func NewVerifyValueIsURIAction(output outputs.Output) inputs.InputResultFunc {
	return func(result interface{}, metadata inputs.IOMetadata) error {
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
	return func(output outputs.Output) inputs.InputResultFunc {
		return NewVerifyValueIsURIAction(output)
	}
}

func NewVerifyResultIsListAction(output outputs.Output) inputs.InputResultFunc {
	return func(result interface{}, metadata inputs.IOMetadata) error {
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
	return func(output outputs.Output) inputs.InputResultFunc {
		return NewVerifyResultIsListAction(output)
	}
}

func NewVerifyResultIsListContentAction(output outputs.Output, verifyFn ActionFunc) inputs.InputResultFunc {
	return func(result interface{}, metadata inputs.IOMetadata) error {
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
	return func(output outputs.Output) inputs.InputResultFunc {
		return NewVerifyResultIsListContentAction(output, verifyFn)
	}
}