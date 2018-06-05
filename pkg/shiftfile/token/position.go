/*
Copyright 2018 The Elasticshift Authors.
*/
package token

import "fmt"

type Position struct {
	Filename string // Name of the configuration file
	Line     int    // Linenumber starts at 1
	Column   int    // Column starts at 1
	Offset   int    // offset, starting at 0
}

func (p *Position) IsValid() bool {
	return p.Line > 0
}

func (p Position) String() string {
	s := p.Filename
	if p.IsValid() {
		if s != "" {
			s += ":"
		}
		s += fmt.Sprintf("%d:%d", p.Line, p.Column)
	}

	if s == "" {
		s = "-"
	}
	return s
}
