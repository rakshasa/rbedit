package objects

import (
	"bufio"
	"fmt"
	"os"

	bencode "github.com/jackpal/bencode-go"
)

type Loader interface {
	WaitResult() (interface{}, error)
}

func WaitLoaderResult(l Loader) (interface{}, error) {
	if l == nil {
		return nil, fmt.Errorf("could not wait for loader results, loader not initalized")
	}

	return l.WaitResult()
}

type decodeResult struct {
	obj interface{}
	err error
}

type fileLoader struct {
	decodeChan chan decodeResult
}

func NewFileLoader(path string) (*fileLoader, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0444)
	if err != nil {
		return nil, fmt.Errorf("could not open file for bencode decoding: %v", err)
	}

	loader := &fileLoader{
		decodeChan: make(chan decodeResult),
	}

	go func() {
		defer file.Close()

		obj, err := bencode.Decode(bufio.NewReader(file))
		if err != nil {
			loader.decodeChan <- decodeResult{err: fmt.Errorf("failed to decode bencoded file: %v", err)}
			return
		}

		loader.decodeChan <- decodeResult{obj: obj}
		return
	}()

	return loader, nil
}

func (l *fileLoader) WaitResult() (interface{}, error) {
	result := <-l.decodeChan
	return result.obj, result.err
}
