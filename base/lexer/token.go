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
type TokenType int16

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

type Token struct {
	Type TokenType
	//Value    string //literal
	Value interface{}
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
	return fmt.Sprintf("{type=%v,value=%v}", t.Type.String(), t.Value)
}

//func (t *Token) IsType() {
//	return t.Value == "bool"
//}
