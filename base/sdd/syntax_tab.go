package sdd

import (
	"fmt"
	"toylang/base/ast"
	"toylang/base/lexer"
)

type Node = ast.Node

func Translate(n Node) (u *Program) {
	var p Program = Program{
		NodeVar: map[Node]*Symbol{},
	}
	var t SymbolTab
	for _, c := range n.GetChildren() {
		tran_stmt(&p, c, &t)
	}
	return &p
}
func tran_stmt(p *Program, u Node, tab *SymbolTab) {
	switch { //assign stmt
	case IsAssignDeclNode(u):
		//var a = 1,  变量定义
		tran_assign_stmt(p, u, tab)
	case IsAssignExpr(u):
		//表达式
		tran_expr_stmt(p, u, tab)
	case IsBlock(u):
		tran_block(p, u, tab)
	case IsIfBlock(u):
		tran_if(p, u, tab)
	case IsFunNode(u):
		tran_fn_decl_stmt(p, u, tab)
	case IsCallFnNode(u):
		tran_call_fn_stmt(p, u, tab)
	case IsReturnStmt(u):
		tran_ret_stmt(p, u, tab)
	default:

	}
}
func tran_expr_stmt(p *Program, node Node, u *SymbolTab) {
	var p0 = u.AddSymbolByToken(node.GetChildren()[0].GetLexeme())
	var rightExpr = node.GetChildren()[1]
	tempVarP1 := _tran_expr(p, rightExpr, u)
	// if p0.Lexeme.Value == "a" {
	// fmt.Printf("%v\n", p0)
	// }
	// _tran_expr(p, rightExpr, u)
	p.AddCmd(&Cmd{
		Type:     ASSIGN,
		Result:   p0,
		Operator: "=",
		Arg1:     tempVarP1,
		Arg2:     nil,
	})
}
func tran_assign_stmt(p *Program, node Node, u *SymbolTab) {
	if u.CurrentHasToken(node.GetChildren()[0].GetLexeme()) {
		//如果 a在当前符号表出现过
		panic(fmt.Sprintf("重复定义符号 %v", node.GetChildren()[0].GetLexeme().Value))
	}
	var p0 = u.AddSymbolByToken(node.GetChildren()[0].GetLexeme())
	var rightExpr = node.GetChildren()[1]
	/*
			=
		a      +
		      b   1

		a = b+1

	*/
	// a*b+2*3
	// p1 = a*b
	// p2 = 2*3
	//p 3 = p1+ p2
	var tempVarP1 = _tran_expr(p, rightExpr, u)
	p.AddCmd(&Cmd{
		Type:     ASSIGN,
		Result:   p0,
		Operator: "=",
		Arg1:     tempVarP1,
		// Arg2:,
	})

}

func _tran_expr(p *Program, node Node, u *SymbolTab) *Symbol {
	if node == nil {
		panic("nil pointer node")
	}
	// if node.
	if IsValueNode(node) { //value

		syb := u.AddSymbolByToken(node.GetLexeme())
		p.NodeVar[node] = syb
		return syb
	} else if IsExprNode(node) { //expr
		for _, ch := range node.GetChildren() {
			_tran_expr(p, ch, u)
		}
		if _, ok := p.NodeVar[node]; !ok {
			//b = a+b,
			//p1 = a+b
			//这里存的是临时变量
			p.NodeVar[node] = u.AddVarSymbol()
		}
		chs := node.GetChildren()
		res := Cmd{
			ASSIGN,
			p.NodeVar[node],
			node.GetLexeme().Value.(string),
			p.NodeVar[chs[0]],
			p.NodeVar[chs[1]],
		}
		p.AddCmd(&res)
		return p.NodeVar[node]
	} else if IsCallFnNode(node) {
		ret := tran_call_fn_stmt(p, node, u)
		p.NodeVar[node] = ret
		return ret
	} else {

		panic("unsupport statement expr")
	}
	// panic("not impl")
	//parse call expr

}
func IsValueNode(p Node) bool {
	t := p.GetLexeme().Type
	switch t {
	case lexer.Number, lexer.Boolean, lexer.String, lexer.Char:
		return true
	case lexer.Variable:
		if IsCallFnNode(p) {
			return false
		}
		// if len(p.GetChildren()) > 0 {
		// return false
		// }
		return true
	default:
		return false
	}
}

func IsExprNode(p Node) bool {
	switch p.(type) {
	case *ast.Expr:
		{

			return true
		}
	}
	return false
}
