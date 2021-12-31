package actions

import (
	"github.com/rakshasa/rbedit/data/templates"
	"github.com/rakshasa/rbedit/types"
)

func NewTemplateExecuteAction(output types.Output, templateText string) types.InputResultFunc {
	return func(metadata types.IOMetadata, object interface{}) error {
		value, err := templates.ExecuteTemplate(metadata, templateText)
		if err != nil {
			return err
		}

		if err := output.Execute(metadata, value); err != nil {
			return err
		}

		return nil
	}
}

func NewTemplateExecute(templateText string) ActionFunc {
	return func(output types.Output) types.InputResultFunc {
		return NewTemplateExecuteAction(output, templateText)
	}
}
