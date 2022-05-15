
## 词法分析器实现记录


### 参考项目

[lexer参考项目](https://github.com/wupeaking/panda)


[hacklang项目视频介绍](https://www.bilibili.com/video/BV1eS4y1G7oC?spm_id_from=333.337.search-card.all.click)

[hacklang代码参考](https://github.com/4ra1n/HacLang)




## 实现读取 变量

```go
func (l *Lexer) readKeywordOrVariableKey_() *Token {

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
```


## 读取字符串

目前字符串 模仿 js那种模式

```go

func (l *Lexer) readString_() *Token {
	var buf bytes.Buffer
	var state uint8 = 0
	//state == 1 => "
	//state == 2 ==> '
	//state == 3 ==> "xx"
	//state == 4 ==> 'xxx'
	//if l.Peek() != '"' {
	//	panic("illegal state  string")
	//}
	for l.HasNext() {
		c := l.Next()
		switch state {
		case 0:
			if c == '"' {
				state = 1
			} else {
				state = 2
			}
			buf.WriteRune(c)
		case 1:
			if c == '"' { //state = 4
				//end
				buf.WriteRune(c)
				return NewToken(String, buf.String(), l.Pos().Line, l.Pos().Offset)
			} else {
				buf.WriteRune(c)
			}
		case 2:
			if c == '\'' { //state = 3
				buf.WriteRune(c)
				return NewToken(String, buf.String(), l.Pos().Line, l.Pos().Offset)
			} else {
				buf.WriteRune(c)
			}

		}
	}
	//没有迭代完成
	panic("illegal state of string")
	return nil
}


```


## 读取 操作符



## 踩坑

golang不支持继承，只支持组合，所以 没办法 对 struct的对象定义的方法进行重写
原先是这样写的：

```go
type BaseLexer struct {
	*scanner.Scanner
}
type LexerWithCache struct {
    *BaseLexer
 
    que *list.Queue
 
}


func (l *LexerWithCache) Peek() Char_ {
    if l.que.ListCnt > 0 {
        return list.QueuePeek(l.que).(char)
    }
    return l.Scanner.Peek()
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

```





上面有一个坑是 ， text/scanner 里面也有个 Peek 方法， 我 在  LexerWithCache 重写了 Peek方法，但是 readKeywordOrVariableKey_() 是 BaseLexer 里面定义的，所以 无法重写这个函数，还是调用的 text/scanner 里面的方法，导致我调试了很久，解决的方法如下



```go
type BaseLexer = LexerWithCache
type LexerWithCache struct {
	//*BaseLexer
	*scanner.Scanner
	que *list.Queue
	//st  *list.Stack
}

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
```



这里没必要用 继承，直接写在一起就可以了



## 参考 otto的实现




[lexer参考链接](https://github.dev/robertkrimen/otto) 

参考 otto/lerxer.go 文件


```go


func digitValue(chr rune) int {
	switch {
	case '0' <= chr && chr <= '9':
		return int(chr - '0')
	case 'a' <= chr && chr <= 'f':
		return int(chr - 'a' + 10)
	case 'A' <= chr && chr <= 'F':
		return int(chr - 'A' + 10)
	}
	return 16 // Larger than any legal digit value
}

func isDigit(chr rune, base int) bool {
	return digitValue(chr) < base
}

// 下面是函数调用

if self.chr == '0' {
		offset := self.chrOffset
		self.read()
		if self.chr == 'x' || self.chr == 'X' {
			// Hexadecimal
			self.read()
			if isDigit(self.chr, 16) {
				self.read()
			} else {
				return token.ILLEGAL, self.str[offset:self.chrOffset]
			}
			self.scanMantissa(16)

			if self.chrOffset-offset <= 2 {
				// Only "0x" or "0X"
				self.error(0, "Illegal hexadecimal number")
			}

			goto hexadecimal
		} else if self.chr == '.' {
			// Float
			goto float
		} else {
			// Octal, Float
			if self.chr == 'e' || self.chr == 'E' {
				goto exponent
			}
			self.scanMantissa(8)
			if self.chr == '8' || self.chr == '9' {
				return token.ILLEGAL, self.str[offset:self.chrOffset]
			}
			goto octal
		}
	}

```

