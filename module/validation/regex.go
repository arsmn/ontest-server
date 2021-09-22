package validation

import "regexp"

var (
	UsernameRegex = regexp.MustCompile("^[A-Za-z][A-Za-z0-9_]*$")
)
