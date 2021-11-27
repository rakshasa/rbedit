package objects

import (
	"net/url"
)

func VerifyAbsoluteURI(uri string) bool {
	u, err := url.Parse(uri)
	return err == nil && u.Scheme != "" && u.Host != ""
}
