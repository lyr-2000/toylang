package evaluator

import (
	"toylang/base/ast"
	"toylang/base/lexer"
	"toylang/base/list"

	"github.com/spf13/cast"
)

type CodeRunner struct {
	Functions []*ast.FuncStmt
	Vars      map[string]interface{}
	Stack     *list.Stack

	stackDeep int
}

func NewCodeRunner() *CodeRunner {
	var r = new(CodeRunner)
	r.Vars = make(map[string]interface{}, 0)
	r.Stack = list.NewStack()
	return r
}
func (h *CodeRunner) GetVar(key string) interface{} {
	return get_var(h.Vars, h.Stack, key)
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
	ts := parse_source_tree(s)
	return c.RunCode(ts)
}
func (c *CodeRunner) RunCode(t ast.Node) int {
	if t == nil {
		return 0
	}
	c.evalNode(t)

	return 0
}
