/*
Copyright 2018 The Elasticshift Authors.
*/
package parser

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"gitlab.com/conspico/elasticshift/pkg/shiftfile/scanner"
	"gitlab.com/conspico/elasticshift/pkg/shiftfile/scope"

	"gitlab.com/conspico/elasticshift/pkg/shiftfile/ast"
	"gitlab.com/conspico/elasticshift/pkg/shiftfile/token"
)

var errEofToken = errors.New("EOF token found")

var (
	validHints = []string{"PARALLEL", "TIMEOUT"}
)

type Parser struct {
	s *scanner.Scanner

	tok     token.Token
	prevTok token.Token

	comments []*ast.Comment

	leadComment []*ast.Comment
	lineComment []*ast.Comment

	f *ast.File

	tokenScanned bool

	ntype  scope.NodeKind // node
	cscope scope.NodeKind // section

	serr error
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

	list, err := p.nodeList()

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

func AST(src []byte) (*ast.File, error) {

	p := New(src)
	return p.Parse()
}

func (p *Parser) nodeList() (*ast.NodeList, error) {

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

	for token.COMMENT == p.tok.Type {

		// if previous token is on the same line as comment
		// then it might be a line comment
		comment := p.grabComment()

		if p.tok.Position.Line == p.prevTok.Position.Line {
			p.lineComment = append(p.lineComment, comment)
		} else if p.tok.Position.Line != p.prevTok.Position.Line {
			p.leadComment = append(p.leadComment, comment)
		}

		p.prevTok = p.tok
		p.tok = p.s.Scan()
		p.unscan()
	}

	return p.tok
}

func (p *Parser) unscan() {
	p.tokenScanned = true
}

func (p *Parser) forceNextScan() {
	p.tokenScanned = false
}

func (p *Parser) nodeItem() (*ast.NodeItem, error) {

	// reset the kind/nodetype
	p.ntype = 0

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

	switch n.Kind {
	case scope.Var, scope.Ver, scope.Nam, scope.Lan, scope.Dir, scope.Frm:
		p.scan()
		n.Value, err = p.literal()
	case scope.Cmd:
		// p.unscan()
		n.Value, err = p.command()
	case scope.Blk:
		n.Value, err = p.block()
	case scope.Img:
		n.Value, err = p.image()
	case scope.Hin:
		n.Value, err = p.hint()
	case scope.Vhl:
		n.Value, err = p.varholder()
	case scope.Prp:
	recheck:
		switch p.tok.Type {
		case token.HINT:
			n.Value, err = p.hint()
		case token.LPAREN:
			n.Value, err = p.varholder()
		case token.LBRACK:
			n.Value, err = p.list()
		case token.STRING:
			n.Value, err = p.literal()
		case token.IDENTIFIER:
			p.scan()
			goto recheck
		}
	}

	if err != nil {
		return nil, err
	}

	// TODO add line comment

	// reset the comment
	p.leadComment = nil
	p.lineComment = nil

	return n, nil
}

// Returns the node key..
// Applicatble for VAR and BLOCK
func (p *Parser) nodeKey() ([]*ast.NodeKey, error) {

	p.ntype = 0

	keycount := 0
	keys := make([]*ast.NodeKey, 0)

	for {

		switch p.tok.Type {
		case token.EOF:
			return keys, errEofToken
		case token.VAR:
			p.kind(scope.Var)
		case token.VERSION:
			p.kind(scope.Ver)
			p.forceNextScan()
			goto exit
		case token.HINT:
			p.kind(scope.Hin)
			p.forceNextScan()
			goto exit
		case token.COMMAND:
			p.kind(scope.Cmd)
			goto exit
		case token.IMAGE:
			p.kind(scope.Img)
		case token.WORKDIR:
			p.kind(scope.Dir)
			p.forceNextScan()
			goto exit
		case token.LANGUAGE:
			p.kind(scope.Lan)
			p.forceNextScan()
			goto exit
		case token.FROM:
			p.kind(scope.Frm)
			p.forceNextScan()
			goto exit
		case token.NAME:
			p.kind(scope.Nam)
			p.forceNextScan()
			goto exit
		case token.IDENTIFIER:
			if p.cscope == scope.Img || p.cscope == scope.Blk {
				p.kind(scope.Prp)
			}
			p.forceNextScan() // avoid buffer
			goto exit
		case token.STRING:
			keys = append(keys, &ast.NodeKey{Key: p.tok})
			p.forceNextScan() // avoid buffer
			keycount++
		case token.COMMA:
		case token.LBRACE:
			if keycount == 0 {
				return keys, &PositionErr{
					Position: p.tok.Position,
					Err:      fmt.Errorf("Expected token: STRING | IDENTIFIER, got: %s", p.tok.Type),
				}
			}
			p.kind(scope.Blk)
			return keys, nil
		case token.LPAREN:
			p.kind(scope.Vhl)
			return keys, nil
		case token.LBRACK:
			return keys, nil
		case token.RBRACE, token.RBRACK, token.RPAREN:
			break
		case token.ILLEGAL:
			return keys, &PositionErr{
				Position: p.tok.Position,
				Err:      fmt.Errorf("Illegal character : %v", p.tok),
			}
		default:
			return keys, &PositionErr{
				Position: p.tok.Position,
				Err:      fmt.Errorf("Expected token: STRING | IDENTIFIER | VAR | LBRACE, got: %s", p.tok.Type),
			}
		}

		// next token
		p.scan()
	}

exit:
	keys = append(keys, &ast.NodeKey{Key: p.tok})
	return keys, nil
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

func (p *Parser) command() (*ast.Command, error) {

	// reads the literal value
	// p.scan()

	lit := &ast.Command{}
	lit.Token = p.tok
	return lit, nil
}

func (p *Parser) literal() (*ast.Literal, error) {

	lit := &ast.Literal{}
	lit.Token = p.tok
	return lit, nil
}

func (p *Parser) image() (*ast.Image, error) {

	p.cscope = scope.Img

	img := &ast.Image{}
	img.Start = p.tok.Position

	nodes, err := p.block()
	if err != nil {
		return nil, err
	}
	img.Node = nodes

	p.cscope = 0

	return img, nil
}

func (p *Parser) block() (*ast.Block, error) {

	if token.LBRACE == p.tok.Type {
		p.scan()
	}

	if p.cscope != scope.Img {
		p.cscope = scope.Blk
	}

	blk := &ast.Block{}
	blk.Lbrace = p.tok.Position

	nodes := make([]ast.Node, 0)
	for {

		n, err := p.nodeItem()
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, n)

		// next token
		p.scan()

		if token.RBRACE == p.tok.Type {
			break
		}
	}

	blk.Node = nodes
	blk.Rbrace = p.tok.Position

	if p.cscope == scope.Blk {
		p.f.BlockCount = p.f.BlockCount + 1
		blk.Number = p.f.BlockCount
	}

	if p.cscope != scope.Img {
		p.cscope = 0
	}

	return blk, nil
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
	return nil, fmt.Errorf("Expected: RPAREN ')' but %v", p.tok.Text)
}

func (p *Parser) hint() (*ast.Hint, error) {

	// after comment is scanned, the next token is uffered
	// set it to false, so that it can read the hint
	if p.prevTok.Type == token.COMMENT {
		p.tokenScanned = false
	}

	hint := &ast.Hint{}
	hint.Token = p.tok

	// reads the hint operation
	p.scan()

	if p.tok.Type == token.IDENTIFIER {

		valid := false
		for _, i := range validHints {
			if strings.EqualFold(i, p.tok.Text) {
				valid = true
				break
			}
		}

		if !valid {
			return nil, fmt.Errorf("Invalid Hint '%s', expecting %s", p.tok.Text, validHints)
		}

	} else {
		return nil, fmt.Errorf("Expected: IDENTIFIER but got: %s", p.tok.Type)
	}

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

func (p *Parser) list() (*ast.List, error) {

	l := &ast.List{}
	l.Lbrack = p.tok.Position

	nodes := make([]ast.Node, 0)
	for {

		p.scan()

		if token.LBRACK == p.tok.Type {
			continue
		}

		if token.RBRACK == p.tok.Type {
			l.RBrack = p.tok.Position
			break
		}

		if token.COMMA == p.tok.Type {
			continue
		}

		if token.STRING != p.tok.Type {
			return nil, fmt.Errorf("Expected a string type, but got %v", p.tok.Type)
		}

		lit := &ast.Literal{}
		lit.Token = p.tok

		nodes = append(nodes, lit)
	}

	l.Node = nodes

	return l, nil
}

func (p *Parser) validate() bool {

	return true
	// return nil, errors.New("Expected: STRING token (the name of the shiftfile, format :'company/name')")
}
