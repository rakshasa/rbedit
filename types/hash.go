package types

import (
	"encoding/hex"
	"fmt"
	"hash"
	"strings"
)

type MD5Hash struct {
	hash []byte
}

func NewMD5HashFromBytes(h []byte) (MD5Hash, error) {
	if len(h) != 16 {
		return MD5Hash{}, fmt.Errorf("invalid raw md5 hash length: %d", len(h))
	}

	return MD5Hash{hash: h}, nil
}

func NewMD5HashFromHasher(hasher hash.Hash) (MD5Hash, error) {
	return NewMD5HashFromBytes(hasher.Sum([]byte{}))
}

func NewMD5HashFromHexString(h string) (MD5Hash, error) {
	bs, err := hex.DecodeString(h)
	if err != nil {
		return MD5Hash{}, err
	}

	return MD5Hash{hash: bs}, nil
}

func (h MD5Hash) Bytes() []byte {
	return h.hash
}

func (h MD5Hash) Hex() string {
	return hex.EncodeToString(h.hash)
}

func (h MD5Hash) HEX() string {
	return strings.ToUpper(hex.EncodeToString(h.hash))
}

func (h MD5Hash) Len() int {
	return len(h.hash)
}
