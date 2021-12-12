package actions

import (
	"fmt"
	"strconv"

	"github.com/rakshasa/rbedit/inputs"
	"github.com/rakshasa/rbedit/objects"
	"github.com/rakshasa/rbedit/outputs"
	"github.com/rakshasa/rbedit/types"
)

func NewGetObjectAction(output outputs.Output, keys []string) inputs.InputResultFunc {
	return func(rootObj interface{}, metadata inputs.IOMetadata) error {
		obj, err := objects.LookupKeyPath(rootObj, keys)
		if err != nil {
			return err
		}
		if err := output.Execute(obj, metadata); err != nil {
			return err
		}

		return nil
	}
}

func NewGetObject(keys []string) ActionFunc {
	return func(output outputs.Output) inputs.InputResultFunc {
		return NewGetObjectAction(output, keys)
	}
}

func NewGetListIndexAction(output outputs.Output, indexString string) inputs.InputResultFunc {
	return func(object interface{}, metadata inputs.IOMetadata) error {
		idx, err := strconv.Atoi(indexString)
		if err != nil || idx < 0 {
			return fmt.Errorf("not a valid list index")
		}

		listObject, ok := objects.AsList(object)
		if !ok {
			return fmt.Errorf("not a list object")
		}
		if idx >= len(listObject) {
			return fmt.Errorf("out-of-bounds")
		}
		if err := output.Execute(listObject[idx], metadata); err != nil {
			return types.PrependKeyStringIfKeysError(err, indexString)
		}

		return nil
	}
}

func NewGetListIndex(indexString string) ActionFunc {
	return func(output outputs.Output) inputs.InputResultFunc {
		return NewGetListIndexAction(output, indexString)
	}
}

func NewGetAnnounceListAppendTrackerAction(output outputs.Output, categoryIdx int, trackers []string) inputs.InputResultFunc {
	return func(rootObj interface{}, metadata inputs.IOMetadata) error {
		obj, err := objects.LookupKeyPath(rootObj, []string{"announce-list"})
		if err != nil {
			return err
		}

		announceList, err := objects.NewAnnounceList(obj)
		if err != nil {
			return fmt.Errorf("could not verify announce-list, %v", err)
		}
		if categoryIdx >= len(announceList.Categories()) {
			return fmt.Errorf("category index out-of-bounds")
		}

		for _, t := range trackers {
			(*announceList.Categories()[categoryIdx]).AppendURI(t)
		}

		rootObj, err = objects.SetObject(rootObj, announceList.ToListObject(), []string{"announce-list"})
		if err != nil {
			return err
		}
		if err := output.Execute(rootObj, metadata); err != nil {
			return err
		}

		return nil
	}
}

func NewPutAction(output outputs.Output, keys []string) inputs.InputResultFunc {
	return func(rootObj interface{}, metadata inputs.IOMetadata) error {
		rootObj, err := objects.SetObject(rootObj, metadata.Value, keys)
		if err != nil {
			return err
		}
		if err := output.Execute(rootObj, metadata); err != nil {
			return err
		}

		return nil
	}
}

func NewPut(keys []string) ActionFunc {
	return func(output outputs.Output) inputs.InputResultFunc {
		return NewPutAction(output, keys)
	}
}

func NewRemoveAction(output outputs.Output, keys []string) inputs.InputResultFunc {
	return func(rootObject interface{}, metadata inputs.IOMetadata) error {
		rootObject, err := objects.RemoveObject(rootObject, keys)
		if err != nil {
			return err
		}
		if err := output.Execute(rootObject, metadata); err != nil {
			return err
		}

		return nil
	}
}

func NewReplaceWithBatchResultAction(output outputs.Output, keys []string, actionsFn ...ActionFunc) inputs.InputResultFunc {
	batch := NewBatch()

	for _, fn := range actionsFn {
		batch.Append(fn)
	}

	return func(rootObject interface{}, metadata inputs.IOMetadata) error {
		resultOutput := outputs.NewResultOutput()
		if err := batch.CreateFunction(resultOutput)(rootObject, metadata); err != nil {
			return err
		}

		rootObject, err := objects.SetObject(rootObject, resultOutput.ResultObject(), keys)
		if err != nil {
			return err
		}
		if err := output.Execute(rootObject, metadata); err != nil {
			return err
		}

		return nil
	}
}

func NewReplaceWithBatchResult(keys []string, actionsFn ...ActionFunc) ActionFunc {
	return func(output outputs.Output) inputs.InputResultFunc {
		return NewReplaceWithBatchResultAction(output, keys, actionsFn...)
	}
}
