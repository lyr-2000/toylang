package lexer

import (
	"bytes"
	"fmt"
	"text/scanner"
	"toylang/base/list"
)

/*
type BaseLexer struct {
	*scanner.Scanner
}

*/
type LexerWithCache struct {
	//*BaseLexer
	*scanner.Scanner
	que *list.Queue
	//st  *list.Stack
}

func (l *LexerWithCache) String() string {
	return l.que.String()
}
func newLexer() *LexerWithCache {
	l := new(LexerWithCache)
	l.que = &list.Queue{}
	l.Scanner = &scanner.Scanner{}
	//l.BaseLexer = new(BaseLexer)
	return l
}
func NewStringLexer(s string) *LexerWithCache {
	return New(bytes.NewBufferString(s))
}
func (l *LexerWithCache) LastPosString() string {
	x := l.Peek()
	if x == -1 {
		return "EOF"
	}
	buf := bytes.Buffer{}
	//buf.WriteRune(x)
	for buf.Len() < 8 && l.HasNext() {
		buf.WriteRune(l.Next())
	}
	return buf.String()
}

func (l *LexerWithCache) PutBackChar(cs ...char) {
	for _, c := range cs {
		l.que = list.QueueAppend(l.que, c)
	}
}
func (l *LexerWithCache) PutBackStr(s string) {
	for _, v := range s {
		l.PutBackChar(v)
	}
}
func (l *LexerWithCache) ClearQueue() {
	list.QueueClear(l.que) //clear queue
}
func (l *LexerWithCache) Peek() Char_ {
	if l.que.ListCnt > 0 {
		return list.QueuePeek(l.que).(char)
	}
	return l.Scanner.Peek()
}
func (l *LexerWithCache) Next() Char_ {
	if l.que.ListCnt > 0 {
		poll := list.QueuePoll(l.que)
		return poll.(char)
	}
	// que empty]
	ch := l.Scanner.Next()

	return ch
}
func (l *LexerWithCache) readBrackets_() *Token {
	c := l.Next()
	return makeToken(Brackets, string(c))

}

func (l *LexerWithCache) peekAndComments_() bool {
	c := l.Peek()
	if c != '/' {
		panic("illegal state of read comments line")
	}
	l.Next()
	h := l.Peek()
	if h == '/' {
		// is //
		for l.HasNext() && l.Next() != '\n' {
		} // delete comment line

	} else if h == '*' {
		// is  /*
		valid := false
		for l.HasNext() {
			tchar := l.Next()
			if tchar == '*' && l.Peek() == '/' {
				valid = true
				l.Next() // delete multi line comment /**/
				break
			}
		}
		if !valid {
			panic("illegal comment matched")
		}
	} else {
		// 不是 // 或者 /* , 则是其他字符，就返回去，重新读取
		l.PutBackChar(c)
		l.PutBackChar(h)
		return false
	}
	return true
}

//读取 token
func (l *LexerWithCache) ReadToken() []*Token {

	var result []*Token
	//defer func() {
	//	err := recover()
	//	fmt.Printf("last char is %d %c ; %v ,%v", int(l.Peek()), l.Peek(), err, result)
	//}()
	for l.HasNext() {
		//l.ClearQueue() // clear token cache
		c := l.Peek()
		if c == -1 {
			break
		}
		if c == '/' && l.peekAndComments_() {
			continue
		}
		if IsSpace(c) {
			l.skipBlankChars_()
			//c = l.Peek()
			continue
		}
		if IsBracket(c) {
			//匹配到括号字符
			result = append(result, l.readBrackets_())
			continue
		}

		if isDoubleQuote(c) { //read string
			result = append(result, l.readString_())
		} else if IsLiteral(c) {
			result = append(result, l.readKeywordOrVariableKey_())
		} else if isSingleQuote(c) {
			//read char
			result = append(result, l.readChar_())
		} else if IsNumber(c) {
			result = append(result, l.readNumber_())
		} else if c == '+' || c == '-' {
			l.Next() //for c next
			spacef := IsSpace(l.Peek())
			if spacef {
				l.skipBlankChars_()
			}
			peekchar := l.Peek()
			if IsNumber(peekchar) {
				//+12,-12
				tk := l.readNumber_()
				if c == '-' {
					//负号就修改内容
					tk.Value = fmt.Sprintf("-%v", tk.Value)
				}
				result = append(result, tk)
				continue
			} else {

				// operator符号，就返回去
				l.PutBackChar(c)
				if spacef {
					l.PutBackChar(' ')
				}
			}
		}
		c = l.Peek()

		if IsOperator(c) {
			//fmt.Println(string(c))
			result = append(result, l.readOperator_())
		}

	}
	return result
}
