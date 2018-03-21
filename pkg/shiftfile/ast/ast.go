/*
Copyright 2018 The Elasticshift Authors.
*/
package ast

import (
	"gitlab.com/conspico/elasticshift/pkg/shiftfile/scope"
	"gitlab.com/conspico/elasticshift/pkg/shiftfile/token"
)

type Node interface {
	node()
	Position() token.Position
}

func (n NodeList) node()  {}
func (n NodeKey) node()   {}
func (n NodeItem) node()  {}
func (b Block) node()     {}
func (l List) node()      {}
func (i Literal) node()   {}
func (h Hint) node()      {}
func (i Image) node()     {}
func (v VarHolder) node() {}

type Command Literal

func (c Command) node() {}

func (c *Command) Position() token.Position {
	return c.Token.Position
}

type File struct {
	Node     Node
	Comments []*Comment
}

func (f *File) Position() token.Position {
	return f.Node.Position()
}

type NodeKey struct {
	Key token.Token
}

type NodeItem struct {
	Kind scope.NodeKind
	Keys []*NodeKey

	Value Node

	LeadComments []*Comment
	LineComments []*Comment
}

func NewNodeItem(kind scope.NodeKind, key *NodeKey) *NodeItem {
	ni := &NodeItem{Kind: kind}
	ni.Keys = append(ni.Keys, key)
	return ni
}

func (n *NodeItem) Position() token.Position {
	return n.Value.Position()
}

func (n *NodeItem) AddKey(key *NodeKey) {
	n.Keys = append(n.Keys, key)
}

type NodeList struct {
	List []*NodeItem
}

func (n *NodeList) Position() token.Position {

	if len(n.List) == 0 {
		return token.Position{}
	}
	return n.List[0].Position()
}

func (n *NodeList) Add(item *NodeItem) {
	n.List = append(n.List, item)
}

// ChildNodes returns the items which have sub-items
func (n *NodeList) ChildNodes() *NodeList {
	var result NodeList
	for _, item := range n.List {
		if len(item.Keys) > 0 {
			result.Add(item)
		}
	}
	return &result
}

// Nodes returns only the items which are at same level
func (n *NodeList) Nodes() *NodeList {
	var result *NodeList
	for _, item := range n.List {
		if len(item.Keys) == 0 {
			result.Add(item)
		}
	}
	return result
}

type Block struct {
	Lbrace token.Position // {
	Rbrace token.Position // }
	Node   []Node
}

func (b *Block) Position() token.Position {
	return b.Lbrace
}

type Comment struct {
	Start token.Position
	Value string
}

func (c *Comment) Position() token.Position {
	return c.Start
}

type Literal struct {
	Token token.Token

	LeadComments []*Comment
	LineComments []*Comment
}

func (i *Literal) Position() token.Position {
	return i.Token.Position
}

type List struct {
	Lbrack token.Position // "["
	RBrack token.Position // "]"
	Node   []Node
}

func (l *List) Position() token.Position {
	return l.Lbrack
}

type Hint struct {
	Token     token.Token
	Operation string
	Value     string
}

func (h *Hint) Position() token.Position {
	return h.Token.Position
}

type Image struct {
	Start token.Position
	Node  Node
}

func (i *Image) Position() token.Position {
	return i.Start
}

type VarHolder struct {
	Token token.Token
}

func (v *VarHolder) Position() token.Position {
	return v.Token.Position
}
