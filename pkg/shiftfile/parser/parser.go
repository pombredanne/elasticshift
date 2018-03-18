/*
Copyright 2018 The Elasticshift Authors.
*/
package parser

import (
	"bytes"
	"errors"
	"fmt"

	"gitlab.com/conspico/elasticshift/pkg/shiftfile/scanner"
	"gitlab.com/conspico/elasticshift/pkg/shiftfile/scope"

	"gitlab.com/conspico/elasticshift/pkg/shiftfile/ast"
	"gitlab.com/conspico/elasticshift/pkg/shiftfile/token"
)

var errEofToken = errors.New("EOF token found")

type Parser struct {
	s *scanner.Scanner

	tok     token.Token
	prevTok token.Token

	comments []*ast.Comment

	leadComment []*ast.Comment
	lineComment []*ast.Comment

	f *ast.File

	tokenScanned bool

	ntype scope.NodeKind
	serr  error
}

func New(src []byte) *Parser {

	p := &Parser{}

	// Replace the \r\n to \n to avoid incorrect behavior
	// because the parse works with \n only as line endings.
	src = bytes.Replace(src, []byte("\r\n"), []byte("\n"), -1)

	errFunc := func(pos token.Position, msg string) {
		p.serr = &PositionErr{Position: pos, Err: errors.New(msg)}
	}
	// Initiatest the new scanner
	p.s = scanner.New(src, errFunc)

	return p
}

func (p *Parser) Parse() (*ast.File, error) {

	p.f = &ast.File{}

	list, err := p.NodeList()

	if p.serr != nil {
		return nil, p.serr
	}

	if err != nil {
		return nil, err
	}

	p.f.Node = list
	p.f.Comments = p.comments

	return p.f, nil
}

func (p *Parser) NodeList() (*ast.NodeList, error) {

	root := &ast.NodeList{}

	for p.s.HasMoreTokens() {

		tok := p.scan()

		if tok.Type == token.EOF {
			break // parsing reached eof
		}

		n, err := p.nodeItem()
		if err != nil && err == errEofToken {
			break
		}

		if err != nil {
			return root, err
		}

		root.Add(n)
	}

	// fmt.Println(fmt.Sprintf("Node list: %q", root))

	return root, nil
}

func (p *Parser) scan() token.Token {

	if p.tokenScanned {
		p.tokenScanned = false
		return p.prevTok
	}

	p.prevTok = p.tok
	p.tok = p.s.Scan()

	// fmt.Printf("Scanned %s", p.tok)

	for token.COMMENT == p.tok.Type {

		// if previous token is on the same line as comment
		// then it might be a line comment
		comment := p.grabComment()

		// fmt.Println(comment)
		if p.tok.Position.Line == p.prevTok.Position.Line {
			p.lineComment = append(p.lineComment, comment)
		} else if p.tok.Position.Line != p.prevTok.Position.Line {
			p.leadComment = append(p.leadComment, comment)
		}

		p.prevTok = p.tok
		p.tok = p.s.Scan()
		p.tokenScanned = true
	}

	// fmt.Println("return token" + p.tok.Text)

	return p.tok
}

func (p *Parser) nodeItem() (*ast.NodeItem, error) {

	keys, err := p.nodeKey()
	if len(keys) > 0 && err != nil {

		// there are some keys availabe, but unfortunately there is also an err
		// ignore the error and proceed to get values
		err = nil
	}

	if len(keys) == 0 && err == errEofToken && token.COMMENT == p.prevTok.Type {

		// Looks like this is a orphan comment, no node is associated with it.
		err = nil
	}

	if err != nil {
		return nil, err
	}

	n := &ast.NodeItem{
		Keys: keys,
		Kind: p.ntype,
	}

	if p.leadComment != nil {
		n.LeadComments = p.leadComment
		p.leadComment = nil
	}

	// fmt.Println(fmt.Sprintf("node item = %q", p.tok))

	switch p.tok.Type {
	case token.VERSION, token.NAME, token.LANGUAGE, token.WORKDIR, token.COMMAND:
		n.Value, err = p.literal()
	case token.LBRACE:
		n.Value, err = p.nodeItem()
	// case token.IMAGE:
	// 	n.Value, err = p.image()
	case token.LPAREN:
		n.Value, err = p.varholder()
	case token.HINT:
		n.Value, err = p.hint()
	// case token.RBRACE:
	// 	break
	// case token.LBRACK:
	// 	n.Value, err = p.properties()
	default:
		switch p.prevTok.Type {
		case token.VAR:
			n.Value, err = p.literal()
			// case token.LPAREN:
			// 	n.Value, err = p.varholder()
		}
	}

	if err != nil {
		return nil, err
	}

	// reset the slice
	p.leadComment = nil
	p.lineComment = nil

	// fmt.Println(fmt.Sprintf("Node: %#v", n))

	return n, nil
}

// Returns the node key..
// Applicatble for VAR and BLOCK
func (p *Parser) nodeKey() ([]*ast.NodeKey, error) {

	keycount := 0
	keys := make([]*ast.NodeKey, 0)

	// reset the kind
	p.kind(0)

	// VAR, IMAGE, BLOCK
	for {

		// scan the next key token
		p.scan()
		// fmt.Println("Key : " + p.tok.Text)

		switch p.tok.Type {
		case token.EOF:
			return keys, errEofToken
		case token.IDENTIFIER:
			keys = append(keys, &ast.NodeKey{Key: p.tok})
			return keys, nil
		case token.STRING:
			keys = append(keys, &ast.NodeKey{Key: p.tok})
			keycount++
		case token.VAR:
			p.kind(scope.Var)
			keys = append(keys, &ast.NodeKey{Key: p.tok})
			return keys, nil
		case token.WORKDIR:
			p.kind(scope.Dir)
			keys = append(keys, &ast.NodeKey{Key: p.tok})
			return keys, nil
		case token.LANGUAGE:
			p.kind(scope.Lan)
			keys = append(keys, &ast.NodeKey{Key: p.tok})
			return keys, nil
		case token.FROM:
			p.kind(scope.Frm)
			keys = append(keys, &ast.NodeKey{Key: p.tok})
			return keys, nil
		case token.VERSION:
			p.kind(scope.Ver)
			keys = append(keys, &ast.NodeKey{Key: p.tok})
			return keys, nil
		case token.NAME:
			p.kind(scope.Nam)
			keys = append(keys, &ast.NodeKey{Key: p.tok})
			return keys, nil
		case token.COMMA:
			p.kind(scope.Blk)
		case token.IMAGE:
			p.kind(scope.Img)
		case token.HINT:
			p.kind(scope.Hin)
			keys = append(keys, &ast.NodeKey{Key: p.tok})
			return keys, nil
		case token.COMMAND:
			p.kind(scope.Cmd)
			keys = append(keys, &ast.NodeKey{Key: p.tok})
			return keys, nil
		case token.LBRACE:
			if keycount == 0 {
				return keys, &PositionErr{
					Position: p.tok.Position,
					Err:      fmt.Errorf("Expected token: STRING | IDENTIFIER, got: %s", p.tok.Type),
				}
			}
			p.kind(scope.Blk)
			return keys, nil
		case token.ILLEGAL:
			return keys, &PositionErr{
				Position: p.tok.Position,
				Err:      fmt.Errorf("Illegal character"),
			}
		case token.LPAREN:
			p.kind(scope.Vhl)
			return keys, nil
		case token.RBRACE:
			break
		default:
			return keys, &PositionErr{
				Position: p.tok.Position,
				Err:      fmt.Errorf("Expected token: STRING | IDENTIFIER | VAR | LBRACE, got: %s", p.tok.Type),
			}
		}
	}
}

func (p *Parser) appendKey(keys []*ast.NodeKey, tok token.Token) {
	keys = append(keys, &ast.NodeKey{Key: tok})
}

func (p *Parser) kind(ntype scope.NodeKind) {
	if p.ntype == 0 {
		p.ntype = ntype
	}
}

func (p *Parser) grabComment() *ast.Comment {

	comment := &ast.Comment{}
	comment.Start = p.tok.Position
	comment.Value = p.tok.Text

	return comment
}

func (p *Parser) literal() (*ast.Literal, error) {

	// reads the literal value
	p.scan()

	lit := &ast.Literal{}
	lit.Token = p.tok
	return lit, nil
}

func (p *Parser) varholder() (*ast.VarHolder, error) {

	// reads the variable placeholder
	p.scan()

	vh := &ast.VarHolder{}
	vh.Token = p.tok

	p.scan()

	if p.tok.Type == token.RPAREN {
		return vh, nil
	}
	return nil, fmt.Errorf("Expected: RBRACK ')' but %v", p.tok.Text)
}

func (p *Parser) hint() (*ast.Hint, error) {

	// after comment is scanned, the next token is buffered
	// set it to false, so that it can read the hint
	if p.prevTok.Type == token.COMMENT {
		p.tokenScanned = false
	}

	// reads the hint operation
	p.scan()

	if p.tok.Type == token.IDENTIFIER {

		//TODO validate the identifier
	} else {
		return nil, fmt.Errorf("Expected: IDENTIFIER but got: %s", p.tok.Type)
	}

	hint := &ast.Hint{}
	hint.Operation = p.tok.Text

	// read the hint delimiter
	p.scan()

	if p.tok.Type != token.HINT_DEL {
		return nil, fmt.Errorf("Expected a Hint delimiter ':' but got %s", p.tok.Type)
	}

	// read the hint operation value
	p.scan()

	if p.tok.Type != token.IDENTIFIER {
		return nil, fmt.Errorf("Expected an identifier but got %s", p.tok.Type)
	}

	hint.Value = p.tok.Text

	return hint, nil
}

func (p *Parser) validate() bool {

	return true
	// return nil, errors.New("Expected: STRING token (the name of the shiftfile, format :'company/name')")
}

// func (p *Parser) variable() (*ast.Variable, error) {

// 	vari := &ast.Variable{}
// 	vari.Start = p.tok.Position

// 	// // next token should be the name of the variable
// 	// p.scan()

// 	// if token.IDENTIFIER == p.tok.Type {
// 	// 	vari.Name = p.tok.Text
// 	// } else {
// 	// 	return nil, errors.New("Expected: IDENTIFIER token (the name of the variable)")
// 	// }

// 	// next token should be the value of the variable
// 	p.scan()

// 	if token.STRING == p.tok.Type {
// 		vari.Value = p.tok.Text
// 	} else {
// 		return nil, fmt.Errorf("Expected: STRING token (the value belongs to the variable '%s')", vari.Name)
// 	}

// 	return vari, nil
// }

// func (p *Parser) grabVersion() (ast.Node, error) {

// 	literal := &ast.Literal{}
// 	literal.Token = p.tok

// 	// n := ast.NewNodeItem(scope.Ver, &ast.NodeKey{p.tok})

// 	// next token should be the STRING type denotes the actual version
// 	p.scan()

// 	if token.STRING == p.tok.Type {
// 		literal.Value = p.tok.Text
// 	} else {
// 		return nil, errors.New("Expected: STRING token (the value of the version)")
// 	}

// 	return literal, nil
// }
// func (p *Parser) grabLanguage() (*ast.Object, error) {

// 	lang := &ast.Language{}
// 	lang.Start = p.tok.Position
// 	lang.Value = p.tok.Text

// 	obj := ast.NewObject(ast.Lan)

// 	// next token should be the STRING type denotes the actual version
// 	p.scan()

// 	if token.STRING == p.tok.Type {
// 		lang.Value = p.tok.Text
// 	} else {
// 		return nil, errors.New("Expected: STRING token (the language used for this project)")
// 	}

// 	obj.Val = lang
// 	obj.AddKey(p.tok)

// 	p.f.Language = lang

// 	return obj, nil
// }

// func (p *Parser) grabWorkDirectory() (*ast.Object, error) {

// 	wdir := &ast.Workdir{}
// 	wdir.Start = p.tok.Position
// 	wdir.Value = p.tok.Text

// 	obj := ast.NewObject(ast.Dir)

// 	// next token should be the STRING type denotes the actual version
// 	p.scan()

// 	if token.STRING == p.tok.Type {
// 		wdir.Value = p.tok.Text
// 	} else {
// 		return nil, errors.New("Expected: STRING token (the working directory of this buil)")
// 	}

// 	obj.Val = wdir
// 	obj.AddKey(p.tok)

// 	p.f.Workdir = wdir

// 	return obj, nil
// }
