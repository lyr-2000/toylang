package evaluator

import (
	"fmt"
	"log"
	"strings"
)

type LabelStack struct {
	Label      string
	// ParamsName []string
	Params     []*RefValue
	top        int
	Line int
	// ReturnValue any
}

func ParamsString(params ...any) string {
	var buf strings.Builder
	for i := range params {
		switch w := params[i].(type) {
		case *RefValue:
			buf.WriteString(fmt.Sprintf("%s %v", w.Type, w.SysName()))
		default:
			buf.WriteString(fmt.Sprintf("%v ", params[i]))
		}
	}
	return buf.String()
}

type UnionStack struct {
	Stack    []*LabelStack
	top      int
	AssignId uint64
}

// func (r *UnionStack) SetGlobalVar(name string, v any) {
// 	if name == "" {
// 		r.AssignId++
// 		name = fmt.Sprintf("@%d",r.AssignId)
// 	}
// 	r.Stack[0].Push(name, v)
// }

// func (r *UnionStack) SetVarAtTop(name string, v any) {
// 	if name == "" {
// 		r.AssignId++
// 		name = fmt.Sprintf("@%d",r.AssignId)
// 	}
// 	r.Stack[r.top].Push(name, v)
// }



func (r *UnionStack) PushWithName(label string, params []*RefValue, names []string) {
	if len(params) != len(names) {
		log.Panicf("params and names length mismatch")
	}
	if r.top+1 >= len(r.Stack) {
		r.Stack = append(r.Stack, &LabelStack{
			Label:      label,
			Params:     params,
			// ParamsName: names,
			top:        len(params) - 1,
		})
		r.top++
		return
	}
	r.Stack[r.top+1] = &LabelStack{
		Label:      label,
		Params:     params,
		// ParamsName: names,
		top:        len(params) - 1,
	}
	r.top++
}

func (r *UnionStack) FreeUnused() {
	for i:=r.top+1;i<len(r.Stack);i++ {
		r.Stack[i].FreeUnused()
	}
	r.Stack = r.Stack[:r.top+1]
}

func (r *UnionStack) Push(label string, params []*RefValue) {
	r.PushWithName(label, params, make([]string, len(params)))
}

func (r *UnionStack) Pop() *LabelStack {
	if r.top < 0 || r.top >= len(r.Stack) {
		log.Panicf("stack pop error")
		return nil
	}

	d := r.Stack[r.top]
	r.top--
	return d
}

func (r *UnionStack) Size() int {
	return r.top + 1
}
func (r *UnionStack) Top() *LabelStack {
	return r.Stack[r.top]
}
func (r *UnionStack) Cap() int {
	return len(r.Stack)
}

func (r *LabelStack) Index() int {
	return r.top
}

func (r *LabelStack) StackLen() int {
	return len(r.Params)
}

func (r *LabelStack) Pop() any {
	d := r.Params[r.top]
	r.top--
	return d
}

func (r *LabelStack) Push(name string, v *RefValue) {
	if r.top+1 >= len(r.Params) {
		r.Params = append(r.Params, v)
		// r.ParamsName = append(r.ParamsName, name)
	}
	//Final PushTOP

	r.Params[r.top+1] = v
	// r.ParamsName[r.top+1] = name
	r.top++
}


func (r *LabelStack) FreeUnused() {
	if r == nil {
		return
	}
	for i:=r.top+1;i<len(r.Params);i++ {
		if r.Params[i] != nil {
			r.Params[i].Free()
		}
		r.Params[i] = nil
	}
	r.Params = r.Params[:r.top+1]
}


func (r *LabelStack) Top() *RefValue {
	if r.top < 0 || r.top >= len(r.Params) {
		return nil
	}
	return r.Params[r.top]
}

func (r *LabelStack) Len() int {
	return r.top + 1
}
