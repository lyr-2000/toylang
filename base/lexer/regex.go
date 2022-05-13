package lexer

import (
	"unicode"
)

/*
import "regexp"

var (
	PatternLetter   = regexp.MustCompile("^[a-zA-Z]$")
	PatternNumber   = regexp.MustCompile("^[0-9]$")
	PatternLiteral  = regexp.MustCompile("^[_a-zA-Z0-9]$")
	PatternOperator = regexp.MustCompile("^[+-\\\\*<>=!&|^%/,]$")
)
*/
type char = rune
type Char_ = char

func IsSpace(c char) bool {

	switch c {
	case '\n', '\t', '\r', ' ':
		return true
	case '\u000b', '\u000c', '\u00a0', '\ufeff':
		return true
	case '\u2028', '\u2029':
		return false
	case '\u0085':
		return false
	//case ' ', '\n', '\r', '\t':
	//	return true
	default:
		return unicode.IsSpace(c)
	}

}
func IsBracket(c char) bool {
	return c == '{' ||
		c == '}' ||
		c == '(' ||
		c == ')' ||
		c == '[' ||
		c == ']'
}
func IsLetter(c char) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z')
}
func IsNumber(c char) bool {
	return c >= '0' && c <= '9'
}

//是否是 字符串
func IsLiteral(c char) bool {
	return c == '_' || IsLetter(c) || IsNumber(c)
}

//是否是操作符
func IsOperator(c char) bool {
	//^[+-\\*<>=!&|^%/]$
	//for _, f := range opcodecharset {
	//	if f == c {
	//		return true
	//	}
	//}
	switch c {
	case '+',
		'-',
		'*',
		'/',
		'\\',
		'!',
		'&',
		'^',
		'%',
		'=',
		',', ';':
		return true
	default:

	}
	return false
}

// func stringToBytes(s *string) []byte {
// 	return *(*[]byte)(unsafe.Pointer(s))
// }

//TODO: 实现高性能判断 keyword
func IsKeyword(s string) bool {
	return s == "var"
}
func HasKeyword(s string) bool {
	return IsKeyword(s)
}

//func IsKeywordBuffer(buf *bytes.Buffer) bool{
//	keyword := buf.String()
//	if keyword == "true" ||
//		keyword == "false" {
//		return
//	}
//}
//func IsKeyword(s *string) bool {
//
//	var value = *(*[]byte)(unsafe.Pointer(&s))
//	if bytes.Equal(value, []byte("true")) || bytes.Equal(value, []byte("false")) {
//		// is boolean
//		return NewToken()
//	}
//
//}