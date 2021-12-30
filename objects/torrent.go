package objects

import (
	"crypto/md5"
	"fmt"

	bencode "github.com/rakshasa/bencode-go"
	"github.com/rakshasa/rbedit/types"
)

func NewTorrentInfo(rootObject interface{}) (types.TorrentInfo, error) {
	if _, ok := AsMap(rootObject); !ok {
		return types.TorrentInfo{}, fmt.Errorf("root is not a map")
	}

	infoMap, ok := AsMap(LookupKey(rootObject, "info"))
	if !ok {
		return types.TorrentInfo{}, fmt.Errorf("missing valid info entry")
	}

	isStrictlyCompliant := true

	if _, ok := AsString(LookupKey(rootObject, "announce")); !ok {
		isStrictlyCompliant = false
	}
	if _, ok := AsInteger(LookupKey(infoMap, "piece length")); !ok {
		isStrictlyCompliant = false
	}
	if _, ok := AsString(LookupKey(infoMap, "pieces")); !ok {
		isStrictlyCompliant = false
	}

	name, ok := AsString(LookupKey(infoMap, "name"))
	if !ok {
		return types.TorrentInfo{}, fmt.Errorf("missing valid info::name entry")
	}

	if _, ok := AsInteger(LookupKey(infoMap, "length")); ok {
		// Do nothing.
	} else if _, ok := AsList(LookupKey(infoMap, "files")); ok {
		// Do nothing.
	} else {
		isStrictlyCompliant = false
	}

	var cachedMD5Hash *types.MD5Hash

	hashFn := func() (types.MD5Hash, error) {
		if cachedMD5Hash == nil {
			hasher := md5.New()
			if err := bencode.Marshal(hasher, infoMap); err != nil {
				return types.MD5Hash{}, fmt.Errorf("failed to calculate info hash: %v", err)
			}
			md5Hash, err := types.NewMD5HashFromHasher(hasher)
			if err != nil {
				return types.MD5Hash{}, err
			}

			cachedMD5Hash = &md5Hash
		}

		return *cachedMD5Hash, nil
	}

	return types.TorrentInfo{
		Name:              name,
		HashFn:            hashFn,
		StrictlyCompliant: isStrictlyCompliant,
	}, nil
}
