package rbeditCmd

import (
	"fmt"

	"github.com/rakshasa/rbedit/objects"
)

// String value:

type stringValue string

func (s *stringValue) String() string {
	return string(*s)
}

func (s *stringValue) Set(arg string) error {
	*s = stringValue(arg)
	return nil
}

func (s *stringValue) Type() string {
	return "string"
}

func (s *stringValue) Valid() bool {
	return len(s.String()) != 0
}

// URI value:

type uriValue string

func (s *uriValue) String() string {
	return string(*s)
}

func (s *uriValue) Set(arg string) error {
	if !objects.VerifyAbsoluteURI(arg) {
		return fmt.Errorf("failed to validate URI")
	}

	*s = uriValue(arg)
	return nil
}

func (s *uriValue) Type() string {
	return "uri"
}

func (s *uriValue) Valid() bool {
	return len(s.String()) != 0
}
