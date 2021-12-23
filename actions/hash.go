package actions

import (
	"crypto/sha1"
	"fmt"

	"github.com/rakshasa/bencode-go"
	"github.com/rakshasa/rbedit/objects"
	"github.com/rakshasa/rbedit/outputs"
	"github.com/rakshasa/rbedit/types"
)

func NewSHA1Action(output outputs.Output, keys []string, target types.ResultTarget) types.InputResultFunc {
	return func(rootObj interface{}, metadata types.IOMetadata) error {
		object, err := objects.LookupKeyPath(rootObj, keys)
		if err != nil {
			return err
		}

		hasher := sha1.New()
		if err := bencode.Marshal(hasher, object); err != nil {
			return err
		}

		result := rootObj

		switch target {
		case types.ObjectResultTarget:
			result = string(hasher.Sum([]byte{}))
		case types.MetadataResultTarget:
			metadata.InfoHash = string(hasher.Sum([]byte{}))
		default:
			return fmt.Errorf("unknown output target type")
		}

		if err := output.Execute(result, metadata); err != nil {
			return err
		}

		return nil
	}
}

func NewSHA1(keys []string, target types.ResultTarget) ActionFunc {
	return func(output outputs.Output) types.InputResultFunc {
		return NewSHA1Action(output, keys, target)
	}
}

func NewCalculateInfoHash() ActionFunc {
	return NewSHA1([]string{"info"}, types.MetadataResultTarget)
}

func NewCachedInfoHashAction(output outputs.Output) types.InputResultFunc {
	return func(rootObj interface{}, metadata types.IOMetadata) error {
		if len(metadata.InfoHash) == 0 {
			return fmt.Errorf("info hash not calculated")
		}

		if err := output.Execute(metadata.InfoHash, metadata); err != nil {
			return err
		}

		return nil
	}
}

func NewCachedInfoHash() ActionFunc {
	return func(output outputs.Output) types.InputResultFunc {
		return NewCachedInfoHashAction(output)
	}
}
