package ast

import (
	"fmt"
	"toylang/base/lexer"
)

func parseFnParameter(t *Tokens) Node {
	var exprNode []Anode
	if t.peek() == "(" {
		t.next()
		for t.peek() != "" && t.peek() != ")" {
			if t.peek() == "," {
				t.next()
				continue
			}
			exprNode = append(exprNode, parseExpr(t))
		}
		t.next()

	} else {
		for t.peek() != ";" && t.peek() != "" {

			if t.peek() == "," {
				t.next()
				continue
			}

			// endl
			exprNode = append(exprNode, parseExpr(t))
		}
		t.next()
	}
	// panic("cannot parse fn params")
	var result = new(FnParam)
	result.Children = exprNode
	return result
}

func parseFn(t *Tokens) Node {
	if !t.hasNext() {
		return nil
	}

	p := t.peek()
	if p == "fn" || p == "function" {
		t.i++
	}
	//function app() { a=11111; }
	// read fnName
	if t.hasNext() == false {
		return nil
	}
	if t.tokens[t.i].Type != lexer.Variable {
		//variable
		panic(fmt.Sprintf("is not an function name at match fn , actual %v", t.tokens[t.i]))
	}

	var parent = new(FuncStmt)
	parent.Lexeme = t.tokens[t.i]
	t.i++

	var child = parseFnParameter(t)
	if t.peek() != "{" {
		panic(fmt.Sprintf("cannot read body actual %v", t.peek()))
	}
	var block Anode = parseBlock(t)
	parent.Children = []Anode{child, block}

	return parent
}
