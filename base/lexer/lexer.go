package lexer

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

//	BaseLexer type BaseLexer struct {
//		*scanner.Scanner
//	}
type BaseLexer = LexerWithCache
type Lexer = BaseLexer

func New(reader io.Reader) *LexerWithCache {
	//lexer := &BaseLexer{
	//	&scanner.Scanner{},
	//}
	l := newLexer()
	//init reader
	//l.BaseLexer.Scanner = &scanner.Scanner{}
	l.Scanner.Init(reader)
	return l
}
func (l *BaseLexer) HasNext() bool {
	return l.que.ListCnt > 0 || l.Peek() != -1
}

//func (l *BaseLexer) next() char {
//	return l.Next()
//}

func (l *BaseLexer) skipBlankChars_() { //需要判断 是不是空格 后才能调用，这个方法不会自己判断
FindSpace:
	for {
		ch := l.Peek()
		switch {
		case IsSpace(ch):
			l.Next()
		default:
			break FindSpace
		}
	}
}

// //bug: bugcheck
// func (l *BaseLexer) skipCommentsLines_() { //需要判断 是不是注释 后才能调用，这个方法不会自己判断
// 	if l.Peek() != '/' {
// 		return
// 	}
// deleteComments: //删除注释
// 	for {

// 		switch l.Peek() {
// 		case scanner.EOF:
// 			break deleteComments
// 		case '\n':
// 			l.Next() // eat \n
// 			if l.Peek() == '\r' {
// 				l.Next()
// 			}
// 			break deleteComments //注释换行了，就退出
// 		default:
// 			l.Next()
// 			if l.Peek() == '\r' {
// 				l.Next()
// 			}
// 			break
// 		}

//		}
//	}
func isSingleQuote(ch char) bool {
	return ch == '\''
}
func isDoubleQuote(ch char) bool {
	return ch == '"'
}
func (l *BaseLexer) readChar_() *Token {
	ch := l.Next()
	if !isSingleQuote(ch) {
		panic("not a char token")
	}
	var buf char
	var cnt uint16
	for !isSingleQuote(l.Peek()) {
		cnt++
		buf = l.Peek()
		if buf == '\\' {
			cnt--
			continue
		}
		l.Next()
	}
	if cnt != 1 {
		panic("not a char token")
	}
	return makeToken(Char, buf)

}

// func (l *BaseLexer) readBigFloat() *Token {
// 	for {
// 		ch := l.Peek()
// 	}
// 	// example: 1e9+7
// 	// 1e9
// 	// 1e+9
// 	// 1e-9
// 	// 1e9+7
// 	// 1e9-7
// 	// 1e9*7
// 	return nil
// }

func (l *BaseLexer) readNumberTok() *Token {
	ch := l.Next()
	switch {
	case ch == '0':
		if l.Peek() != '.' {
			// not  0.xxx, return 0
			pos := l.Pos()
			return NewToken(Number, "0", pos.Line, pos.Offset) //转 int
			//return Token{Type: Number, Value: 0, Line: pos.Line, Position: pos.Offset}
		}
		//eat .
		l.Next()
		var buf = bytes.Buffer{}
		buf.WriteString("0.") // 0. xxx
	Label:
		for {
			peekchar := l.Peek()
			switch {
			case 'e' == peekchar || 'E' == peekchar:
				l.Next()
				buf.WriteRune(ch)
				continue Label
			case IsNumber(peekchar):
				l.Next() //for->next
				buf.WriteRune(peekchar)
			default:
				break Label
			}
		}

		if buf.Len() < 3 {
			return NewTokenPos(Illegal, buf.String(), l.Pos())
		}
		return NewTokenPos(Number, buf.String(), l.Pos())
	case IsNumber(ch):
		var buf bytes.Buffer
		buf.WriteRune(ch)
		hasDot := false
		eNumber := false
		specialOpChar := 0
	readNumber:
		for {
			fc := l.Peek()
			if fc == '.' && !hasDot {
				//first meet dot
				hasDot = true
				l.Next()
				buf.WriteRune(fc)
				continue readNumber
			}
			if IsNumber(fc) {
				l.Next()
				buf.WriteRune(fc)
				continue readNumber
			}
			if fc == 'e' || fc == 'E' {
				eNumber = true
				l.Next()
				buf.WriteRune(fc)
				continue readNumber
			}
			if eNumber && specialOpChar == 0 && fc == '+' || fc == '-' {
				l.Next()
				specialOpChar++
				buf.WriteRune(fc)
				continue readNumber
			}
			break readNumber // not a number
		}
		return NewToken(Number, buf.String(), l.Pos().Line, l.Pos().Offset)
	default:

	}
	return NewToken(Illegal, nil, l.Pos().Line, l.Pos().Offset)
}

func (l *BaseLexer) readKeywordOrVariableKey_() *Token {

	var buf = strings.Builder{} //""
	//buf.WriteRune(ch)
	for l.HasNext() {
		readchar := l.Peek()
		if IsLiteral(readchar) {
			buf.WriteRune(readchar) //append char to keyword
		} else {
			break
		}
		l.Next()
	}
	literal := buf.String()
	if IsKeyword(literal) {
		//系统关键字，var,fun, 等
		return NewToken(Keyword, literal, l.Pos().Line, l.Pos().Offset)
	}
	switch literal {
	//case "null":
	case "true", "false":
		//零值
		return NewToken(Boolean, literal, l.Pos().Line, l.Pos().Offset)
	default:
		return NewToken(Variable, literal, l.Pos().Line, l.Pos().Offset)
	}

}

func (l *BaseLexer) readString_() *Token {
	var buf bytes.Buffer
	//var state uint8 = 0
	var L = l.Next()
	if L != '"' {
		panic("illegal state literal")
	}

	for l.HasNext() && l.Peek() != L {
		c := l.Next()
		if c == '\\' { //转义字符串
			// readNext
			c = l.Next()
			buf.WriteRune(c)
		} else {
			buf.WriteRune(c)
		}
	}
	//it must be l.peek() == L
	l.Next()
	//没有迭代完成
	return makeToken(String, buf.String())
}

func mixOperator_(l, r char) string {
	if r == -1 {
		return string(l)
	}
	if l == '+' {
		if r == '+' {
			return "++"
		} else if r == '=' {
			return "+="
		}
		return "+"
	} else if l == '-' {
		if r == '-' {
			return "--"
		}
		if r == '=' {
			return "-="
		}
		if r == '>' {
			return "->"
		}
		return "-"
	} else if l == '*' {
		if r == '=' {
			return "*="
		}

		return "*"
	} else if l == '/' {
		if r == '=' {
			return "/="
		}
		return "/"
	} else if l == '%' {
		if r == '=' {
			return "%="
		}
		return "%"
	} else if l == '&' {
		if r == '&' {
			return "&&"
		} else if r == '=' {
			return "&="
		}
		return "&"

	} else if l == '|' {
		if r == '=' {
			return "|="
		}
		if r == '|' {
			return "||"
		}
		return "|"
	} else if l == '!' {
		if r == '=' {

			return "!="
		}
		return "!"
	} else if l == '^' {
		if r == '=' {
			return "^="
		}
		return "^"
	} else if l == ',' {
		return ","
	} else if l == ';' {
		return ";"
	} else if l == '<' {
		if r == '<' {
			return "<<"
		} else if r == '=' {
			return "<="
		}
		return "<"
	} else if l == '>' {
		if r == '>' {
			return ">>"
		} else if r == '=' {
			return ">="
		}
		return ">"
	} else if l == '=' {
		if r == '>' {
			return "=>"
		} else if r == '=' {
			return "=="
		}
		return "="
	}
	//panic(fmt.Sprintf("illegal state operator %v,%v", string(l), string(r)))
	//return ""
	return string(l)
}
func (l *BaseLexer) readOperator_() *Token {
	//state := 0
	var L = l.Peek()
	if !IsOperator(L) {
		panic(fmt.Sprintf("illegal state %v %v", string(L), l.Pos()))
	}
	l.Next()
	var R = l.Peek()
	if IsOperator(R) { //mix operator
		var res = mixOperator_(L, R)
		//l.Next() // eat r
		if len(res) == 2 {
			l.Next() //eat R
		}
		return makeToken(Operator, res)
	}

	return makeToken(Operator, string(L))
}
