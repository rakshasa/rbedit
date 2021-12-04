package objects

import (
	"bufio"
	"fmt"
	"os"

	bencode "github.com/rakshasa/bencode-go"
)

type Input interface {
	// Executed once for every distinct root bencoded data object in
	// the input.
	Execute(fn func(obj interface{}) error) error
}

type inputObjectError struct {
	obj interface{}
	err error
}

type fileInput struct {
	waitChan chan inputObjectError
}

func NewFileInput(path string) (*fileInput, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0444)
	if err != nil {
		return nil, fmt.Errorf("could not open file for bencode decoding: %v", err)
	}

	input := &fileInput{
		waitChan: make(chan inputObjectError),
	}

	go func() {
		defer file.Close()

		obj, err := bencode.Decode(bufio.NewReader(file))
		if err != nil {
			input.waitChan <- inputObjectError{err: fmt.Errorf("failed to decode bencoded file: %v", err)}
			return
		}

		input.waitChan <- inputObjectError{obj: obj}

		close(input.waitChan)
		return
	}()

	return input, nil
}

func (f *fileInput) Execute(fn func(obj interface{}) error) error {
	objErr, ok := <-f.waitChan
	if !ok {
		return fmt.Errorf("failed to execute, channel is closed")
	}
	if objErr.err != nil {
		return fmt.Errorf("failed to execute: %v", objErr.err)
	}

	return fn(objErr.obj)
}
