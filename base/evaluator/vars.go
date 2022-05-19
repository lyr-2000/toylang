package evaluator

import "toylang/base/list"

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

*/
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
			mp[key] = value
		}
		return
	}
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
