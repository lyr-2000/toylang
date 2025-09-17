package ast

import (
	"github.com/lyr-2000/toylang/base/lexer"
	"strings"
)

type Stmt struct {
}

func (t *Tokens) peek() string {
	if t.i >= len(t.tokens) {
		return ""
	}
	return t.tokens[t.i].Value.(string)
}
func (t *Tokens) next() *Token {
	if t.i >= len(t.tokens) {
		return nil
	}
	cur := t.tokens[t.i]
	t.i++
	return cur
}
func (t *Tokens) back(i int) {
	t.i = i
}
func (t *Tokens) hasNext() bool {
	return t.i < len(t.tokens)
}
func ParseStmt(t *Tokens) Anode {
	bn := new(BlockNode)
	for t.hasNext() {
		bn.Children = append(bn.Children, parseStmt(t))
	}
	return bn
}

var (
	keywordMap = map[string]struct{}{
		"break":      {},
		"continue":   {},
		"debugger":   {},
		"printstack": {},
		// "fatal":{},
	}
)

func parseStmt(t *Tokens) Anode {
	for t.hasNext() {
		if t.hasNext() && t.current().Value == ";" {
			t.i++
		} else {
			break // quit
		}
	}

	if !t.hasNext() {
		return nil
	}
	if t.peek() == "}" {
		return nil
	}

	var prev = t.i
	var token = t.next()
	var head = t.peek()
	t.back(prev)
	if _, ok := keywordMap[token.Value.(string)]; ok {
		cont := new(KeywordStmt)
		cont.Lexeme = token
		t.next()
		return cont
	}
	if strings.EqualFold(token.Value.(string), "for") || strings.EqualFold(token.Value.(string), "while") {
		return parseForStmt(t)
	}
	if strings.EqualFold(token.Value.(string), "if") {
		return parseIf(t)
	}
	if strings.EqualFold(token.Value.(string), "return") {
		return parseReturnStmt(t)
	}
	if token.Value == "{" {
		return parseBlock(t)
	}
	if token.Value == "." {
		// call func
		return parseCallFuncStmt(t)
	}
	//a=1
	if token.Type == lexer.Variable {
		//case: a=1, or  b = 2
		if head == "=" {
			return parseExpr(t)
		}
	} else if token.Value == "var" {
		//var a = 1
		//define
		return parseDeclareStmt(t)
	} else if token.Value == "fn" {
		return parseFn(t)
	}

	return parseExpr(t)
}
