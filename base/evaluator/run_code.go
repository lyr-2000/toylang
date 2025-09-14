package evaluator

import (
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/lyr-2000/toylang/base/ast"
	"github.com/lyr-2000/toylang/base/lexer"
	"github.com/lyr-2000/toylang/base/list"

	"github.com/spf13/cast"
)

type CodeRunner struct {
	Functions []*ast.FuncStmt
	Inlines map[string]func([]interface{}) interface{}
	Vars      map[string]interface{}
	Stack     *list.Stack
	ExitCode uint8
	PrevEval uint8
	stackDeep int

	DebugLog *log.Logger
	vmOutput io.Writer
	VmLogger *log.Logger
}
func (r *CodeRunner) SetOutput(w io.Writer) {
	r.vmOutput = w
	if r.VmLogger == nil {
		r.VmLogger = log.New(w, "[VM]", log.LstdFlags)
	}else {
		r.VmLogger.SetOutput(w)
	}
}

func (r *CodeRunner) SetFunc(name string,fn func([]interface{}) interface{}){
	if r.Inlines == nil {
		r.Inlines = make(map[string]func([]interface{}) interface{})
	}
	r.Inlines[name] = fn
}

func NewCodeRunner() *CodeRunner {
	var r = new(CodeRunner)
	r.Vars = make(map[string]interface{}, 0)
	r.Stack = list.NewStack()
	r.SetFunc("array", func(params []interface{}) interface{} {
		return params
	})
	r.SetFunc("len", func(params []interface{}) interface{} {
		if arr,ok := params[0].([]interface{});ok {
			return len(arr)
		}
		if mp,ok := params[0].(map[string]interface{});ok {
			return len(mp)
		}
		if str,ok := params[0].(string);ok {
			return len(str)
		}
		r.VmLogger.Panicf("len param must be array or map or string,but got %T",params[0])
		return nil
	})
	r.SetFunc("del",func(params []interface{}) interface{} {
		if arr,ok := params[0].([]interface{});ok {
			arr = append(arr[:params[1].(int)],arr[params[1].(int)+1:]...)
			return arr
		}
		if mp,ok := params[0].(map[string]interface{});ok {
			delete(mp,params[1].(string))
			return mp
		}
		r.VmLogger.Panicf("del param must be array or map")
		return nil
	})
	r.SetFunc("map", func(params []interface{}) interface{} {
		mp := make(map[string]interface{})
		for i:=0;i<len(params);i+=2{
			key :=  params[i]
			value := params[i+1]
			mp[cast.ToString(key)] = value
		}
		return mp
	})
	r.SetFunc("set", func(params []interface{}) interface{} {
		params[0].([]interface{})[cast.ToInt(params[1])] = params[2]
		return nil
	})
	r.SetFunc("get", func(params []interface{}) interface{} {
		return params[0].([]interface{})[cast.ToInt(params[1])]
	})
	r.SetFunc("add", func(params []interface{}) interface{} {
		f,ok :=params[0].([]interface{})
		if ok {
			f = append(f, params[1])
		}
		return f
	})
	r.SetFunc("pop", func(params []interface{}) interface{} {
		f,ok :=params[0].([]interface{})
		if ok {
			f = f[:len(f)-1]
		}
		return f
	})
	r.SetFunc("contains",func(params []interface{}) interface{} {
		w := slices.Contains(params[0].([]interface{}),params[1])
		return w
	})
	r.SetFunc("trimPrefix",func(params []interface{}) interface{} {
		return strings.TrimPrefix(params[0].(string),params[1].(string))
	})
	r.SetFunc("trimSuffix",func(params []interface{}) interface{} {
		return strings.TrimSuffix(params[0].(string),params[1].(string))
	})
	r.SetFunc("trim",func(params []interface{}) interface{} {
		return strings.Trim(params[0].(string),params[1].(string))
	})
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
			h.SetVar(k, h.evalNode(v.(ast.Anode)), isStackAssign)
		}

		return
	}

	set_var(h.Vars, h.Stack, k, v, isStackAssign)
}

func (c *CodeRunner) ParseAndRun(s string) int {
	ts := parseSourceTree(s)
	return c.RunCode(ts)
}

func (c *CodeRunner) ParseAndRunRecover(s string) (int,error) {
	ts := parseSourceTree(s)
	return c.RunCodeRecover(ts)
}

func (c *CodeRunner) RunCode(t ast.Node) int {
	if t == nil {
		return 0
	}
	c.evalNode(t)
	return int(c.ExitCode)
}

func (c *CodeRunner) RunCodeRecover(t ast.Node) (ex int,err error) {
	defer func() {
		if r := recover(); r != nil {
			c.VmLogger.Printf("panic: %v", r)
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	return c.RunCode(t),err
}