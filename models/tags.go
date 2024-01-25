package models

import (
	"net/url"
)

type Tags []Tag

type Tag struct {
	Name     string
	Count    int
	Selected bool
}

// TODO: Extend this method to disambiguate tags with spaces from tags with pluses
// e.g. "react native" vs "react+native"
// e.g. transform "react native" to "react__native"
func (t Tag) GetEscapedName() string {
	return url.QueryEscape(t.Name)
}
