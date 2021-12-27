package common

import (
	"fmt"
	"strings"

	"github.com/rakshasa/bencode-go"
	"github.com/rakshasa/rbedit/objects"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

const (
	inputFlagName         = "input"
	inputBatchFlagName    = "batch"
	outputInplaceFlagName = "inplace"
	outputFlagName        = "output"

	bencodeValueFlagName = "bencode"
	integerValueFlagName = "int"
	jsonValueFlagName    = "json"
	stringValueFlagName  = "string"
)

// Add command flags:

func addInputFlags(cmd *cobra.Command) {
	cmd.Flags().VarP(&nonEmptyString{}, inputFlagName, "i", "Input source")
	cmd.Flags().Bool(inputBatchFlagName, false, "Input as batch of filenames")
}

func addOutputFlags(cmd *cobra.Command) {
	cmd.Flags().Bool(outputInplaceFlagName, false, "Output to source file, replacing it")
	cmd.Flags().VarP(&nonEmptyString{}, outputFlagName, "o", "Output to file")
}

func addAnyValueFlags(cmd *cobra.Command) {
	cmd.Flags().Var(&nonEmptyString{}, bencodeValueFlagName, "Bencoded value")
	cmd.Flags().Int64(integerValueFlagName, 0, "Integer value")
	cmd.Flags().Var(&nonEmptyString{}, jsonValueFlagName, "JSON value")
	cmd.Flags().String(stringValueFlagName, "", "String value")
}

// Get value from flags:

func getBencodeValueFromFlag(flags *flag.FlagSet, name string) (interface{}, bool, error) {
	if !flags.Changed(name) {
		return nil, false, nil
	}

	value, err := flags.GetString(name)
	if err != nil {
		return nil, false, err
	}

	object, err := bencode.Decode(strings.NewReader(value))
	if err != nil {
		return nil, false, err
	}

	return object, true, nil
}

func getIntegerValueFromFlag(flags *flag.FlagSet, name string) (interface{}, bool, error) {
	if !flags.Changed(name) {
		return nil, false, nil
	}

	value, err := flags.GetInt64(name)
	if err != nil {
		return nil, false, err
	}

	return value, true, nil
}

func getJSONValueFromFlag(flags *flag.FlagSet, name string) (interface{}, bool, error) {
	if !flags.Changed(name) {
		return nil, false, nil
	}

	data, err := flags.GetString(name)
	if err != nil {
		return nil, false, err
	}

	object, err := objects.ConvertJSONToBencodeObject(data)
	if err != nil {
		return nil, false, err
	}

	return object, true, nil
}

func getStringValueFromFlag(flags *flag.FlagSet, name string) (interface{}, bool, error) {
	if !flags.Changed(name) {
		return nil, false, nil
	}

	value, err := flags.GetString(name)
	if err != nil {
		return nil, false, err
	}

	return value, true, nil
}

// Flag types:

type nonEmptyString struct {
	value string
}

func (v *nonEmptyString) String() string {
	return string(v.value)
}

func (v *nonEmptyString) Set(s string) error {
	if len(v.value) != 0 {
		return fmt.Errorf("duplicate flags")
	}
	if len(s) == 0 {
		return fmt.Errorf("empty string")
	}

	*v = nonEmptyString{value: s}
	return nil
}

func (v *nonEmptyString) Type() string {
	return "string"
}
