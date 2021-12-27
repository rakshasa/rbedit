package inputs

import (
	"github.com/rakshasa/rbedit/types"
)

// SequentialBatchInput:

type sequentialBatchInput struct {
	decodeFn types.DecodeFunc
	inputFn  types.InputFunc
}

func NewSequentialBatchInput(decodeFn types.DecodeFunc, inputFn types.InputFunc) *sequentialBatchInput {
	return &sequentialBatchInput{
		decodeFn: decodeFn,
		inputFn:  inputFn,
	}
}

func (o *sequentialBatchInput) Execute(metadata types.IOMetadata, resultFn types.InputResultFunc) error {
	for {
		metadata, data, err := o.inputFn(metadata)
		if err != nil {
			return err
		}
		if data == nil {
			break
		}

		object, err := o.decodeFn(data)
		if err != nil {
			return err
		}

		if err := resultFn(object, metadata); err != nil {
			return err
		}
	}

	return nil
}

// ParallelBatchInput:

// TODO: Slow on macos, test on linux.

type parallelBatchInput struct {
	decodeFn types.DecodeFunc
	inputFn  types.InputFunc
}

func NewParallelBatchInput(decodeFn types.DecodeFunc, inputFn types.InputFunc) *parallelBatchInput {
	return &parallelBatchInput{
		decodeFn: decodeFn,
		inputFn:  inputFn,
	}
}

func (o *parallelBatchInput) Execute(metadata types.IOMetadata, resultFn types.InputResultFunc) error {
	taskChan := make(chan string, 32)
	defer close(taskChan)

	resultChan := make(chan error, 32)
	defer close(resultChan)

	readingBatch := true
	taskCount := 0
	workerCount := 0

	for readingBatch || taskCount != 0 {
		select {
		case err := <-resultChan:
			if err != nil {
				return err
			}
			taskCount--
		default:
		}

		if !readingBatch {
			continue
		}

		metadata, filename, err := o.inputFn(metadata)
		if err != nil {
			return err
		}
		if filename == nil {
			readingBatch = false
			continue
		}

		taskChan <- string(filename)
		taskCount++

		if taskCount > workerCount && workerCount < 10 {
			go func() {
				for filename := range taskChan {
					metadata, data, err := readFileInput(metadata, string(filename))
					if err != nil {
						resultChan <- err
						return
					}

					object, err := o.decodeFn(data)
					if err != nil {
						resultChan <- err
						return
					}

					if err := resultFn(object, metadata); err != nil {
						resultChan <- err
						return
					}

					resultChan <- nil
				}
			}()

			workerCount++
		}
	}

	return nil
}
