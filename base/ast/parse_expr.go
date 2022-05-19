package ast

import (
	"fmt"
	"toylang/base/lexer"
)

// type AstNode = Anode
type Node = Anode

type Tokens struct {
	tokens []*Token
	i      int
}

func NewTokens(ts []*lexer.Token) *Tokens {
	return &Tokens{
		i:      0,
		tokens: ts,
	}
}
func (t *Tokens) peekNext(id int) string {
	if id+t.i >= len(t.tokens) {
		return ""
	}
	return t.tokens[t.i+id].Value.(string)
}
func (t *Tokens) current() *Token {
	if !t.hasNext() {
		return nil
	}
	return t.tokens[t.i]
}
func (t *Tokens) NextMust(idx int, must string) {
	if t.i+idx >= len(t.tokens) {
		panic("illegal match")
	}
	if t.tokens[t.i+idx].Value != must {
		panic(fmt.Sprintf("cannot match %v", must))
	}

}
func parseExpr(t *Tokens) Node {
	return expr_parse_(t)
}
func ParseExpr_(t []*Token) Node {
	return expr_parse_(&Tokens{t, 0})
}
func expr_parse_(t *Tokens) Node {
	var e Expr
	node := e.parseE(t, 0)
	// log.Printf("read info %+v, %v\n", t.i, len(t.tokens))
	return node
}

var (
	tables = [][]string{
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
	}
)

//1.处理伪运算符，  -> , \\.  , [] 等 ，优先级最高
//2. 处理单目运算符， ++ ,  次高

func (e *Expr) parseE(t *Tokens, k int) Node {
	if k >= len(tables) {

		return e.parseTop2(t)
	}
	var a = e.parseE(t, k+1)
	for t.i < len(t.tokens) {
		cur := t.tokens[t.i]
		if k == 0 {
			if cur.Value == "{" || cur.Value == "," {
				break
			}
			if cur == nil || cur.Value == ";" { //end line
				// t.i++
				break
			}
		}
		match := false
		for _, v := range tables[k] {
			if v == cur.Value {
				//match it
				t.i++
				c := newBinaryExpr(cur)
				b := e.parseE(t, k+1)
				c.Children = []Anode{a, b}
				a = c
				match = true
			}

		}
		if !match {
			break
		}
	}

	return a
}
func (e *Expr) parseTop1(t *Tokens) Node {
	token := t.tokens[t.i]
	if token.Value == "(" {
		t.i++
		child := e.parseE(t, 0)
		if child == nil {
			//fail
			return child
		}
		if t.tokens[t.i].Value != ")" {
			panic(fmt.Sprintf("illegal match at ), actual %v", t.tokens[t.i]))
		}
		t.i++
		return child
	} else if token.Value == "!" || token.Value == "-" || token.Value == "+" {
		t.i++
		child := e.parseE(t, 0)
		if child == nil {
			return child
		}
		expr := newUninaryExpr(token)
		expr.Children = []Anode{
			child,
		}
		return expr
	}
	if token.Value == "." {
		//call func
		return parseCallFuncStmt(t)
	}

	if token.Type == lexer.Number || token.Type == lexer.String ||
		token.Type == lexer.Char ||
		token.Type == lexer.Boolean {
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

	if token.Value == ";" || token.Value == "\n" {
		t.i++ //skip end line
		return nil
	}
	panic("cannot explain ast tree")
	// return nil
}
func (e *Expr) parseTop2(t *Tokens) Node {
	if t.i >= len(t.tokens) {
		return nil
	}

	a := e.parseTop1(t)
	if a == nil || !t.hasNext() {
		return a
	}
	token := a.GetLexeme()
	if token == nil {
		//!warning: function call param node
		return nil
	}
	if token.Type != lexer.Variable {
		// is variable
		return a
	}
	//is variable
	if t.tokens[t.i].Value == "[" {
		t.i++
		bnode := e.parseE(t, 0)
		if t.tokens[t.i].Value != "]" {
			panic("illegal match brackets [] at right")
		}
		t.i++ //eat ]
		newParent := new(MapIndexNode)
		newParent.Children = []Anode{
			a, bnode,
		} // a  [1]
		//index node
		return newParent

	}
	if t.tokens[t.i].Value == "." {
		t.i++
		bnode := e.parseE(t, 0)
		// if t.tokens[t.i].Value != "]" {
		// 	panic("illegal match brackets [] at right")
		// }
		if bnode == nil {
			panic("illegal syntax")
		}
		if bnode.GetLexeme().Type != lexer.Variable {
			panic("parse at . error")
		}
		newParent := new(MapIndexNode)
		newParent.Children = []Anode{
			a, bnode,
		} // a  [1]
		//index node
		return newParent
	}
	cur := t.tokens[t.i]
	if cur.Value == "++" || cur.Value == "--" {
		t.i++
		// bnode := e.parseE(t, 0)
		newParent := newUninaryExpr(cur)
		newParent.Children = []Anode{a}
		return newParent
	}
	return a
}
func newBinaryExpr(tk *Token) *Expr {
	var ptr = new(Expr)
	ptr.Lexeme = tk
	ptr.NodeType = BINARY_EXPR
	return ptr
}
func newUninaryExpr(tk *Token) *Expr {
	var ptr = new(Expr)
	ptr.Lexeme = tk
	ptr.NodeType = BINARY_EXPR
	return ptr
}

// func (e *Expr) parseW(t *Tokens) Node {
// 	var a = e.parseE(t)
// 	var b Node = nil
// 	for t.i < len(t.tokens) {
// 		token := t.tokens[t.i]
// 		if token.Value == "==" || token.Value == "!=" {
// 			t.i++
// 			c := newBinaryExpr(token)
// 			b = e.parseE(t)
// 			c.Children = append(c.Children, a, b)
// 			a = c

// 		} else {
// 			break
// 		}
// 	}
// 	return a
// }

// func (e *Expr) parseE(t *Tokens) Node {
// 	var a = e.parseT(t)
// 	var b Node = nil
// 	for t.i < len(t.tokens) {
// 		if t.tokens[t.i].Value == "+" {
// 			c := newBinaryExpr(t.tokens[t.i])
// 			t.i++
// 			//consume
// 			b = e.parseT(t)
// 			c.Children = append(c.Children, a, b)
// 			a = c

// 		} else if t.tokens[t.i].Value == "-" {
// 			c := newBinaryExpr(t.tokens[t.i])
// 			t.i++
// 			b = e.parseT(t)
// 			c.Children = append(c.Children, a, b)
// 			a = c
// 		} else {
// 			break
// 		}
// 	}
// 	return a
// }
// func newBinaryExpr(lc *Token) *Expr {
// 	a := new(Expr)
// 	a.NodeType = BINARY_EXPR
// 	a.Lexeme = lc
// 	return a
// }
// func newUninaryExpr(lc *Token) *Expr {
// 	a := new(Expr)
// 	a.NodeType = UNARY_EXPR
// 	//例如 a++
// 	a.Lexeme = lc
// 	return a
// }
// func (e *Expr) parseT(t *Tokens) Node {
// 	var a = e.parseUnary_(t)
// 	var b Node = nil
// 	for t.i < len(t.tokens) {

// 		token := t.tokens[t.i]

// 		if t.tokens[t.i].Value == "*" {
// 			c := newBinaryExpr(token)
// 			t.i++
// 			b = e.parseUnary_(t)
// 			c.Children = append(c.Children, a, b)
// 			a = c
// 		} else if t.tokens[t.i].Value == "/" {
// 			c := newBinaryExpr(token)
// 			t.i++
// 			b = e.parseUnary_(t)
// 			c.Children = append(c.Children, a, b)
// 			a = c
// 		} else {
// 			break
// 		}
// 	}
// 	return a
// }
// func (e *Expr) parseUnary_(t *Tokens) Node {
// 	// 解决 a++ 的问题 , 然后 parseF
// 	if t.i >= len(t.tokens) {
// 		return nil
// 	}
// 	var anode = e.parseF(t)

// 	if t.i >= len(t.tokens) || anode == nil {
// 		return anode
// 	}
// 	//else
// 	// i< len
// 	token := t.tokens[t.i]

// 	if anode.GetLexeme().Type == lexer.Variable {
// 		// 只有变量 才能自增或者复制， 并且 目前设计是只能后置++
// 		if token.Value == "++" {
// 			// if anode.GetLexeme()
// 			t.i++
// 			top := newUninaryExpr(token)
// 			top.Children = append(top.Children, anode)
// 			anode = top
// 		} else if token.Value == "--" {
// 			t.i++
// 			top := newUninaryExpr(token)
// 			top.Children = append(top.Children, anode)
// 			anode = top
// 		} else if token.Value == "+=" {
// 			t.i++
// 			bnode := e.parseF(t)
// 			temp := newBinaryExpr(token)
// 			temp.Children = append(temp.Children, anode, bnode)
// 			anode = temp
// 		} else if token.Value == "-=" {
// 			t.i++
// 			bnode := e.parseF(t)
// 			temp := newBinaryExpr(token)
// 			temp.Children = append(temp.Children, anode, bnode)
// 			anode = temp

// 		} else if token.Value == "[" {
// 			t.i++
// 			propnode := e.parseF(t)
// 			node := new(MapIndexNode)
// 			// node.PropName = propnode
// 			if t.tokens[t.i].Value != "]" {
// 				panic("illegal state of match ]")
// 			}
// 			t.i++ // for next
// 			node.Children = append(node.Children, anode, propnode)
// 			// node.Variable = anode.(*Variable)
// 			anode = node
// 		}
// 	}

// 	return anode
// }

/*
func (e *Expr) parseF(t *Tokens) Node {
	if t.i >= len(t.tokens) {
		return nil
	}
	token := t.tokens[t.i]
	// tk := t.tokens[t.i]
	if token.Value == "(" {
		t.i++
		a := e.parseW(t)
		if t.tokens[t.i].Value != ")" {
			panic("illegal state of bracket match (")
		}
		t.i++
		return a
	}

	if token.Value == "!" {
		t.i++
		a := e.parseW(t)
		// 形如:  !(a)
		node := newUninaryExpr(token)
		node.Children = append(node.Children, a)
		return node
	}

	if token.Type == lexer.Number || token.Type == lexer.String || token.Type == lexer.Char {
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
	// if {
	// 	scalar :=
	// }
	panic("illegal node at parse")
	// scalar := new(Scalar)
	// scalar.Lexeme = t.tokens[t.i]
	// // scalar.NodeType = SCALAR
	// t.i++

}


*/
