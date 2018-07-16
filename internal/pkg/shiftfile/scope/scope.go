/*
Copyright 2018 The Elasticshift Authors.
*/
package scope

type NodeKind int

const (
	Bad NodeKind = iota // error handling
	Ver
	Frm
	Nam
	Wdi
	Img
	Var
	Lan
	Blk
	Prp
	Vhl
	Hin
	Cmd
	Cac
	Dir
)

var nodeKindStrings = [...]string{
	Bad: "bad",
	Ver: "VERSION",
	Frm: "FROM",
	Nam: "NAME",
	Wdi: "WORKDIR",
	Img: "IMAGE",
	Var: "VAR",
	Lan: "LANGUAGE",
	Blk: "BLOCK",
	Prp: "PROPERTY",
	Vhl: "VARHOLDER",
	Hin: "HINT",
	Cmd: "COMMAND",
	Cac: "CACHE",
	Dir: "DIRECTORY",
}

func (k NodeKind) String() string {
	return nodeKindStrings[k]
}
