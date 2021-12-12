package types

import (
	"fmt"
	"net/url"
)

func VerifyAbsoluteURI(uri string) bool {
	u, err := url.Parse(uri)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func EscapeURIString(str string) string {
	return string(EscapeURIBytes([]byte(str)))
}

func EscapeURIBytes(data []byte) []byte {
	escaped := make([]byte, 0, len(data)*2)

	for _, c := range data {
		if c < 0x20 || c >= 0x7f {
			escaped = append(escaped, []byte(fmt.Sprintf("\\x%02x", int(c)))...)
		} else {
			escaped = append(escaped, c)
		}
	}

	return escaped
}
