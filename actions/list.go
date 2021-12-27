package actions

import (
	"fmt"
	"strconv"

	"github.com/rakshasa/rbedit/objects"
	"github.com/rakshasa/rbedit/outputs"
	"github.com/rakshasa/rbedit/types"
)

func NewAppendFromListOfBatchResultsAction(output types.Output, indexString string, actionsFn ...ActionFunc) types.InputResultFunc {
	batch := NewBatch()

	for _, fn := range actionsFn {
		batch.Append(fn)
	}

	return func(rootObject interface{}, metadata types.IOMetadata) error {
		rootList, ok := objects.AsList(rootObject)
		if !ok {
			return fmt.Errorf("not a list object")
		}

		resultOutput := outputs.NewResultOutput()
		if err := batch.CreateFunction(resultOutput)(rootList, metadata); err != nil {
			return err
		}
		resultList, ok := objects.AsList(resultOutput)
		if !ok {
			return fmt.Errorf("not a list object")
		}

		rootList = append(rootList, resultList...)
		if err := output.Execute(rootList, metadata); err != nil {
			return err
		}

		return nil
	}
}

func NewAppendFromListOfBatchResults(indexString string, actionsFn ...ActionFunc) ActionFunc {
	return func(output types.Output) types.InputResultFunc {
		return NewAppendFromListOfBatchResultsAction(output, indexString, actionsFn...)
	}
}

func NewReplaceIndexWithBatchResultAction(output types.Output, indexString string, actionsFn ...ActionFunc) types.InputResultFunc {
	batch := NewBatch()

	for _, fn := range actionsFn {
		batch.Append(fn)
	}

	return func(rootObject interface{}, metadata types.IOMetadata) error {
		idx, err := strconv.Atoi(indexString)
		if err != nil || idx < 0 {
			return fmt.Errorf("not a valid list index")
		}

		rootList, ok := objects.AsList(rootObject)
		if !ok {
			return fmt.Errorf("not a list object")
		}
		if idx >= len(rootList) {
			return fmt.Errorf("out-of-bounds")
		}
		resultOutput := outputs.NewResultOutput()
		if err := batch.CreateFunction(resultOutput)(rootList, metadata); err != nil {
			return err
		}

		rootList[idx] = resultOutput.ResultObject()
		if err := output.Execute(rootList, metadata); err != nil {
			return err
		}

		return nil
	}
}

func NewReplaceIndexWithBatchResult(indexString string, actionsFn ...ActionFunc) ActionFunc {
	return func(output types.Output) types.InputResultFunc {
		return NewReplaceIndexWithBatchResultAction(output, indexString, actionsFn...)
	}
}
