package common

import (
	"fmt"
	"log"
	"strings"

	"github.com/rakshasa/bencode-go"
	"github.com/rakshasa/rbedit/objects"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

const (
	inputFlagName           = "input"
	inputBatchFlagName      = "batch"
	inputNotTorrentFlagName = "not-torrent"
	outputFlagName          = "output"
	outputInplaceFlagName   = "inplace"
	outputTemplateFlagName  = "output-template"

	bencodeValueFlagName = "bencode"
	integerValueFlagName = "int"
	jsonValueFlagName    = "json"
	stringValueFlagName  = "string"
)

// Add command flags:

func addInputFlags(cmd *cobra.Command) {
	cmd.Flags().VarP(&nonEmptyString{}, inputFlagName, "i", "Input filename")
	cmd.Flags().Bool(inputBatchFlagName, false, "Input as batch of filenames")
	cmd.Flags().Bool(inputNotTorrentFlagName, false, "Disable torrent verification on input")
}

func addDataOutputFlags(cmd *cobra.Command) {
	cmd.Flags().VarP(&nonEmptyString{}, outputFlagName, "o", "Output to filename")
	cmd.Flags().Var(&nonEmptyString{}, outputTemplateFlagName, "Output to template filename")
}

func addFileOutputFlags(cmd *cobra.Command) {
	addDataOutputFlags(cmd)
	cmd.Flags().Bool(outputInplaceFlagName, false, "Output to source filename, replacing it")
}

func addAnyValueFlags(cmd *cobra.Command) {
	cmd.Flags().Var(&nonEmptyString{}, bencodeValueFlagName, "Bencoded value")
	cmd.Flags().Int64(integerValueFlagName, 0, "Integer value")
	cmd.Flags().Var(&nonEmptyString{}, jsonValueFlagName, "JSON value")
	cmd.Flags().String(stringValueFlagName, "", "String value")
}

// Get value from flags:

func hasChangedFlags(cmd *cobra.Command) bool {
	var result bool

	cmd.Flags().Visit(func(f *flag.Flag) {
		result = true
	})

	return result
}

func getChangedString(flagSet *flag.FlagSet, flagName string) (string, bool) {
	if !flagSet.Changed(flagName) {
		return "", false
	}

	str, err := flagSet.GetString(flagName)
	if err != nil {
		log.Fatalf("getChangedString(flagSet, \"%s\") received unexpected error: %v", err)
	}

	return str, true
}

func getChangedTrue(flagSet *flag.FlagSet, flagName string) bool {
	if !flagSet.Changed(flagName) {
		return false
	}

	flag, err := flagSet.GetBool(flagName)
	if err != nil {
		log.Fatalf("getBoolString(flagSet, \"%s\") received unexpected error: %v", err)
	}

	return flag
}

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
