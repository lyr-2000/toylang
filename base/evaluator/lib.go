package evaluator

import (
	"fmt"
	"strings"

	"github.com/lyr-2000/toylang/base/ast"
	"github.com/spf13/cast"
)

func IsLibFn(fn string) bool {
	return fn == "print" || fn == "echo" || fn == "printf" || fn == "println" || fn == "exit"
}

// 打印内置函数
func (h *CodeRunner) libPrint_(p []ast.Anode, ln bool) int {
	var input strings.Builder
	//第一个参数是 print, 后面的是参数
	for _, v := range p[1:] {
		a := h.evalNode(v)
		input.WriteString(fmt.Sprintf("%v", a))
	}
	if ln {
		fmt.Fprintln(h.vmOutput, input.String())
	} else {
		fmt.Fprint(h.vmOutput, input.String())
	}
	return 0
}
func (h *CodeRunner) libExit(p []ast.Anode) uint8 {
	for _, v := range p[1:] {
		a := h.evalNode(v)
		h.ExitCode = cast.ToUint8(a)
	}
	return h.ExitCode
}

func (h *CodeRunner) CallInline(fnName string, params []ast.Anode) (interface{}, bool) {
	if h.Inlines != nil {
		fn := h.Inlines[fnName]
		if fn != nil {
			var args []interface{}
			for _, v := range params[1:] {
				args = append(args, h.evalNode(v))
			}
			return fn(args), true
		}
	}
	return nil, false

}

func (h *CodeRunner) libFnCall(fnName string, paramNode []ast.Anode) interface{} {
	switch fnName {
	case "echo":
		return h.libPrint_(paramNode, false)
	case "print":
		return h.libPrint_(paramNode, true)
	case "println":

		return h.libPrint_(paramNode, true)
	case "exit":
		return h.libExit(paramNode)
	default:
	}
	return 0
}

func (h *CodeRunner) fn_define(n *ast.FuncStmt) {
	for _, v := range h.Functions {
		if v == n { //重复定义
			return
		}
	}
	h.Functions = append(h.Functions, n)
}
