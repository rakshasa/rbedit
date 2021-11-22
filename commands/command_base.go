package commands

import (
	"flag"
	"fmt"

	"github.com/rakshasa/rbedit/objects"
)

type CommandBase struct {
	loader      objects.Loader
	includeRoot bool
}

func (c *CommandBase) loadFile(path string) error {
	if c.loader != nil {
		return fmt.Errorf("could not load file, bencode loader already selected")
	}

	loader, err := objects.NewFileLoader(path)
	if err != nil {
		return err
	}

	c.loader = loader
	return nil
}

func (c *CommandBase) commonInputFlags(f *flag.FlagSet) {
	f.Func("file", "Input file", c.loadFile)
	f.BoolVar(&c.includeRoot, "include-root", false, "When passing no arguments, get root object")
}

func (c *CommandBase) lookupKeyPath(keys []string) (interface{}, error) {
	rootObj, err := objects.WaitLoaderResult(c.loader)
	if err != nil {
		return nil, &exitFailureError{msg: err.Error()}
	}

	obj, err := objects.LookupKeyPath(rootObj, keys)
	if err != nil {
		return nil, &exitFailureError{msg: err.Error()}
	}

	return obj, err
}
