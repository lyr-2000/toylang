package lexer

import (
	"encoding/base64"
	"fmt"
	"text/scanner"
	"unicode"
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

func (r *Token) VarName() string {
	if r == nil {
		return "@error!nil"
	}
	if r.Type == Variable {
		return r.Value.(string)
	}
	return "@error!"
}
func (r *Token) ToString() string {
	if r == nil {
		return "@eof"
	}
	if r.Type == String {
		if onlyLetter(r.Value.(string)) {
			return fmt.Sprintf("%v %v", r.Type.String(), r.Value)
		}
		buf := base64.StdEncoding.EncodeToString([]byte(r.Value.(string)))
		return fmt.Sprintf("%v %v", "strb64", buf)
	}
	return fmt.Sprintf("%v %v", r.Type.String(), r.Value)
}

func onlyLetter(s string) bool {
	for _, v := range s {
		if unicode.IsSpace(v) {
			return false
		}
		if unicode.IsNumber(v) {
			continue
		}
		if !unicode.IsLetter(v) {
			return false
		}
	}
	return true
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
