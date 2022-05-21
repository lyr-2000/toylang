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
	/*

	   {
	   			"=", "+=", "-=", "*=", "/=", "%=", "&=", "^=", "!=",
	   		},
	   		{
	   			"||",
	   		},
	   		{
	   			"&&",
	   		},
	   		{
	   			"|",
	   		}, {
	   			"^",
	   		}, {
	   			"&",
	   		}, {
	   			"==", "!=",
	   		}, {
	   			">=", "<=", ">", "<",
	   		}, {
	   			">>", "<<",
	   		}, {
	   			"+", "-",
	   		}, {
	   			"*", "/", "%",
	   		},


	*/

	switch c {
	case '+', '-', '*', '/', '%',
		'\\',
		// '!',
		'&', '|',
		'^',
		'=', '>', '<', '!',
		',', ';',
		'.':
		return true
	default:

	}
	return false
}

// func stringToBytes(s *string) []byte {
// 	return *(*[]byte)(unsafe.Pointer(s))
// }

var (
	// keywords = map[string]bool{
	// 	"bool":     true,
	// 	"var":      true,
	// 	"if":       true,
	// 	"else":     true,
	// 	"const":    true,
	// 	"def":      true,
	// 	"function": true,
	// 	"fn":       true,
	// }
	keywords = []string{
		"bool",
		"var",
		"global", "lookup", // 全局变量使用
		"if",
		"for", "while", "break",
		"else",
		"const",
		"def",
		"fn",
		"function",
		"int",
		"float",
		"double",
		"long",
		"number",
		"object",
		"void",
		"return",
	}
)

func IsKeyword(s string) bool {
	for i, _ := range keywords {
		if s == keywords[i] {
			return true
		}
	}
	return false
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
