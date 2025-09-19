package ast

import (
	"fmt"
	"github.com/lyr-2000/toylang/base/lexer"
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
func (t *Tokens) peekNextToken(id int) *Token {
	if id+t.i >= len(t.tokens) {
		return nil
	}
	return t.tokens[t.i+id]
}
func (t *Tokens) peekNextString(id int) string {
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
	return pareseExprFunc(t)
}
func ParseExpr_(t []*Token) Node {
	return pareseExprFunc(&Tokens{t, 0})
}
func pareseExprFunc(t *Tokens) Node {
	var e Expr
	node := e.parseE(t, 0)
	// log.Printf("read info %+v, %v\n", t.i, len(t.tokens))
	return node
}

var (
	OpTables = [][]string{
		0: {
			"=", "+=", "-=", "*=", "/=", "%=", "&=", "|=", "^=", "!=",
		},
		1: {
			"||",
		},
		2: {
			"&&",
		},
		3: {
			"|",
		},
		4: {
			"^",
		},
		5: {
			"&",
		},
		6: {
			"==", "!=",
		},
		7: {
			">=", "<=", ">", "<",
		},
		8: {
			">>", "<<",
		},
		9: {
			"+", "-",
		},
		10: {
			"*", "/", "%",
		},
	}
)

func RegisterOpSign(priorityIndex int, sign string) {
	for _, v := range OpTables[priorityIndex] {
		if v == sign {
			return
		}
	}
	OpTables[priorityIndex] = append(OpTables[priorityIndex], sign)
}

//1.处理伪运算符，  -> , \\.  , [] 等 ，优先级最高
//2. 处理单目运算符， ++ ,  次高

func (e *Expr) parseE(t *Tokens, k int) Node {
	if k >= len(OpTables) {

		return e.parseTop2(t)
	}
	var a = e.parseE(t, k+1)

	for t.i < len(t.tokens) {
		cur := t.tokens[t.i]
		if k == 0 {
			if cur == nil || cur.Value == ";" || cur.Value == ")" { //end line
				// t.i++
				break
			}
			if cur.Value == "{" || cur.Value == "," {
				break
			}

		}
		match := false
		for _, v := range OpTables[k] {
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
func isMixToken(token *Token) bool {
	sop := token.Value.(string)
	if len(sop) == 2 && lexer.IsOperator(lexer.CharNum(sop[0])) && lexer.MixOpDefine != nil {
		l, r := lexer.CharNum(sop[0]), lexer.CharNum(sop[1])
		// d := lexer.MakeString(l, r)
		if str := lexer.MixOpDefine(l, r); str != "" {
			return true
		}
	}
	return false
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
	} else if token.Value == "!" || token.Value == "-" || token.Value == "+" || isMixToken(token) {
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

	if token.Type == lexer.Variable {
		ne := t.peekNextString(1)
		if ne == "(" {
			//ok
			return parseCallFuncStmt(t)
		}
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
	if token.Value == "," {
		return nil
	}
	if token.Value == ")" {
		return nil
	}
	
	panic(fmt.Sprintf("cannot explain ast tree %+v", token)) // return nil
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
