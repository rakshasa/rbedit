package rbeditCmd

import (
	"bytes"

	"github.com/rakshasa/bencode-go"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

const (
	bencodeValueFlagName = "bencode"
	inplaceFlagName      = "inplace"
	inputFlagName        = "input"
	stringValueFlagName  = "string"
)

// Add command flags:

func addInputFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(inputFlagName, "i", "", "Input a file by path")
}

func addOutputFlags(cmd *cobra.Command) {
	cmd.Flags().Bool(inplaceFlagName, false, "Output inplace to input file")
}

func addAnyValueFlags(cmd *cobra.Command) {
	cmd.Flags().String(bencodeValueFlagName, "", "Bencoded value")
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

	object, err := bencode.Decode(bytes.NewReader([]byte(value)))
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
