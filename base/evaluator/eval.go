package evaluator

import (
	"fmt"
	"strings"
	"toylang/base/ast"
	"toylang/base/lexer"

	"github.com/spf13/cast"
)

const (
	returnKey = "@return"
)

// func (h *CodeRunner) getVar2(key string) (interface{},present) {
// 	return
// }
func (h *CodeRunner) fn_call(fn *ast.FuncStmt, caller *ast.CallFuncStmt) interface{} {
	if fn == nil || caller == nil {
		return nil
	}
	if h.Stack.Len() > 444 {
		panic("stack overflow error !!!")
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
			if i+1 >= len(caller.Children) {
				break
				// it is no param for next
			}
			paramKey := v.GetLexeme().Value.(string)

			paramVal := h.evalNode(caller.Children[i+1])
			//fmt.Printf("k %v,v %v \n", paramKey, paramVal)
			h.SetVar(paramKey, paramVal, true)
		}
	}
	// call stmt
	h.evalNode(body)
	return h.GetVar(returnKey)

}

func (h *CodeRunner) evalNode(n ast.Node) interface{} {
	if n == nil {
		return nil
	}
	if h.stackDeep > 0 {
		// fn call
		_, exists := h.GetStackVar(returnKey)
		if exists {
			//	fmt.Printf("%+v\n", aaa)
			return nil
		}
	}

	switch n.(type) {
	case *ast.DeclareStmt:
		w := n.(*ast.DeclareStmt)
		fmt.Printf("%#v\n", w)
		if len(w.Children) == 2 {
			h.SetVar(w.Children[0].GetLexeme().Value.(string), h.evalNode(w.Children[1]), true)
		} else {
			h.SetVar(w.Children[0].GetLexeme().Value.(string), nil, true)
		}

	case *ast.Scalar:
		w := n.(*ast.Scalar)
		return cast_scalar_node_type(w)
	case *ast.Variable:
		return h.GetVar(n.GetLexeme().Value.(string))
	case *ast.BlockNode:
		for _, c := range n.GetChildren() {
			if ret, isRet := c.(*ast.ReturnStmt); isRet {
				// block  里面 有 return，就不执行后面的代码
				return h.evalNode(ret)
			}
			h.evalNode(c)
		}
	case *ast.IfStmt:
		stmt := n.(*ast.IfStmt)
		cond := stmt.GetCondition()
		body := stmt.GetBody()
		elseNode := stmt.GetElseNode()

		val := h.evalNode(cond)
		ok := cast.ToBool(val)
		if ok {
			return h.evalNode(body)
		}
		return h.evalNode(elseNode)

	case *ast.Expr:
		//left and right
		return h.evalExpr(n)
	case *ast.FuncStmt:
		//定义函数
		h.fn_define(n.(*ast.FuncStmt))
		// h.Functions = append(h.Functions, n.(*ast.FuncStmt))
	case *ast.CallFuncStmt:

		var fnName = n.GetLexeme().Value.(string)
		if IsLibFn(fnName) {
			// lib call
			return h.libFnCall(fnName, n.GetChildren())
		}
		match := false
		var res interface{}

		// call  user fun
		for _, v := range h.Functions {
			if v.Lexeme.Value.(string) == n.GetLexeme().Value.(string) {
				// match it
				h.stackDeep++
				defer func() {
					h.stackDeep--
				}()
				res = h.fn_call(v, n.(*ast.CallFuncStmt))
				match = true
				break
			}
		}
		if !match {
			// fmt.Printf("cannot call fn\n")
			panic("cannot call fn")
		}
		return res
	case *ast.ReturnStmt:
		//return statement
		m := n.(*ast.ReturnStmt)
		if len(m.Children) == 1 {
			h.SetVar(returnKey, h.evalNode(m.Children[0]), true)
		} else {
			h.SetVar(returnKey, nil, true)
		}

	default:

	}

	return nil
}

// func parseVar1(a, b interface{}, op string) interface{} {
// 	return nil
// }

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
	case "||":
		l := h.evalNode(ch[0])
		lb := cast.ToBool(l)
		if lb {
			return true
		}
		//短路 或
		r := h.evalNode(ch[1])
		return cast.ToBool(r)
	case "&&":
		panic("unsupport operation")
	case "+=":
		//a+=1 =>  a = a+ 1
		l := h.evalNode(ch[0])
		r := h.evalNode(ch[1])
		var res interface{}
		switch l.(type) {
		case string:
			res = fmt.Sprintf("%v%v", l, r)
		default:
			res = cast.ToFloat64(l) + cast.ToFloat64(r)
		}
		h.SetVar(ch[0].(*ast.Variable).Lexeme.Value.(string), res, false)
		return res
	case "<=":
		l := h.evalNode(ch[0])
		r := h.evalNode(ch[1])
		var res bool
		switch l.(type) {

		case string:
			// a.compare b => -1 => a-b == -1  ==> a < b
			res = strings.Compare(l.(string), fmt.Sprintf("%+v", r)) <= 0

		default:
			res = cast.ToFloat64(l) <= cast.ToFloat64(r)
		}
		return res
	case ">=":
		l := h.evalNode(ch[0])
		r := h.evalNode(ch[1])
		var res bool
		switch l.(type) {

		case string:
			// a.compare b => -1 => a-b == -1  ==> a < b
			res = strings.Compare(l.(string), fmt.Sprintf("%+v", r)) >= 0

		default:
			res = cast.ToFloat64(l) >= cast.ToFloat64(r)
		}
		return res
	case "==":
		l := h.evalNode(ch[0])
		r := h.evalNode(ch[1])
		var res bool
		switch l.(type) {

		case string:
			// a.compare b => -1 => a-b == -1  ==> a < b
			res = strings.Compare(l.(string), fmt.Sprintf("%+v", r)) == 0

		default:
			res = cast.ToFloat64(l) == cast.ToFloat64(r)
		}
		return res
	case ">":
		l := h.evalNode(ch[0])
		r := h.evalNode(ch[1])
		var res bool
		switch l.(type) {

		case string:
			// a.compare b => -1 => a-b == -1  ==> a < b
			res = strings.Compare(l.(string), fmt.Sprintf("%+v", r)) > 0

		default:
			res = cast.ToFloat64(l) > cast.ToFloat64(r)
		}
		return res
	case "<":
		l := h.evalNode(ch[0])
		r := h.evalNode(ch[1])
		var res bool
		switch l.(type) {

		case string:
			// a.compare b => -1 => a-b == -1  ==> a < b
			res = strings.Compare(l.(string), fmt.Sprintf("%+v", r)) < 0

		default:
			res = cast.ToFloat64(l) < cast.ToFloat64(r)
		}
		return res
	case "=":
		r := h.evalNode(ch[1])
		variable := ch[0].(*ast.Variable)
		switch r.(type) {
		case *ast.Scalar:
			r1 := r.(*ast.Scalar)
			//常量
			h.SetVar(variable.Lexeme.Value.(string), r1.Lexeme.Value, false)
		case *ast.Variable:
			r1 := r.(*ast.Variable)
			w := h.GetVar(r1.Lexeme.String())
			h.SetVar(variable.Lexeme.Value.(string), w, false)
		// case *ast.Expr:
		// 	h.SetVar(variable.Lexeme.Value.(string), 1, true)
		// r = h.evalExpr(r)
		// goto repeat
		default:
			h.SetVar(variable.Lexeme.Value.(string), r, false)
			// panic(fmt.Sprintf("illegal node support %T", r))

		}

		// h.SetVar(variable.Lexeme.String(), r, false)
	default:
		panic(fmt.Sprintf("unsupport operation %+v ,%+v", n, word))

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
