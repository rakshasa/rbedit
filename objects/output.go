package objects

import (
	"bytes"
	"fmt"
	"os"

	bencode "github.com/rakshasa/bencode-go"
)

type Output interface {
	Execute(obj interface{}, path string) error
}

type outputWrite struct {
	path string
	obj  interface{}
}

type outputWait struct {
	path string
	err  error
}

type fileOutput struct {
	writeChan chan outputWrite
	waitChan  chan outputWait
}

func NewFileOutput() (*fileOutput, error) {
	output := &fileOutput{
		writeChan: make(chan outputWrite),
		waitChan:  make(chan outputWait),
	}

	go func() {
		pathObj, ok := <-output.writeChan
		if !ok {
			output.waitChan <- outputWait{path: pathObj.path, err: fmt.Errorf("failed to output file, channel is closed")}
			return
		}

		file, err := os.OpenFile(pathObj.path, os.O_WRONLY, 0666)
		if err != nil {
			output.waitChan <- outputWait{path: pathObj.path, err: fmt.Errorf("could not open file for bencode decoding: %v", err)}
			return
		}
		defer file.Close()

		var buf bytes.Buffer

		err = bencode.Marshal(&buf, pathObj.obj)
		if err != nil {
			output.waitChan <- outputWait{path: pathObj.path, err: fmt.Errorf("failed to encode data: %v", err)}
			return
		}

		reader := bytes.NewReader(buf.Bytes())

		count, err := file.ReadFrom(reader)
		if err != nil {
			output.waitChan <- outputWait{path: pathObj.path, err: fmt.Errorf("failed to write to output: %v", err)}
			return
		}
		if count != int64(reader.Size()) {
			output.waitChan <- outputWait{path: pathObj.path, err: fmt.Errorf("failed to write whole file to output: %d != %d", count, reader.Size())}
			return
		}

		output.waitChan <- outputWait{path: pathObj.path}

		close(output.waitChan)
		return
	}()

	return output, nil
}

func (f *fileOutput) Execute(obj interface{}, path string) error {
	f.writeChan <- outputWrite{path: path, obj: obj}

	pathObj, ok := <-f.waitChan
	if !ok {
		return fmt.Errorf("failed to output file, channel is closed")
	}
	if pathObj.err != nil {
		return pathObj.err
	}

	return nil
}
