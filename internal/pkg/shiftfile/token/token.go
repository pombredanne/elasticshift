/*
Copyright 2018 The Elasticshift Authors.
*/
package token

import (
	"fmt"
	"strconv"
)

var eof = rune(0)

type Type int

// Token structure
type Token struct {
	Type     Type
	Text     string
	Position Position
}

//  Here is the list of ten types
const (
	ILLEGAL Type = iota
	EOF
	COMMENT // # this is the comment tag
	NEWLINE // \n
	WHITESPACE

	literal_beg
	IDENTIFIER // literals
	INT        //123
	STRING     // "abc"
	BOOL       // true | false
	FLOAT      // 1.23

	literal_end

	operator_beg
	ASSIGN // =

	LPAREN // (
	LBRACK // [
	LBRACE // {
	COMMA  // ,
	PERIOD // .

	RPAREN   // )
	RBRACK   // ]
	RBRACE   // }
	ENV      // $
	HINT     // //
	LHINT    // /*
	RHINT    // */
	HINT_DEL // :
	ARGUMENT // @
	SECRET   // ^
	operator_end

	keyword_beg
	KEYWORD
	FROM
	IMAGE
	NAME
	VAR
	VERSION
	WORKDIR
	SCRIPT
	LANGUAGE
	COMMAND
	CACHE
	DIRECTORY
	keyword_end
)

// Token keyword
var tokens = [...]string{

	ILLEGAL:    "ILLEGAL",
	EOF:        "EOF",
	COMMENT:    "COMMENT",
	NEWLINE:    "NEWLINE",
	WHITESPACE: "WHITESPACE",

	IDENTIFIER: "IDENTIFIER",
	INT:        "INT",
	STRING:     "STRING",
	BOOL:       "BOOL",
	FLOAT:      "FLOAT",

	ASSIGN: "ASSIGN",

	LPAREN: "LPAREN",
	LBRACK: "LBRACK",
	LBRACE: "LBRACE",
	COMMA:  "COMMA",
	PERIOD: "PERIOD",

	RPAREN:   "RPAREN",
	RBRACK:   "RBRACK",
	RBRACE:   "RBRACE",
	ENV:      "ENV",
	HINT:     "HINT",
	LHINT:    "LHINT",
	RHINT:    "RHINT",
	HINT_DEL: "HINT_DEL",
	ARGUMENT: "ARGUMENT",
	SECRET:   "SECRET",

	FROM:      "FROM",
	IMAGE:     "IMAGE",
	NAME:      "NAME",
	VAR:       "VAR",
	VERSION:   "VERSION",
	WORKDIR:   "WORKDIR",
	SCRIPT:    "SCRIPT",
	LANGUAGE:  "LANGUAGE",
	COMMAND:   "COMMAND",
	CACHE:     "CACHE",
	DIRECTORY: "DIRECTORY",
}

var keywords map[string]Type

func init() {
	keywords = make(map[string]Type)
	for i := keyword_beg + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}
}

func Lookup(iden string) Type {
	if tok, isKeyword := keywords[iden]; isKeyword {
		return tok
	}
	return IDENTIFIER
}

// Returns the string representation of the given ten
func (t Type) String() string {

	s := ""
	if 0 <= t && t < Type(len(tokens)) {
		s = tokens[t]
	}
	if s == "" {
		s = "ten(" + strconv.Itoa(int(t)) + ")"
	}
	return s
}

// Returns the value of the token
func (t Token) Value() interface{} {

	switch t.Type {
	case INT:
		value, err := strconv.ParseInt(t.Text, 0, 64)
		if err != nil {
			panic(err)
		}
		return int64(value)
	case BOOL:
		if t.Text == "true" {
			return true
		} else if t.Text == "false" {
			return false
		}
		panic(fmt.Sprintf("Unknown bool type %s", t.Text))
	case STRING:
		if t.Text == "" {
			return ""
		}
		return t.Text
	case FLOAT:
		value, err := strconv.ParseFloat(t.Text, 64)
		if err != nil {
			panic(err)
		}
		return float64(value)
	case IDENTIFIER:
		return t.Text
	default:
		panic(fmt.Sprintf("Unknown type %s", t.Type))
	}
	return nil
}

// IsLiteral returns true for tens corresponding to identifiers
// and basic type literals; it returns false otherwise.
//
func (t Type) IsLiteral() bool { return literal_beg < t && t < literal_end }

// IsOperator returns true for tens corresponding to operators and
// delimiters; it returns false otherwise.
//
func (t Type) IsOperator() bool { return operator_beg < t && t < operator_end }

// IsKeyword returns true for tokens corresponding to keywords;
// it returns false otherwise.
//
func (t Type) IsKeyword() bool { return keyword_beg < t && t < keyword_end }
