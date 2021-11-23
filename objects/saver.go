package objects

import (
	"bytes"
	"fmt"
	"os"

	bencode "github.com/jackpal/bencode-go"
)

type Saver interface {
	WaitResult() error
}

func WaitSaverResult(s Saver) error {
	if s == nil {
		return fmt.Errorf("could not wait for saver results, saver not initalized")
	}

	return s.WaitResult()
}

type encodeResult struct {
	err error
}

type fileSaver struct {
	encodeChan chan encodeResult
}

func NewFileSaver(path string, rootObj interface{}) (*fileSaver, error) {
	if len(path) == 0 {
		return nil, fmt.Errorf("failed to create file saver, path is empty")
	}

	file, err := os.OpenFile(path, os.O_WRONLY, 0222)
	if err != nil {
		return nil, fmt.Errorf("could not open file for bencode encoding: %v", err)
	}

	saver := &fileSaver{
		encodeChan: make(chan encodeResult),
	}

	go func() {
		defer file.Close()

		var buf bytes.Buffer

		err := bencode.Marshal(&buf, rootObj)
		if err != nil {
			saver.encodeChan <- encodeResult{err: fmt.Errorf("failed to encode bencode file: %v", err)}
			return
		}

		reader := bytes.NewReader(buf.Bytes())

		count, err := file.ReadFrom(reader)
		if err != nil {
			saver.encodeChan <- encodeResult{err: fmt.Errorf("failed to write to output: %v", err)}
			return
		}
		if count != int64(reader.Size()) {
			saver.encodeChan <- encodeResult{err: fmt.Errorf("failed to write whole file to output: %d != %d", count, reader.Size())}
			return
		}

		saver.encodeChan <- encodeResult{}
		return
	}()

	return saver, nil
}

func (l *fileSaver) WaitResult() error {
	result := <-l.encodeChan
	return result.err
}
