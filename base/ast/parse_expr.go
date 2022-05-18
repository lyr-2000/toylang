package ast

import (
	"toylang/base/lexer"
)

// type AstNode = Anode
type Node = Anode

type Tokens struct {
	tokens []*Token
	i      int
}

func parseExpr(t *Tokens) Node {
	return expr_parse_(t)
}
func ParseExpr_(t []*Token) Node {
	return expr_parse_(&Tokens{t, 0})
}
func expr_parse_(t *Tokens) Node {
	var e Expr
	node := e.parseE(t)
	return node
	//todo: solve it
}

func (e *Expr) parseE(t *Tokens) Node {
	var a = e.parseT(t)
	var b Node = nil
	for t.i < len(t.tokens) {
		if t.tokens[t.i].Value == "+" {
			c := newBinaryExpr(t.tokens[t.i])
			t.i++
			//consume
			b = e.parseT(t)
			c.Children = append(c.Children, a, b)
			a = c

		} else if t.tokens[t.i].Value == "-" {
			c := newBinaryExpr(t.tokens[t.i])
			t.i++
			b = e.parseT(t)
			c.Children = append(c.Children, a, b)
			a = c
		} else {
			break
		}
	}
	return a
}
func newBinaryExpr(lc *Token) *Expr {
	a := new(Expr)
	a.NodeType = BINARY_EXPR
	a.Lexeme = lc
	return a
}
func newUninaryExpr(lc *Token) *Expr {
	a := new(Expr)
	a.NodeType = UNARY_EXPR
	//例如 a++
	a.Lexeme = lc
	return a
}
func (e *Expr) parseT(t *Tokens) Node {
	var a = e.parseUnary_(t)
	var b Node = nil
	for t.i < len(t.tokens) {

		token := t.tokens[t.i]

		if t.tokens[t.i].Value == "*" {
			c := newBinaryExpr(token)
			t.i++
			b = e.parseUnary_(t)
			c.Children = append(c.Children, a, b)
			a = c
		} else if t.tokens[t.i].Value == "/" {
			c := newBinaryExpr(token)
			t.i++
			b = e.parseUnary_(t)
			c.Children = append(c.Children, a, b)
			a = c
		} else {
			break
		}
	}
	return a
}
func (e *Expr) parseUnary_(t *Tokens) Node {
	// 解决 a++ 的问题 , 然后 parseF
	if t.i >= len(t.tokens) {
		return nil
	}
	var anode = e.parseF(t)

	if t.i >= len(t.tokens) || anode == nil {
		return anode
	}
	//else
	// i< len
	token := t.tokens[t.i]

	if anode.GetLexeme().Type == lexer.Variable {
		// 只有变量 才能自增或者复制， 并且 目前设计是只能后置++
		if token.Value == "++" {
			// if anode.GetLexeme()
			t.i++
			top := newUninaryExpr(token)
			top.Children = append(top.Children, anode)
			anode = top
		} else if token.Value == "--" {
			t.i++
			top := newUninaryExpr(token)
			top.Children = append(top.Children, anode)
			anode = top
		} else if token.Value == "+=" {
			t.i++
			bnode := e.parseF(t)
			temp := newBinaryExpr(token)
			temp.Children = append(temp.Children, anode, bnode)
			anode = temp
		} else if token.Value == "-=" {
			t.i++
			bnode := e.parseF(t)
			temp := newBinaryExpr(token)
			temp.Children = append(temp.Children, anode, bnode)
			anode = temp

		}
	}

	return anode
}

func (e *Expr) parseF(t *Tokens) Node {
	if t.i >= len(t.tokens) {
		return nil
	}

	if t.tokens[t.i].Value == "(" {
		t.i++
		a := e.parseE(t)
		if t.tokens[t.i].Value != ")" {
			panic("illegal state of bracket match")
		}
		t.i++
		return a
	}

	token := t.tokens[t.i]
	// if token.Type == lexer.Operator {
	// 	// log.Printf("result %v", token)
	// 	// 判断 自增操作，例如 ++, -- 等
	// 	switch token.Value {
	// 	case "++", "--", "!":
	// 		t.i++ // eat token
	// 		var expr = new(Expr)
	// 		expr.NodeType = UNARY_EXPR
	// 		expr.Lexeme = token
	// 		expr.Children = append(expr.Children, e.parseE(t))
	// 		return expr
	// 	default:

	// 	}
	// }
	if token.Type == lexer.Number {
		// if t.tokens[t.i].Type != lexer.Number {
		// 	// panic(t.tokens[t.i])
		// }
		scalar := new(Scalar)
		scalar.Lexeme = t.tokens[t.i]
		// scalar.NodeType = SCALAR
		t.i++
		return scalar
	}
	if token.Type == lexer.Variable {
		variable := new(Variable)
		variable.Lexeme = token
		t.i++
		return variable
	}
	panic("illegal node at parse")
	// scalar := new(Scalar)
	// scalar.Lexeme = t.tokens[t.i]
	// // scalar.NodeType = SCALAR
	// t.i++

}
