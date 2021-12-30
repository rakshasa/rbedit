package inputs

import (
	"bytes"
	"fmt"

	bencode "github.com/rakshasa/bencode-go"
	"github.com/rakshasa/rbedit/objects"
	"github.com/rakshasa/rbedit/types"
)

// DecodeFunc:

func NewDecodeGenericBencode() types.DecodeFunc {
	return func(metadata types.IOMetadata, data []byte) (types.IOMetadata, interface{}, error) {
		object, err := bencode.Decode(bytes.NewReader(data))
		if err != nil {
			return types.IOMetadata{}, nil, fmt.Errorf("failed to decode bencode as generic object: %v", err)
		}

		return metadata, object, nil
	}
}

func NewDecodeTorrentBencode() types.DecodeFunc {
	return func(metadata types.IOMetadata, data []byte) (types.IOMetadata, interface{}, error) {
		object, err := bencode.Decode(bytes.NewReader(data))
		if err != nil {
			return types.IOMetadata{}, nil, fmt.Errorf("failed to decode bencode as torrent object: %v", err)
		}

		torrentInfo, err := objects.NewTorrentInfo(object)
		if err != nil {
			return types.IOMetadata{}, nil, fmt.Errorf("failed to decode bencode as valid torrent object: %v", err)
		}
		metadata.InputTorrentInfo = &torrentInfo

		return metadata, object, nil
	}
}
