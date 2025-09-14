package lexer

import (
	"fmt"
	"text/scanner"
)

/*
- keyword
- variable
- operator
- brackets
- string
- number /float64
- boolean
*/

//go:generate stringer -type TokenType
type TokenType int8

const (
	_        TokenType = iota
	Keyword            //var ,for ,goto,
	Variable           // var a
	Operator           // +-*/
	Brackets           // () {} []
	String             // "a"
	Char               // 'a'
	Number             //1.34
	Boolean            // true false
	EOF                // eof

	Illegal //illegal state
)
func (t TokenType) String() string {
	switch t {
	case Keyword:
		return "keyword"
	case Variable:
		return "variable"
	case Operator:
		return "operator"
	case Brackets:
		return "brackets"
	case String:
		return "string"
	case Char:
		return "char"
	case Number:
		return "number"
	case Boolean:
		return "boolean"
	case EOF:
		return "eof"
	case Illegal:
		return "illegal"
	default:
		return "unknown"
	}
	return ""
}

type TokenValue = interface{}

type Token struct {
	Type TokenType
	//Value    string //literal
	Value TokenValue
	//Line     int, 不需要 line 和 col ,浪费内存
	//Position int
}

func NewTokenPos(t TokenType, v interface{}, _ scanner.Position) *Token {
	return &Token{t, v}
}
func makeToken(t TokenType, v interface{}) *Token {
	return &Token{t, v}
}

func NewToken(t TokenType, v interface{}, _, _ int) *Token {
	return &Token{t, v}
}

func (t *Token) String() string { // convert Token to string
	if t == nil {
		return "<NIL>"
	}
	return fmt.Sprintf("{type:%v,value:%v}", t.Type.String(), t.Value)
}

//func (t *Token) IsType() {
//	return t.Value == "bool"
//}
