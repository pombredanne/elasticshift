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
	//return fmt.Sprintf("%x%x%x%x%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:]), err
	return fmt.Sprintf("%x", u[0:]), err
}
