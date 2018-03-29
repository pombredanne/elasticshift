/*
Copyright 2018 The Elasticshift Authors.
*/
package ast

import (
	"strings"

	"gitlab.com/conspico/elasticshift/pkg/shiftfile/keys"
	"gitlab.com/conspico/elasticshift/pkg/shiftfile/scope"
)

func (f *File) Version() string {
	return f.value(scope.Ver)
}

func (f *File) From() string {
	return f.value(scope.Frm)
}

func (f *File) Name() string {
	return f.value(scope.Nam)
}

func (f *File) Language() string {
	return f.value(scope.Lan)
}

func (f *File) WorkDir() string {
	return f.value(scope.Dir)
}

// TODO GET method

func (f *File) Vars() map[string]string {

	vars := make(map[string]string)
	for _, item := range items(f.Node) {
		if item.Kind == scope.Var {
			vars[item.Keys[0].Key.Text] = item.Value.(*Literal).Token.Text
		}
	}
	return vars
}

func (f *File) Var(name string) string {

	for _, item := range items(f.Node) {
		if item.Kind == scope.Var && strings.EqualFold(item.Keys[0].Key.Text, name) {
			return item.Value.(*Literal).Token.Text
		}
	}
	return ""
}

func (f *File) Image() map[string]interface{} {

	var n *NodeItem
	iMap := make(map[string]interface{})
	for _, item := range items(f.Node) {
		if scope.Img == item.Kind {
			n = item
			break
		}
	}

	if n == nil {
		return iMap
	}

	// Image name
	iMap[keys.NAME] = n.Keys[0].Key.Text

	img := n.Value.(*Image)
	if img.Node == nil {
		return iMap
	}

	// read the properties
	f.properties(img.Node, iMap)

	return iMap
}

func (f *File) HasMoreBlocks() bool {
	return f.currentBlock < f.BlockCount
}

func (f *File) NextBlock() map[string]interface{} {

	if f.currentBlock == 0 && f.BlockCount > 0 {
		f.currentBlock = 1
	} else {
		f.currentBlock += 1
	}

	blk := make(map[string]interface{})
	for _, i := range items(f.Node) {

		if scope.Blk == i.Kind {

			n := i.Value.(*Block)

			if f.currentBlock == n.Number {

				// block name
				blk[keys.NAME] = i.Keys[0].Key.Text

				// block description
				blk[keys.DESC] = i.Keys[1].Key.Text

				// block properties
				f.properties(n, blk)

				break
			}
		}
	}

	blk[keys.BLOCK_NUMBER] = f.currentBlock

	return blk
}

func (f *File) properties(ni Node, props map[string]interface{}) map[string]interface{} {

	// hint map
	hmap := make(map[string]string)
	cmdSlice := []string{}

	for _, item := range ni.(*Block).Node {

		n := item.(*NodeItem)

		// key (identifier)
		key := n.Keys[0].Key.Text

		switch n.Value.(type) {

		case *Hint:
			hint := n.Value.(*Hint)
			hmap[hint.Operation] = hint.Value

		case *Literal:
			props[key] = n.Value.(*Literal).Token.Text

		case *List:

			listSlice := []string{}
			for _, i := range n.Value.(*List).Node {
				listSlice = append(listSlice, i.(*Literal).Token.Text)
			}
			props[key] = listSlice

		case *VarHolder:
			props[key] = f.Var(n.Value.(*VarHolder).Token.Text)

		case *Command:
			cmdSlice = append(cmdSlice, n.Value.(*Command).Token.Text)
		}
	}

	if len(hmap) > 0 {
		props[keys.HINT] = hmap
	}

	if len(cmdSlice) > 0 {
		props[keys.COMMAND] = cmdSlice
	}

	return props
}

func items(n Node) []*NodeItem {
	return n.(*NodeList).List
}

func (f *File) value(kind scope.NodeKind) string {

	for _, item := range items(f.Node) {
		if kind == item.Kind {
			return item.Keys[0].Key.Text
		}
	}
	return ""
}