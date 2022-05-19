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
func (h *CodeRunner) SetVar(k string, v interface{}, isStk bool) {
	if v == nil {
		set_var(h.Vars, h.Stack, k, v, isStk)
		return
	}
	if n, ok := v.(ast.Anode); ok {
		switch n.(type) {
		case *ast.Variable:
			varnode := n.(*ast.Variable)
			h.SetVar(k, h.GetVar(varnode.Lexeme.Value.(string)), isStk)
		case *ast.Scalar:
			varnode := n.(*ast.Scalar)
			switch varnode.Lexeme.Type {
			case lexer.Number:
				h.SetVar(k, cast.ToFloat64(varnode.Lexeme.Value), isStk)

			default:
				h.SetVar(k, varnode.Lexeme.Value.(string), isStk)
			}
			// h.SetVar(k, varnode, isStk)
		default:
			h.SetVar(k, h.evalNode(v.(ast.Anode)), isStk)
		}

		return
	}

	set_var(h.Vars, h.Stack, k, v, isStk)
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

	//run main
	// switch t.(type) {
	// case *ast.Factor, *ast.Scalar: // skip
	// //bad node
	// case *ast.Expr:

	// default:

	// }
	return 0
}
