package actions

import (
	"fmt"

	"github.com/rakshasa/rbedit/inputs"
	"github.com/rakshasa/rbedit/objects"
	"github.com/rakshasa/rbedit/outputs"
)

func NewVerifyResultIsURI(output outputs.Output) inputs.InputResultFunc {
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

func NewVerifyResultIsURIFunction() ActionFunc {
	return func(output outputs.Output) inputs.InputResultFunc {
		return NewVerifyResultIsURI(output)
	}
}

func NewVerifyValueIsURI(output outputs.Output) inputs.InputResultFunc {
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

func NewVerifyValueIsURIFunction() ActionFunc {
	return func(output outputs.Output) inputs.InputResultFunc {
		return NewVerifyValueIsURI(output)
	}
}

func NewVerifyResultIsList(output outputs.Output) inputs.InputResultFunc {
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

func NewVerifyResultIsListFunction() ActionFunc {
	return func(output outputs.Output) inputs.InputResultFunc {
		return NewVerifyResultIsList(output)
	}
}

func NewVerifyResultIsListContent(output outputs.Output, verifyFn ActionFunc) inputs.InputResultFunc {
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

func NewVerifyResultIsListContentFunction(verifyFn ActionFunc) ActionFunc {
	return func(output outputs.Output) inputs.InputResultFunc {
		return NewVerifyResultIsListContent(output, verifyFn)
	}
}
