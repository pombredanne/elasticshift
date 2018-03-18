package parser

import (
	"fmt"

	"gitlab.com/conspico/elasticshift/pkg/shiftfile/token"
)

type PositionErr struct {
	Position token.Position
	Err      error
}

func (pe *PositionErr) Error() string {
	return fmt.Sprintf("At %s:%s", pe.Position, pe.Err)
}
