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
			return item.Keys[0].Key.Text
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
	iMap[keys.IMAGE_NAME] = n.Keys[0].Key.Text

	if n.Value == nil {
		return iMap
	}

	// read the properties
	props := f.properties(n.Value)

	iMap[keys.PROPERTIES] = props

	return iMap
}

func (p *File) properties(ni Node) map[string]interface{} {

	props := make(map[string]interface{})
	for _, item := range ni.(*Block).Node {

		n := item.(*NodeItem)

		// key (identifier)
		key := n.Keys[0].Key.Text

		var val interface{}
		switch n.Value.(type) {

		case *Hint:
		case *Literal:
		case *List:
		}

		props[key] = val
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
