package evaluator

import (
	"toylang/base/ast"
	"toylang/base/lexer"
	"toylang/base/list"

	"github.com/spf13/cast"
)

// var (
// 	global_variables map[string]interface{}
// 	stack            *list.Stack
// )

// func init() {
// 	stack = list.NewStack()
// 	global_variables = make(map[string]interface{}, 8)
// }

// func GetVar(key string) interface{} {
// 	return get_var(key)
// }
// func SetVar(key string, value interface{}, isStackVar bool) {
// 	set_var(key, value, isStackVar)
// }
func get_var(global_variables map[string]interface{}, stack *list.Stack, key string) interface{} {
	if stack.Len() <= 0 {
		return global_variables[key]
	}

	h := stack.Queue.Head
	for h != nil {
		if h.Value != nil {
			if mp, ok := h.Value.(map[string]interface{}); ok {
				res, ok := mp[key]
				if ok {
					return res
				}
			}
		}
		h = h.Next
	}
	return global_variables[key]
}

//这里的逻辑有点难处理，
/*
规定 使用 var 声明的，就在当前的栈顶开辟内存，否则在全局区域开辟内存

isStackVar 为true，就在当前堆栈查找 变量，有变量就覆盖，没有就修改全局变量
		   为 false， 优先 遍历链表，如果有变量 就覆盖，否则就修改全局变量

*/

func cast_scalar_node_type(a *ast.Scalar) interface{} {
	w := a.Lexeme.Value
	switch a.Lexeme.Type {
	case lexer.Boolean:
		return cast.ToBool(w)
	case lexer.Char:
		return cast.ToInt32(w)
	case lexer.Number:
		return cast.ToFloat64(w)
	case lexer.String:
		return w

	default:

	}
	return nil
}

func set_var(global_variables map[string]interface{}, stack *list.Stack, key string, value interface{}, isStackVar bool) {
	if isStackVar {
		h := stack.Queue.Head
		if h == nil {
			// head is null, then to global
			global_variables[key] = value
		} else {
			if h.Value == nil {
				h.Value = make(map[string]interface{}, 8)
			}
			mp := h.Value.(map[string]interface{})
			// var 关键字使用，会开辟新变量
			mp[key] = value
		}
		return
	}
	//优先修改变量
	// set to global
	h := stack.Queue.Head
	for h != nil {
		if h.Value == nil {
			h = h.Next
			continue
		}
		mp := h.Value.(map[string]interface{})
		//h.value is not null
		if _, ok := mp[key]; ok {
			mp[key] = value
			return
		}
		h = h.Next
	}
	//cannot set to stack
	global_variables[key] = value
}
