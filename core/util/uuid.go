package util

import (
	"fmt"

	"github.com/nu7hatch/gouuid"
)

// NewUUID ..
// Creates a new UUID and returns string
func NewUUID() (string, error) {
	u, err := uuid.NewV4()
	if err != nil {
		return "", nil
	}

	// trim the uuid without hyphen
	return fmt.Sprintf("%x", u[0:]), err
}
