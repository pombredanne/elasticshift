package util

import (
	"fmt"

	"github.com/nu7hatch/gouuid"
	"github.com/palantir/stacktrace"
)

// NewUUID ..
// Creates a new UUID and returns string
func NewUUID() (string, error) {
	u, err := uuid.NewV4()
	if err != nil {
		return "", stacktrace.Propagate(err, "Can't generate new UUID")
	}

	// trim the uuid without hyphen
	return fmt.Sprintf("%x", u[0:]), err
}
