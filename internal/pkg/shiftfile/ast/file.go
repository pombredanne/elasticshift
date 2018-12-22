/*
Copyright 2018 The Elasticshift Authors.
*/
package ast

import (
	"strings"

	"github.com/elasticshift/elasticshift/internal/pkg/shiftfile/keys"
	"github.com/elasticshift/elasticshift/internal/pkg/shiftfile/scope"
)

var (
	PREFIX_SECRET   = "^"
	PREFIX_ARGUMENT = "@"
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
	return f.value(scope.Wdi)
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

func (f *File) ImageNames() []string {

	var n *NodeItem
	for _, item := range items(f.Node) {
		if scope.Img == item.Kind {
			n = item
			break
		}
	}

	if n == nil {
		return nil
	}

	keys := []string{}
	for _, k := range n.Keys {
		keys = append(keys, k.Key.Text)
	}
	return keys
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

	nkeys := len(n.Keys)
	if nkeys > 1 {
		keyarr := []string{}
		for _, k := range n.Keys {
			keyarr = append(keyarr, k.Key.Text)
		}
		iMap[keys.NAME] = strings.Join(keyarr, ",")
	} else {
		iMap[keys.NAME] = n.Keys[0].Key.Text
	}

	// Image name

	var img *Image
	if n.Value != nil {
		img = n.Value.(*Image)
	}

	if img == nil || img.Node == nil {
		return iMap
	}

	// read the properties
	f.properties(img.Node, iMap)

	return iMap
}

func (f *File) CacheDirectories() []string {

	var n *NodeItem
	for _, item := range items(f.Node) {
		if scope.Cac == item.Kind {
			n = item
			break
		}
	}

	if n == nil {
		return nil
	}

	directories := []string{}

	for _, item := range n.Value.(*Cache).Node.(*Block).Node {

		n := item.(*NodeItem)

		switch n.Value.(type) {
		case *Directory:
			directories = append(directories, n.Value.(*Directory).Token.Text)
		}

	}

	return directories
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

		case *Argument:
			props[key] = PREFIX_ARGUMENT + n.Value.(*Argument).Token.Text

		case *Secret:
			props[key] = PREFIX_SECRET + n.Value.(*Secret).Token.Text
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

func (f *File) IsSecretOrArgument(value string) bool {
	return value != "" && (strings.HasPrefix(value, PREFIX_SECRET) || strings.HasPrefix(value, PREFIX_ARGUMENT))
}

func items(n Node) []*NodeItem {
	return n.(*NodeList).List
}

func (f *File) value(kind scope.NodeKind) string {

	var n *NodeItem
	for _, item := range items(f.Node) {
		if kind == item.Kind {
			n = item
			break
		}
	}

	if n == nil {
		return ""
	}

	return n.Value.(*Literal).Token.Text
}
