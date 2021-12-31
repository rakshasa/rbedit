package types

import (
	"fmt"
	"path"
	"strings"
)

type IOMetadata struct {
	Input             Input
	InputFilename     string
	InputTorrentInfo  *TorrentInfo
	OutputTorrentInfo *TorrentInfo
	Value             interface{}
}

type FileInfo struct {
	filename string
}

type TorrentInfo struct {
	name              string
	hashFn            func() (MD5Hash, error)
	strictlyCompliant bool
}

// FileInfo:

func NewFileInfo(filename string) (*FileInfo, error) {
	if len(filename) == 0 {
		return nil, fmt.Errorf("missing input filename")
	}

	return &FileInfo{filename: filename}, nil
}

func (f *FileInfo) Filename() (string, error) {
	if len(f.filename) == 0 {
		return "", fmt.Errorf("missing input filename")
	}

	return f.filename, nil
}

func (f *FileInfo) Basename() (string, error) {
	if len(f.filename) == 0 {
		return "", fmt.Errorf("missing input filename")
	}

	return path.Base(f.filename), nil
}

func (f *FileInfo) BasenameWithoutTorrent() (string, error) {
	if len(f.filename) == 0 {
		return "", fmt.Errorf("missing input filename")
	}

	basename := path.Base(f.filename)

	if !strings.HasSuffix(strings.ToLower(basename), ".torrent") {
		return "", fmt.Errorf("missing '.torrent' suffix in input filename")
	}

	return basename[:len(basename)-len(".torrent")], nil
}

func (f *FileInfo) Dirname() (string, error) {
	if len(f.filename) == 0 {
		return "", fmt.Errorf("missing input filename")
	}

	return path.Dir(f.filename), nil
}

// TorrentInfo:

func NewTorrentInfo(name string, hashFn func() (MD5Hash, error), strictlyCompliant bool) (*TorrentInfo, error) {
	if len(name) == 0 {
		return nil, fmt.Errorf("missing torrent name")
	}

	return &TorrentInfo{
		name:              name,
		hashFn:            hashFn,
		strictlyCompliant: strictlyCompliant,
	}, nil
}

func (t *TorrentInfo) Hash() (MD5Hash, error) {
	return t.hashFn()
}

func (t *TorrentInfo) Name() string {
	return t.name
}

func (t *TorrentInfo) StrictlyCompliant() bool {
	return t.strictlyCompliant
}
