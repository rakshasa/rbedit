package common

import (
	"strings"

	"github.com/rakshasa/bencode-go"
	"github.com/rakshasa/rbedit/objects"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

const (
	bencodeValueFlagName = "bencode"
	integerValueFlagName = "int"
	inplaceFlagName      = "inplace"
	inputFlagName        = "input"
	jsonValueFlagName    = "json"
	outputFlagName       = "output"
	stringValueFlagName  = "string"
)

// Add command flags:

func addInputFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(inputFlagName, "i", "", "Input a file by path")
}

func addOutputFlags(cmd *cobra.Command) {
	cmd.Flags().Bool(inplaceFlagName, false, "Replace input file with output")
	cmd.Flags().String(outputFlagName, "", "Output to file")
}

func addAnyValueFlags(cmd *cobra.Command) {
	cmd.Flags().String(bencodeValueFlagName, "", "Bencoded value")
	cmd.Flags().Int64(integerValueFlagName, 0, "Integer value")
	cmd.Flags().String(jsonValueFlagName, "", "JSON value")
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
