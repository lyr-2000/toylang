package evaluator

import (
	"fmt"
	"strings"

	"github.com/lyr-2000/toylang/base/ast"
	"github.com/spf13/cast"
)

func IsLibFn(fn string) bool {
	return fn == "print" || fn == "printf" || fn == "println" || fn == "exit"
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
func (h *CodeRunner) libExit(p []ast.Anode) uint8 {
	for _, v := range p[1:]{
		a := h.evalNode(v)
		h.ExitCode = cast.ToUint8(a)
	}
	return h.ExitCode
}

func (h *CodeRunner) libFnCall(fnName string, paramNode []ast.Anode) interface{} {
	switch fnName {
	case "print":
		return h.libPrint_(paramNode, false)
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
