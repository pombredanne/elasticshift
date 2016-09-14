package util

import (
	"regexp"
)

// IsAlphaNumericOnly ..
// Check to see if the given text is alpha-numeric only
func IsAlphaNumericOnly(str string) (bool, error) {
	return regexp.MatchString("^[A-Za-z0-9]*$", str)
}
