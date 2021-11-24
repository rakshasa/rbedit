package commands

import (
	"flag"
	"fmt"

	"github.com/rakshasa/rbedit/objects"
)

type CommandBase struct {
	loader objects.Loader

	path  string
	value interface{}

	includeRoot bool
	inplace     bool
}

func (c *CommandBase) commonInputFlags(f *flag.FlagSet) {
	f.Func("file", "Input file", c.fileLoader)
	f.BoolVar(&c.includeRoot, "include-root", false, "When passing no arguments, get root object")
}

func (c *CommandBase) commonOutputFlags(f *flag.FlagSet) {
	f.BoolVar(&c.inplace, "inplace", false, "Inplace replace file on write")
}

// TODO: Convert to exitstatus.
func (c *CommandBase) fileLoader(path string) error {
	if c.loader != nil {
		return fmt.Errorf("failed to initialize, multiple bencode loaders selected")
	}

	loader, err := objects.NewFileLoader(path)
	if err != nil {
		return err
	}

	c.loader = loader
	c.path = path

	return nil
}

func (c *CommandBase) bencodeWriter(rootObj interface{}) (objects.Saver, ExitStatusError) {
	var path string

	if c.inplace {
		if len(c.path) == 0 {
			return nil, &exitUsageError{msg: "failed to save inplace, not loaded from a file"}
		}

		path = c.path
	} else {
		return nil, &exitUsageError{msg: "no output target selected"}
	}

	saver, err := objects.NewFileSaver(path, rootObj)
	if err != nil {
		return nil, &exitFailureError{msg: err.Error()}
	}

	return saver, nil
}

func (c *CommandBase) loadRootWithKeyPath(keys []string) (interface{}, interface{}, ExitStatusError) {
	rootObj, err := objects.WaitLoaderResult(c.loader)
	if err != nil {
		return nil, nil, &exitFailureError{msg: err.Error()}
	}

	obj, err := objects.LookupKeyPath(rootObj, keys)
	if err != nil {
		return nil, nil, &exitFailureError{msg: err.Error()}
	}

	return rootObj, obj, nil
}

func (c *CommandBase) saveRootWithKeyPath(rootObj, setObj interface{}, keys []string) (interface{}, ExitStatusError) {
	if setObj == nil {
		return nil, &exitFailureError{msg: "cannot save a key path with nil value"}
	}

	rootObj, err := objects.SetKeyPath(rootObj, setObj, keys)
	if err != nil {
		return nil, &exitFailureError{msg: err.Error()}
	}

	saver, statusErr := c.bencodeWriter(rootObj)
	if statusErr != nil {
		return nil, statusErr
	}

	if err := objects.WaitSaverResult(saver); err != nil {
		return nil, &exitFailureError{msg: fmt.Sprintf("failed to write encoded output: %v", err)}
	}

	return rootObj, nil
}

func (c *CommandBase) putString(str string) error {
	if c.value != nil {
		return &exitUsageError{msg: "cannot pass multiple value"}
	}

	c.value = str
	return nil
}
