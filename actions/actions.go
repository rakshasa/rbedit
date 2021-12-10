package actions

import (
	"fmt"
	"strconv"

	"github.com/rakshasa/rbedit/inputs"
	"github.com/rakshasa/rbedit/objects"
	"github.com/rakshasa/rbedit/outputs"
)

// TODO: Use short name for the more used version:

func NewGetObject(output outputs.Output, keys []string) inputs.InputResultFunc {
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

func NewGetObjectFunction(keys []string) ActionFunc {
	return func(output outputs.Output) inputs.InputResultFunc {
		return NewGetObject(output, keys)
	}
}

func NewGetListIndex(output outputs.Output, indexString string) inputs.InputResultFunc {
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
			return err
		}

		return nil
	}
}

func NewGetListIndexFunction(indexString string) ActionFunc {
	return func(output outputs.Output) inputs.InputResultFunc {
		return NewGetListIndex(output, indexString)
	}
}

func NewGetAbsoluteURI(output outputs.Output, keys []string) inputs.InputResultFunc {
	return func(rootObj interface{}, metadata inputs.IOMetadata) error {
		obj, err := objects.LookupKeyPath(rootObj, keys)
		if err != nil {
			return err
		}
		if _, ok := objects.AsAbsoluteURI(obj); !ok {
			return fmt.Errorf("not a valid absolute path URI")
		}
		if err := output.Execute(obj, metadata); err != nil {
			return err
		}

		return nil
	}
}

func NewGetAnnounceList(output outputs.Output, keys []string) inputs.InputResultFunc {
	return func(rootObj interface{}, metadata inputs.IOMetadata) error {
		obj, err := objects.LookupKeyPath(rootObj, keys)
		if err != nil {
			return err
		}
		if _, err := objects.NewAnnounceList(obj); err != nil {
			return fmt.Errorf("could not verify announce-list, %v", err)
		}
		if err := output.Execute(obj, metadata); err != nil {
			return err
		}

		return nil
	}
}

func NewGetAnnounceListFunction(keys []string) ActionFunc {
	return func(output outputs.Output) inputs.InputResultFunc {
		return NewGetAnnounceList(output, keys)
	}
}

func NewGetAnnounceListAppendTracker(output outputs.Output, categoryIdx int, trackers []string) inputs.InputResultFunc {
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

func NewPut(output outputs.Output, keys []string) inputs.InputResultFunc {
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

func NewPutFunction(keys []string) ActionFunc {
	return func(output outputs.Output) inputs.InputResultFunc {
		return NewPut(output, keys)
	}
}

func NewPutAbsoluteURI(output outputs.Output, keys []string) inputs.InputResultFunc {
	return func(rootObj interface{}, metadata inputs.IOMetadata) error {
		if _, ok := objects.AsAbsoluteURI(metadata.Value); !ok {
			return fmt.Errorf("not a valid absolute path URI")
		}

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

func NewRemove(output outputs.Output, keys []string) inputs.InputResultFunc {
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

func NewReplaceListIndexWithBatchResult(output outputs.Output, indexString string, actionsFn ...ActionFunc) inputs.InputResultFunc {
	batch := NewBatch()

	for _, fn := range actionsFn {
		batch.Append(fn)
	}

	return func(rootObject interface{}, metadata inputs.IOMetadata) error {
		idx, err := strconv.Atoi(indexString)
		if err != nil || idx < 0 {
			return fmt.Errorf("not a valid list index")
		}

		listObject, ok := objects.AsList(rootObject)
		if !ok {
			return fmt.Errorf("not a list object")
		}
		if idx >= len(listObject) {
			return fmt.Errorf("out-of-bounds")
		}

		resultOutput := outputs.NewResultOutput()
		if err := batch.CreateFunction(resultOutput)(listObject, metadata); err != nil {
			return err
		}

		listObject[idx] = resultOutput.ResultObject()
		if err := output.Execute(listObject, metadata); err != nil {
			return err
		}

		return nil
	}
}

func NewReplaceListIndexWithBatchResultFunction(indexString string, actionsFn ...ActionFunc) ActionFunc {
	return func(output outputs.Output) inputs.InputResultFunc {
		return NewReplaceListIndexWithBatchResult(output, indexString, actionsFn...)
	}
}

func NewReplaceWithBatchResult(output outputs.Output, keys []string, actionsFn ...ActionFunc) inputs.InputResultFunc {
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

func NewReplaceWithBatchResultFunction(keys []string, actionsFn ...ActionFunc) ActionFunc {
	return func(output outputs.Output) inputs.InputResultFunc {
		return NewReplaceWithBatchResult(output, keys, actionsFn...)
	}
}
