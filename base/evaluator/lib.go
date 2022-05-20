package evaluator

import (
	"fmt"
	"strings"
	"toylang/base/ast"
)

func IsLibFn(fn string) bool {
	return fn == "print" || fn == "printf" || fn == "println"
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
		fmt.Println(input.String())
	} else {
		fmt.Print(input.String())
	}
	return 0
}

func (h *CodeRunner) libFnCall(fnName string, paramNode []ast.Anode) interface{} {
	switch fnName {
	case "print":
		return h.libPrint_(paramNode, false)
	case "println":

		return h.libPrint_(paramNode, true)
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
