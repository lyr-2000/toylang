package evaluator

import (
	"fmt"
	"toylang/base/ast"
	"toylang/base/lexer"

	"github.com/spf13/cast"
)

func (h *CodeRunner) fn_call(fn *ast.FuncStmt, caller *ast.CallFuncStmt) interface{} {
	if fn == nil || caller == nil {
		return nil
	}

	var mp = make(map[string]interface{}, 0)

	h.Stack.Push(mp)
	defer h.Stack.Pop()
	if len(fn.Children) != 2 {
		panic("error fncall")
	}
	param := fn.Children[0].(*ast.FnParam)
	body := fn.Children[1].(*ast.BlockNode)

	// lib fn ,case printf ,echo then
	// fnName := fn.Lexeme.Value.(string)

	//register parameter
	for i, v := range param.Children {
		//set env
		if i < len(caller.Children) {
			h.SetVar(v.GetLexeme().Value.(string), h.evalNode(caller.Children[i]), true)
		}
	}
	// call stmt
	h.evalNode(body)
	return h.GetVar("$return")

}

func (h *CodeRunner) evalNode(n ast.Node) interface{} {
	if n == nil {
		return nil
	}
	switch n.(type) {
	case *ast.Scalar:
		w := n.(*ast.Scalar)
		if w.Lexeme.Type == lexer.Number {
			// todo: 实现转换
			return cast.ToFloat64(w.Lexeme.Value)
		}
		return w.Lexeme.Value
	case *ast.Variable:
		return h.GetVar(n.GetLexeme().Value.(string))
	case *ast.BlockNode:
		for _, c := range n.GetChildren() {
			h.evalNode(c)
		}
	case *ast.Expr:
		//left and right
		return h.evalExpr(n)
	case *ast.FuncStmt:
		//define func
		h.Functions = append(h.Functions, n.(*ast.FuncStmt))
	case *ast.CallFuncStmt:
		match := false
		var res interface{}
		for _, v := range h.Functions {
			if v.Lexeme.Value.(string) == n.GetLexeme().Value.(string) {
				res = h.fn_call(v, n.(*ast.CallFuncStmt))
				match = true
				break
			}
		}
		if !match {
			fmt.Printf("cannot call fn\n")
		}
		return res
	case *ast.ReturnStmt:
		//return statement
		m := n.(*ast.ReturnStmt)
		if len(m.Children) == 1 {
			h.SetVar("$return", h.evalNode(m.Children[0]), true)
		} else {
			h.SetVar("$return", nil, true)
		}

	default:

	}

	return nil
}

func parseVar1(a, b interface{}, op string) interface{} {
	return nil
}

func (h *CodeRunner) evalExpr(n ast.Node) interface{} {
	if n == nil {
		return nil
	}
	ch := n.GetChildren()
	word := n.GetLexeme()
	//left or right
	if len(ch) == 0 {
		return nil
	}
	if len(ch) == 1 || word == nil || word.Value == nil {
		panic("not support")
	}
	if len(ch) != 2 {
		panic("illegal state")
	}
	switch word.Value {
	case "+":
		l := h.evalNode(ch[0])
		r := h.evalNode(ch[1])
		switch l.(type) {
		case string:
			return fmt.Sprintf("%v%v", l, r)
		default:
			return cast.ToFloat64(l) + cast.ToFloat64(r)
		}
	case "-":
		return cast.ToFloat64(h.evalNode(ch[0])) - cast.ToFloat64(h.evalNode(ch[1]))
	case "*":
		return cast.ToFloat64(h.evalNode(ch[0])) * cast.ToFloat64(h.evalNode(ch[1]))
	case "/":
		b := cast.ToFloat64(h.evalNode(ch[1]))
		if b == 0 {
			panic("cannot divide by zero")
		}
		return cast.ToFloat64(h.evalNode(ch[0])) / b
	case "=":
		// a = b  => parent is = , left is a ,right is b
		// l := h.evalNode(ch[0]) // variable
		r := h.evalNode(ch[1])
		// repeat:
		// h.SetVar(cast.ToString(l))
		variable := ch[0].(*ast.Variable)
		switch r.(type) {
		case *ast.Scalar:
			r1 := r.(*ast.Scalar)
			//常量
			h.SetVar(variable.Lexeme.Value.(string), r1.Lexeme.Value, true)
		case *ast.Variable:
			r1 := r.(*ast.Variable)
			w := h.GetVar(r1.Lexeme.String())
			h.SetVar(variable.Lexeme.Value.(string), w, true)
		// case *ast.Expr:
		// 	h.SetVar(variable.Lexeme.Value.(string), 1, true)
		// r = h.evalExpr(r)
		// goto repeat
		default:
			h.SetVar(variable.Lexeme.Value.(string), r, true)
			// panic(fmt.Sprintf("illegal node support %T", r))

		}
		// h.SetVar(variable.Lexeme.String(), r, false)
	default:

	}
	return nil
}

func parse_source_tree(s string) ast.Anode {
	var lx = lexer.NewStringLexer(s)

	tt := lx.ReadTokens()
	b := ast.NewTokens(tt)
	tree := ast.ParseStmt(b)
	return tree

}
