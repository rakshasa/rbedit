package templates

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/rakshasa/rbedit/types"
)

type notPrintable struct {
	executeErr *error
}

type Root struct {
	notPrintable
	Input  Input
	Output Output
}

type Input struct {
	notPrintable
	File    *FileInfo
	Torrent *TorrentInfo
}

type Output struct {
	notPrintable
	Torrent *TorrentInfo
}

type FileInfo struct {
	notPrintable
	types.FileInfo
}

type TorrentInfo struct {
	notPrintable
	base types.TorrentInfo
}

type MD5Hash struct {
	types.MD5Hash
}

func (t *notPrintable) String() string {
	if *t.executeErr == nil {
		*t.executeErr = fmt.Errorf("called print un-printable type")
	}

	return ""
}

func (t *TorrentInfo) Hash() (*MD5Hash, error) {
	hash, err := t.base.Hash()
	if err != nil {
		return nil, err
	}

	return &MD5Hash{hash}, nil
}

func (t *MD5Hash) String() string {
	return t.Hex()
}

func ExecuteTemplate(metadata types.IOMetadata, templateText string) (string, error) {
	var executeErr error
	var inputFileInfo *FileInfo
	var inputTorrentInfo *TorrentInfo
	var outputTorrentInfo *TorrentInfo

	if fileInfo, err := types.NewFileInfo(metadata.InputFilename); err != nil {
		inputFileInfo = &FileInfo{
			notPrintable{&executeErr},
			*fileInfo,
		}
	}
	if metadata.InputTorrentInfo != nil {
		inputTorrentInfo = &TorrentInfo{
			notPrintable{&executeErr},
			*metadata.InputTorrentInfo,
		}
	}
	if metadata.OutputTorrentInfo != nil {
		outputTorrentInfo = &TorrentInfo{
			notPrintable{&executeErr},
			*metadata.OutputTorrentInfo,
		}
	}

	data := Root{
		notPrintable{&executeErr},
		Input{
			notPrintable{&executeErr},
			inputFileInfo,
			inputTorrentInfo,
		},
		Output{
			notPrintable{&executeErr},
			outputTorrentInfo,
		},
	}

	tmpl, err := template.New("top").Parse(templateText)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, &data); err != nil {
		return "", err
	}
	if executeErr != nil {
		return "", executeErr
	}

	return buf.String(), nil
}
