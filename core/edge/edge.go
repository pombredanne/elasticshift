package edge

import (
	"context"
)

// Edge ..
// Defines a rest end point
type Edge func(ctx context.Context, req interface{}) (res interface{}, err error)
