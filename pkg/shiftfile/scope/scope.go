package scope

type NodeKind int

const (
	Bad NodeKind = iota // error handling
	Ver
	Frm
	Nam
	Dir
	Img
	Var
	Lan
	Blk
	Prp
	Vhl
	Hin
	Cmd
)

var nodeKindStrings = [...]string{
	Bad: "bad",
	Ver: "VERSION",
	Frm: "FROM",
	Nam: "NAME",
	Dir: "DIRECTORY",
	Img: "IMAGE",
	Var: "VAR",
	Lan: "LANGUAGE",
	Blk: "BLOCK",
	Prp: "PROPERTY",
	Vhl: "VARHOLDER",
	Hin: "HINT",
	Cmd: "COMMAND",
}

func (k NodeKind) String() string {
	return nodeKindStrings[k]
}
