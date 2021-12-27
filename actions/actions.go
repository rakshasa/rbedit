package actions

import (
	"fmt"
	"strconv"

	"github.com/rakshasa/rbedit/objects"
	"github.com/rakshasa/rbedit/outputs"
	"github.com/rakshasa/rbedit/types"
)

func NewGetObjectAction(output types.Output, keys []string) types.InputResultFunc {
	return func(rootObj interface{}, metadata types.IOMetadata) error {
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
	return func(output types.Output) types.InputResultFunc {
		return NewGetObjectAction(output, keys)
	}
}

func NewGetListIndexAction(output types.Output, indexString string) types.InputResultFunc {
	return func(object interface{}, metadata types.IOMetadata) error {
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
	return func(output types.Output) types.InputResultFunc {
		return NewGetListIndexAction(output, indexString)
	}
}

func NewGetAnnounceListAppendTrackerAction(output types.Output, categoryIdx int, trackers []string) types.InputResultFunc {
	return func(rootObj interface{}, metadata types.IOMetadata) error {
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

func NewPutAction(output types.Output, keys []string) types.InputResultFunc {
	return func(rootObj interface{}, metadata types.IOMetadata) error {
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
	return func(output types.Output) types.InputResultFunc {
		return NewPutAction(output, keys)
	}
}

func NewRemoveAction(output types.Output, keys []string) types.InputResultFunc {
	return func(rootObject interface{}, metadata types.IOMetadata) error {
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

func NewReplaceWithBatchResultAction(output types.Output, keys []string, actionsFn ...ActionFunc) types.InputResultFunc {
	batch := NewBatch()

	for _, fn := range actionsFn {
		batch.Append(fn)
	}

	return func(rootObject interface{}, metadata types.IOMetadata) error {
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
	return func(output types.Output) types.InputResultFunc {
		return NewReplaceWithBatchResultAction(output, keys, actionsFn...)
	}
}
