package evaluator

import (
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"reflect"
	"slices"
	"strings"

	"github.com/lyr-2000/toylang/base/ast"
	"github.com/lyr-2000/toylang/base/lexer"
	"github.com/lyr-2000/toylang/base/list"

	"github.com/spf13/cast"
)

type CodeRunner struct {
	Functions []*ast.FuncStmt
	Inlines   map[string]func([]interface{}) interface{}
	Vars      map[string]interface{}
	Stack     *list.Stack
	ExitCode  uint8
	ErrCode   uint64
	ErrMsg    string
	PrevEval  uint8
	stackDeep uint32

	DebugLog *log.Logger
	vmOutput io.Writer
	VmLogger *log.Logger
}

func (r *CodeRunner) SetOutput(w io.Writer) {
	r.vmOutput = w
	if r.VmLogger == nil {
		r.VmLogger = log.New(w, "[VM]", log.LstdFlags)
	} else {
		r.VmLogger.SetOutput(w)
	}
}

func (r *CodeRunner) Alias(from, to string) {
	r.Inlines[to] = r.Inlines[from]
}
func (r *CodeRunner) SetFunc(name string, fn func([]interface{}) interface{}) {
	if r.Inlines == nil {
		r.Inlines = make(map[string]func([]interface{}) interface{})
	}
	r.Inlines[name] = fn
}

func init() {
	lexer.MixOpDefine = func(l, r lexer.CharNum) string {
		df := lexer.MakeString(l, r)
		_, ok := ExtraOp[df]
		if ok {
			return df
		}
		return ""
	}
	SetExtrapOp("!=", func(h *CodeRunner, node ast.Node) any {
		l := h.EvalNode(node.GetChildren()[0])
		r := h.EvalNode(node.GetChildren()[1])
		return l != r
	})
	// a<<1
	// b>>1
	SetExtrapOp("<<", func(h *CodeRunner, node ast.Node) any {
		l := h.EvalNode(node.GetChildren()[0])
		r := h.EvalNode(node.GetChildren()[1])
		d := cast.ToInt64(l) << cast.ToInt64(r)
		return d
	})
	SetExtrapOp(">>", func(h *CodeRunner, node ast.Node) any {
		l := h.EvalNode(node.GetChildren()[0])
		r := h.EvalNode(node.GetChildren()[1])
		d := cast.ToInt64(l) >> cast.ToInt64(r)
		return d
	})
	SetExtrapOp("%=", func(h *CodeRunner, node ast.Node) any {
		left := node.GetChildren()[0]
		l := h.EvalNode(node.GetChildren()[0])
		r := h.EvalNode(node.GetChildren()[1])
		d := cast.ToInt64(l) % cast.ToInt64(r)
		v := ast.AsVariable(left)
		if v != nil {
			h.SetVar(v.GetVarName(), d, true)
		}
		return d
	})
	SetExtrapOp("*=", func(h *CodeRunner, node ast.Node) any {
		left := node.GetChildren()[0]
		l := h.EvalNode(node.GetChildren()[0])
		r := h.EvalNode(node.GetChildren()[1])
		d := cast.ToInt64(l) * cast.ToInt64(r)
		v := ast.AsVariable(left)
		if v != nil {
			h.SetVar(v.GetVarName(), d, true)
		}
		return d
	})
	SetExtrapOp("/=", func(h *CodeRunner, node ast.Node) any {
		left := node.GetChildren()[0]
		l := h.EvalNode(node.GetChildren()[0])
		r := h.EvalNode(node.GetChildren()[1])
		d := cast.ToInt64(l) / cast.ToInt64(r)
		v := ast.AsVariable(left)
		if v != nil {
			h.SetVar(v.GetVarName(), d, true)
		}
		return d
	})
	SetExtrapOp("++", func(h *CodeRunner, node ast.Node) any {
		left := node.GetChildren()[0]
		l := h.EvalNode(node.GetChildren()[0])
		v := ast.AsVariable(left)
		if v != nil {
			h.SetVar(v.GetVarName(), cast.ToInt64(l)+1, true)
		}
		return l
	})
	SetExtrapOp("--", func(h *CodeRunner, node ast.Node) any {
		left := node.GetChildren()[0]
		l := h.EvalNode(node.GetChildren()[0])
		v := ast.AsVariable(left)
		if v != nil {
			h.SetVar(v.GetVarName(), cast.ToInt64(l)-1, true)
		}
		return l
	})
	SetExtrapOp("~=", func(h *CodeRunner, node ast.Node) any {
		l := h.EvalNode(node.GetChildren()[0])
		r := h.EvalNode(node.GetChildren()[1])
		return l != r
	})
	SetExtrapOp("^=", func(h *CodeRunner, node ast.Node) any {
		l := h.EvalNode(node.GetChildren()[0])
		r := h.EvalNode(node.GetChildren()[1])
		l1 := cast.ToInt64(l)
		r2 := cast.ToInt64(r)
		l1 ^= r2
		v := ast.AsVariable(node.GetChildren()[0])
		if v != nil {
			h.SetVar(v.GetVarName(), l1, true)
		}
		return l1
	})
	SetExtrapOp("&=", func(h *CodeRunner, node ast.Node) any {
		l := h.EvalNode(node.GetChildren()[0])
		r := h.EvalNode(node.GetChildren()[1])
		l1 := cast.ToInt64(l)
		r2 := cast.ToInt64(r)
		l1 &= r2
		v := ast.AsVariable(node.GetChildren()[0])
		if v != nil {
			h.SetVar(v.GetVarName(), l1, true)
		}
		return l1
	})
	SetExtrapOp("|=", func(h *CodeRunner, node ast.Node) any {
		l := h.EvalNode(node.GetChildren()[0])
		r := h.EvalNode(node.GetChildren()[1])
		l1 := cast.ToInt64(l)
		r2 := cast.ToInt64(r)
		l1 |= r2
		v := ast.AsVariable(node.GetChildren()[0])
		if v != nil {
			h.SetVar(v.GetVarName(), l1, true)
		}
		return l1
	})
}

func setFunc(r *CodeRunner) {
	r.SetFunc("fatal", func(params []interface{}) interface{} {
		msg := fmt.Sprintf("fatal call: %v", params[0])
		r.ErrCode = 1
		r.ErrMsg = msg
		r.VmLogger.Panicf(msg)
		return nil
	})
	r.SetFunc("recover", func(params []interface{}) interface{} {
		if r.ErrCode == 0 {
			return nil
		}
		ret := map[string]any{
			"errCode": r.ErrCode,
			"errMsg":  r.ErrMsg,
		}
		r.ErrCode = 0
		r.ErrMsg = ""
		return ret
	})
	// throw(1001,"errorMsg")
	r.SetFunc("throw", func(params []interface{}) interface{} {
		if len(params) <= 0 {
			r.ErrCode = 1
			r.ErrMsg = "throw error"
		}
		if len(params) <= 1 {
			r.ErrCode = 1
			r.ErrMsg = cast.ToString(params[0])
			if r.ErrMsg == "" {
				r.ErrMsg = "throw error"
			}
		}
		r.ErrCode = cast.ToUint64(params[0])
		r.ErrMsg = cast.ToString(params[1])
		if r.ErrCode == 0 {
			r.ErrCode = 1
		}
		ret := map[string]any{
			"errCode": r.ErrCode,
			"errMsg":  r.ErrMsg,
		}
		return ret
	})
	r.SetFunc("int64", func(params []interface{}) interface{} {
		p := params[0]
		return cast.ToInt64(p)
	})
	r.Alias("int64", "int")
	r.SetFunc("float64", func(params []interface{}) interface{} {
		p := params[0]
		switch val := p.(type) {
		case string:
			if strings.Contains(val, "e") || strings.Contains(val, "E") {
				bg := big.NewFloat(0)
				bg.SetString(val)
				w, _ := bg.Float64()
				return w
			}
			return cast.ToFloat64(val)
		}
		return cast.ToFloat64(p)
	})
	r.Alias("float64", "float")
	r.SetFunc("typeof", func(params []interface{}) interface{} {
		return reflect.TypeOf(params[0]).String()
	})
	r.SetFunc("array", func(params []interface{}) interface{} {
		return params
	})
	r.SetFunc("len", func(params []interface{}) interface{} {
		if arr, ok := params[0].([]interface{}); ok {
			return len(arr)
		}
		if mp, ok := params[0].(map[string]interface{}); ok {
			return len(mp)
		}
		if str, ok := params[0].(string); ok {
			return len(str)
		}
		r.VmLogger.Panicf("len param must be array or map or string,but got %T", params[0])
		return nil
	})
	r.SetFunc("del", func(params []interface{}) interface{} {
		if arr, ok := params[0].([]interface{}); ok {
			arr = append(arr[:params[1].(int)], arr[params[1].(int)+1:]...)
			return arr
		}
		if mp, ok := params[0].(map[string]interface{}); ok {
			delete(mp, params[1].(string))
			return mp
		}
		r.VmLogger.Panicf("del param must be array or map")
		return nil
	})
	r.SetFunc("map", func(params []interface{}) interface{} {
		mp := make(map[string]interface{})
		for i := 0; i+1 < len(params); i += 2 {
			key := params[i]
			value := params[i+1]
			mp[cast.ToString(key)] = value
		}
		return mp
	})
	r.SetFunc("set", func(params []interface{}) interface{} {
		if arr, ok := params[0].([]interface{}); ok {
			arr[cast.ToInt(params[1])] = params[2]
			return arr
		}
		if mp, ok := params[0].(map[string]interface{}); ok {
			mp[cast.ToString(params[1])] = params[2]
			return mp
		}

		val := reflect.Indirect(reflect.ValueOf(params[0]))
		typ := val.Kind()
		switch typ {
		case reflect.Struct:
			val.FieldByName(cast.ToString(params[1])).Set(reflect.ValueOf(params[2]))
			return params[0]
		}
		r.VmLogger.Panicf("set param must be array or map")
		return nil
	})
	r.SetFunc("get", func(params []interface{}) interface{} {
		if arr, ok := params[0].([]interface{}); ok {
			return arr[cast.ToInt(params[1])]
		}
		if mp, ok := params[0].(map[string]interface{}); ok {
			return mp[cast.ToString(params[1])]
		}
		val := reflect.Indirect(reflect.ValueOf(params[0]))
		typ := val.Kind()
		switch typ {
		case reflect.Struct:
			return val.FieldByName(cast.ToString(params[1])).Interface()
		}
		r.VmLogger.Panicf("get param must be array or map")
		return nil
	})
	r.SetFunc("add", func(params []interface{}) interface{} {
		f, ok := params[0].([]interface{})
		if ok {
			f = append(f, params[1])
		}
		return f
	})
	r.SetFunc("pop", func(params []interface{}) interface{} {
		f, ok := params[0].([]interface{})
		if ok {
			f = f[:len(f)-1]
		}
		return f
	})
	r.SetFunc("contains", func(params []interface{}) interface{} {
		w := slices.Contains(params[0].([]interface{}), params[1])
		return w
	})
	r.SetFunc("trimPrefix", func(params []interface{}) interface{} {
		return strings.TrimPrefix(params[0].(string), params[1].(string))
	})
	r.SetFunc("trimSuffix", func(params []interface{}) interface{} {
		return strings.TrimSuffix(params[0].(string), params[1].(string))
	})
	r.SetFunc("trim", func(params []interface{}) interface{} {
		return strings.Trim(params[0].(string), params[1].(string))
	})
	r.SetFunc("min", func(params []interface{}) interface{} {
		d := (arrayVal{Value: params})
		return d.Min()
	})
	r.SetFunc("max", func(params []interface{}) interface{} {
		d := arrayVal{Value: params}
		return d.Max()
	})

}
func NewCodeRunner() *CodeRunner {
	var r = new(CodeRunner)
	r.Vars = make(map[string]interface{}, 0)
	r.Stack = list.NewStack()
	setFunc(r)
	r.DebugLog = log.New(io.Discard, "[Evaluator]", log.LstdFlags)
	r.SetOutput(os.Stdout)
	return r
}

func (h *CodeRunner) GetVar(key string) interface{} {
	return get_var(h.Vars, h.Stack, key)
}

func (h *CodeRunner) GetVar2(key string) (interface{}, bool) {
	return get_var2(h.Vars, h.Stack, key)
}
func (h *CodeRunner) DelVar(key string) (interface{}, bool) {
	return del_var2(h.Vars, h.Stack, key)
}
func (h *CodeRunner) GetStackVar(key string) (interface{}, bool) {
	return get_stack_var(h.Vars, h.Stack, key)
}

func (h *CodeRunner) SetVar(k string, v interface{}, isStackAssign bool) {
	if k == "" {
		return
	}
	if v == nil {
		set_var(h.Vars, h.Stack, k, v, isStackAssign)
		return
	}
	if n, ok := v.(ast.Anode); ok {
		switch n.(type) {
		case *ast.Variable:
			varnode := n.(*ast.Variable)
			h.SetVar(k, h.GetVar(varnode.Lexeme.Value.(string)), isStackAssign)
		case *ast.Scalar:
			varnode := n.(*ast.Scalar)
			switch varnode.Lexeme.Type {
			case lexer.Number:
				h.SetVar(k, cast.ToFloat64(varnode.Lexeme.Value), isStackAssign)

			default:
				h.SetVar(k, varnode.Lexeme.Value.(string), isStackAssign)
			}
			// h.SetVar(k, varnode, isStk)
		default:
			h.SetVar(k, h.EvalNode(v.(ast.Anode)), isStackAssign)
		}

		return
	}

	set_var(h.Vars, h.Stack, k, v, isStackAssign)
}

func (c *CodeRunner) ParseAndRun(s string) int {
	ts := parseSourceTree(s)
	return c.RunCode(ts)
}

func (c *CodeRunner) ParseAndRunRecover(s string) (int, error) {
	ts := parseSourceTree(s)
	return c.RunCodeRecover(ts)
}

func (c *CodeRunner) RunCode(t ast.Node) int {
	if t == nil {
		return 0
	}
	c.EvalNode(t)
	return int(c.ExitCode)
}

func (c *CodeRunner) RunCodeRecover(t ast.Node) (ex int, err error) {
	defer func() {
		if r := recover(); r != nil {
			c.VmLogger.Printf("panic: %v", r)
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	return c.RunCode(t), err
}
