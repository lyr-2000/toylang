package evaluator

import (
	"fmt"
	"github.com/lyr-2000/toylang/base/ast"
	"github.com/lyr-2000/toylang/base/lexer"
	"strings"

	"github.com/spf13/cast"
)

const (
	returnKey = "@return"
	breakKey  = "@break"
)

func panicf(s string, o ...interface{}) {
	panic(fmt.Sprintf(s, o...))
}

type arrayVal struct {
	Value []interface{}
}

func (r *arrayVal) Min() interface{} {
	if len(r.Value) == 0 {
		panic("array is empty")
	}
	var f = cast.ToFloat64(r.Value[0])
	for _, v := range r.Value {
		d := cast.ToFloat64(v)
		if d < f {
			f = d
		}
	}
	return f
}
func (r *arrayVal) Max() interface{} {
	if len(r.Value) == 0 {
		panic("array is empty")
	}
	var f = cast.ToFloat64(r.Value[0])
	for _, v := range r.Value {
		d := cast.ToFloat64(v)
		if d > f {
			f = d
		}
	}
	return f
}

func (h *CodeRunner) evalBool(bh ast.Node) bool {
	if bh == nil {
		return false
	}

	ret := h.EvalNode(bh)
	switch ret.(type) {
	case float64:
		return cast.ToBool(int(ret.(float64)))
	default:

	}
	return cast.ToBool(ret)

}

type RefValue struct {
	Value   interface{}
	setFunc (func(interface{}))
}

func (r *RefValue) SetValue(value interface{}) {
	r.Value = value
	if r.setFunc != nil {
		r.setFunc(value)
	}
}

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

			paramVal := h.EvalNode(caller.Children[i+1])
			//fmt.Printf("k %v,v %v \n", paramKey, paramVal)
			h.SetVar(paramKey, paramVal, true)
		}
	}
	// call stmt
	h.EvalNode(body)
	return h.GetVar(returnKey)

}

func (h *CodeRunner) EvalNode(n ast.Node) interface{} {
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
		h.DebugLog.Printf("declare %#v\n", w)
		if len(w.Children) == 2 {
			h.SetVar(w.Children[0].GetLexeme().Value.(string), h.EvalNode(w.Children[1]), true)
		} else {
			h.SetVar(w.Children[0].GetLexeme().Value.(string), nil, true)
		}

	case *ast.Scalar:
		w := n.(*ast.Scalar)
		h.DebugLog.Printf("evalScalar %+v\n", w)
		return cast_scalar_node_type(w)
	case *ast.Variable:
		return h.GetVar(n.GetLexeme().Value.(string))
	case *ast.BlockNode:
		for _, c := range n.GetChildren() {
			if ret, isRet := c.(*ast.ReturnStmt); isRet {
				// block  里面 有 return，就不执行后面的代码
				return h.EvalNode(ret)
			}
			h.EvalNode(c)
		}
	case *ast.MapIndexNode:
		m := n.(*ast.MapIndexNode)
		arr := h.EvalNode(m.GetChildren()[0])
		key := h.EvalNode(m.GetChildren()[1])
		h.DebugLog.Printf("evalMapIndexNode %+v, %+v\n", arr, key)
		if _, ok := arr.([]interface{}); ok {
			return arr.([]interface{})[cast.ToInt(key)]
		}
		if mp, ok := arr.(map[string]interface{}); ok {
			ret := mp[cast.ToString(key)]
			return ret
		}
		panic(fmt.Sprintf("illegal array %+v", arr))
	case *ast.IfStmt:

		stmt := n.(*ast.IfStmt)
		cond := stmt.GetCondition()
		body := stmt.GetBody()
		elseNode := stmt.GetElseNode()
		h.DebugLog.Printf("evalIfStmt %+v, %+v, %+v\n", cond, body, elseNode)
		val := h.EvalNode(cond)
		ok := cast.ToBool(val)
		if ok {
			return h.EvalNode(body)
		}
		return h.EvalNode(elseNode)

	case *ast.Expr:
		//left and right
		pe := h.evalExpr(n)
		h.DebugLog.Printf("!astExpr: %+v, %+v\n", n, pe)
		if f, ok := pe.(bool); ok {
			h.setPrevEval(f)
		}
		return pe
	case *ast.ExprGroups:
		var last interface{}
		ln := len(n.GetChildren())
		for i, v := range n.GetChildren() {
			if i == ln-1 {
				last = h.EvalNode(v)
			} else {
				h.EvalNode(v)
			}
		}
		// a==1,b==2,c==3, 逗号表达式，最后返回 c==3 结果
		return last
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
		res1, ok := h.CallInline(fnName, n.GetChildren())
		if ok {
			return res1
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
			h.DebugLog.Panicf("cannot call fn " + fnName)
		}
		return res
	case *ast.ReturnStmt:
		//return statement
		m := n.(*ast.ReturnStmt)
		if len(m.Children) == 1 {
			h.SetVar(returnKey, h.EvalNode(m.Children[0]), true)
		} else {
			h.SetVar(returnKey, nil, true)
		}

	case *ast.ForStmt:
		m := n.(*ast.ForStmt)
		init := m.GetInitNode()
		cond := m.GetCondition()
		le := m.GetEachLoopAction()
		body := m.GetBody()

		for cond == nil || h.evalBool(cond) {
			if init != nil {
				h.EvalNode(init)
				init = nil
			}
			h.EvalNode(body)
			_, exists := h.GetStackVar(returnKey)
			if exists {

				// return
				return nil
			}

			_, exists = h.GetVar2(breakKey)
			if exists {
				h.DelVar(breakKey)
				//	fmt.Printf("%+v\n", aaa)
				break
				// return nil
			}
			//for next
			h.EvalNode(le)
		}
	case *ast.BreakFlagStmt:
		h.SetVar(breakKey, 1, true)
		h.DebugLog.Printf("evalBreakFlagStmt %+v\n", n)

	default:

	}

	return nil
}

func (h *CodeRunner) setPrevEval(b bool) {
	if b {
		h.PrevEval = 1
	} else {
		h.PrevEval = 0
	}
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
	if /* len(ch) == 1 || */ word == nil || word.Value == nil {
		h.DebugLog.Panicf("not support" + fmt.Sprintf("vl=%+v", n))
	}
	if len(ch) == 1 {
		//处理一元运算符 ，i++
		if n.GetLexeme().Value == "++" {
			// r := h.evalNode(ch[1])
			r := h.EvalNode(ch[0])
			var res = cast.ToFloat64(r) + 1
			// switch r.(type) {
			// case string:
			// 	res = fmt.Sprintf("%v%v", l, r)
			// default:
			// 	res = cast.ToFloat64(l) + cast.ToFloat64(r)
			// }
			h.SetVar(ch[0].(*ast.Variable).Lexeme.Value.(string), res, false)
			//处理自增
			return res
		}
		if ExtraOp != nil {
			exp := cast.ToString(n.GetLexeme().Value)
			if fn, ok := ExtraOp[exp]; ok {
				return fn(h, n)
			}
		}
	}

	if len(ch) != 2 {
		panic("illegal state")
	}
	switch word.Value {
	case "+":
		l := h.EvalNode(ch[0])
		r := h.EvalNode(ch[1])
		switch l.(type) {
		case string:
			return fmt.Sprintf("%v%v", l, r)
		default:
			return cast.ToFloat64(l) + cast.ToFloat64(r)
		}
	case "-":
		return cast.ToFloat64(h.EvalNode(ch[0])) - cast.ToFloat64(h.EvalNode(ch[1]))
	case "*":
		return cast.ToFloat64(h.EvalNode(ch[0])) * cast.ToFloat64(h.EvalNode(ch[1]))
	case "/":
		b := cast.ToFloat64(h.EvalNode(ch[1]))
		if b == 0 {
			panic("cannot divide by zero")
		}
		return cast.ToFloat64(h.EvalNode(ch[0])) / b
	case "||":
		l := h.EvalNode(ch[0])
		lb := cast.ToBool(l)
		if lb {
			return true
		}
		//短路 或
		r := h.EvalNode(ch[1])
		return cast.ToBool(r)
	case "&&":
		// panic("unsupport operation")
		l := h.EvalNode(ch[0])
		lb := cast.ToBool(l)
		if !lb {
			return false
		}
		//短路 或
		r := h.EvalNode(ch[1])
		ret := cast.ToBool(r)
		return ret
	case "+=":
		//a+=1 =>  a = a+ 1
		l := h.EvalNode(ch[0])
		r := h.EvalNode(ch[1])
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
		l := h.EvalNode(ch[0])
		r := h.EvalNode(ch[1])
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
		l := h.EvalNode(ch[0])
		r := h.EvalNode(ch[1])
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
		l := h.EvalNode(ch[0])
		r := h.EvalNode(ch[1])
		var res bool
		switch l.(type) {
		case nil:
			return r == nil

		case string:
			// a.compare b => -1 => a-b == -1  ==> a < b
			res = strings.Compare(l.(string), fmt.Sprintf("%+v", r)) == 0

		case float64, int64, uint64, uint32, int32, uint16, int16, uint8, int8:
			res = cast.ToFloat64(l) == cast.ToFloat64(r)
		default:
			res = cast.ToString(l) == cast.ToString(r)
		}
		return res
	case ">":
		l := h.EvalNode(ch[0])
		r := h.EvalNode(ch[1])
		var res bool
		switch l.(type) {
		case string:
			// a.compare b => -1 => a-b == -1  ==> a < b
			res = strings.Compare(l.(string), fmt.Sprintf("%+v", r)) > 0
		case nil:
			return false

		default:
			res = cast.ToFloat64(l) > cast.ToFloat64(r)
		}
		return res
	case "<":
		l := h.EvalNode(ch[0])
		r := h.EvalNode(ch[1])
		var res bool
		switch l.(type) {
		case nil:
			return false
		case string:
			// a.compare b => -1 => a-b == -1  ==> a < b
			res = strings.Compare(l.(string), fmt.Sprintf("%+v", r)) < 0

		default:
			res = cast.ToFloat64(l) < cast.ToFloat64(r)
		}
		return res
	case "=":
		r := h.EvalNode(ch[1])
		variable := ch[0].(*ast.Variable)
		switch r.(type) {
		case *ast.Scalar:
			r1 := r.(*ast.Scalar)
			//常量
			h.SetVar(variable.Lexeme.Value.(string), r1, true)
		case *ast.Variable:
			r1 := r.(*ast.Variable)
			w := h.GetVar(r1.Lexeme.String())
			h.SetVar(variable.Lexeme.Value.(string), w, true)
			// goto repeat
		default:
			h.SetVar(variable.Lexeme.Value.(string), r, true)
			// panic(fmt.Sprintf("illegal node support %T", r))

		}

		// h.SetVar(variable.Lexeme.String(), r, false)
	default:
		if ExtraOp != nil {
			exp := cast.ToString(word.Value)
			if fn, ok := ExtraOp[exp]; ok {
				return fn(h, n)
			}
		}
		panic(fmt.Sprintf("unsupport operation %+v ,%+v", n, word))

	}
	return nil
}

var (
	ExtraOp map[string]func(h *CodeRunner, node ast.Node) any
)

func SetExtrapOp(op string, fn func(h *CodeRunner, node ast.Node) any) {
	if ExtraOp == nil {
		ExtraOp = make(map[string]func(h *CodeRunner, node ast.Node) any)
	}
	ExtraOp[op] = fn
}

func ParseSourceTree(s string) ast.Anode {
	return parseSourceTree(s)
}

func parseSourceTree(s string) ast.Anode {
	var lx = lexer.NewStringLexer(s)

	tt := lx.ReadTokens()
	// checkToken
	CheckToken(tt)
	b := ast.NewTokens(tt)
	tree := ast.ParseStmt(b)
	return tree

}
func CheckToken(tokens []*lexer.Token) bool {
	var con uint8
	for _, v := range tokens {
		if v.Type == lexer.Variable {
			con++
		} else {
			con = 0
		}
		if con >= 3 {
			panicf("illegal token %+v", tokens)
		}
	}
	return true
}

func ParseTree(s string) ast.Anode {
	return parseSourceTree(s)
}
