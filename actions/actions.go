package actions

import (
	"fmt"

	"github.com/rakshasa/rbedit/inputs"
	"github.com/rakshasa/rbedit/objects"
	"github.com/rakshasa/rbedit/outputs"
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

func NewGetObjectActionFunc(keys []string) ActionFunc {
	return func(output outputs.Output) inputs.InputResultFunc {
		return NewGetObjectAction(output, keys)
	}
}

func NewGetAbsoluteURIAction(output outputs.Output, keys []string) inputs.InputResultFunc {
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

func NewGetAnnounceListAction(output outputs.Output, keys []string) inputs.InputResultFunc {
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

func NewGetAnnounceListActionFunc(keys []string) ActionFunc {
	return func(output outputs.Output) inputs.InputResultFunc {
		return NewGetAnnounceListAction(output, keys)
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

func NewPutAbsoluteURIAction(output outputs.Output, keys []string) inputs.InputResultFunc {
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

func NewRemoveAction(output outputs.Output, keys []string) inputs.InputResultFunc {
	return func(rootObj interface{}, metadata inputs.IOMetadata) error {
		rootObj, err := objects.RemoveObject(rootObj, keys)
		if err != nil {
			return err
		}
		if err := output.Execute(rootObj, metadata); err != nil {
			return err
		}

		return nil
	}
}
