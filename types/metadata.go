package types

import (
	"bytes"
	"fmt"
	"path"
	"strings"
	"text/template"
)

type IOMetadata struct {
	Input             Input
	InputFilename     string
	InputTorrentInfo  *TorrentInfo
	OutputTorrentInfo *TorrentInfo
	Value             interface{}
	InfoHash          string // TODO: Remove!
}

type FilenameInfo struct {
	filename string
}

type TorrentInfo struct {
	Name              string
	HashFn            func() (MD5Hash, error)
	StrictlyCompliant bool
}

func (m *IOMetadata) ExecuteTemplate(t string) (string, error) {
	type outputTemplate struct {
		TorrentInfo *TorrentInfo
	}
	type inputTemplate struct {
		TorrentInfo *TorrentInfo
		FilenameInfo
	}
	type ioMetadataTemplate struct {
		Input  inputTemplate
		Output outputTemplate
	}

	data := ioMetadataTemplate{
		inputTemplate{
			TorrentInfo:  m.InputTorrentInfo,
			FilenameInfo: FilenameInfo{m.InputFilename},
		},
		outputTemplate{
			TorrentInfo: m.OutputTorrentInfo,
		},
	}

	tmpl, err := template.New("top").Parse(t)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, &data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (t *TorrentInfo) Hash() (MD5Hash, error) {
	return t.HashFn()
}

func (f *FilenameInfo) Filename() (string, error) {
	if len(f.filename) == 0 {
		return "", fmt.Errorf("missing input filename")
	}

	return f.filename, nil
}

func (f *FilenameInfo) Basename() (string, error) {
	if len(f.filename) == 0 {
		return "", fmt.Errorf("missing input filename")
	}

	return path.Base(f.filename), nil
}

func (f *FilenameInfo) BasenameWithoutTorrent() (string, error) {
	if len(f.filename) == 0 {
		return "", fmt.Errorf("missing input filename")
	}

	basename := path.Base(f.filename)

	if !strings.HasSuffix(strings.ToLower(basename), ".torrent") {
		return "", fmt.Errorf("missing '.torrent' suffix in input filename")
	}

	return basename[:len(basename)-len(".torrent")], nil
}

func (f *FilenameInfo) Dirname() (string, error) {
	if len(f.filename) == 0 {
		return "", fmt.Errorf("missing input filename")
	}

	return path.Dir(f.filename), nil
}
