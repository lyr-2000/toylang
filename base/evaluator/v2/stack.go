package evaluator

import (
	"fmt"
	"log"
	"strings"
)

type LabelStack struct {
	Label     string
	StackType byte ///function or block ('F',or 'B')
	Params    []*RefValue
	top       int
	Line      int
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
	FuncStackCount uint16
	// todo count function stack size
}

func (r *UnionStack) pushStackx(label string, params []*RefValue, stackType byte) {
	if stackType == 'F' {
		r.FuncStackCount++
	}
	if r.top+1 >= len(r.Stack) {
		r.Stack = append(r.Stack, &LabelStack{
			Label:     label,
			Params:    params,
			top:       len(params) - 1,
			StackType: stackType,
		})
		r.top++
		return
	}

	r.Stack[r.top+1] = &LabelStack{
		Label:     label,
		Params:    params,
		top:       len(params) - 1,
		StackType: stackType,
	}
	r.top++
}

func (r *UnionStack) FreeUnused() {
	for i := r.top + 1; i < len(r.Stack); i++ {
		r.Stack[i].FreeUnused()
	}
	r.Stack = r.Stack[:r.top+1]
}

// function stack
func (r *UnionStack) PushF(label string, params []*RefValue) {
	r.pushStackx(label, params,'F')
}

// block stack
func (r *UnionStack) PushB(label string, params []*RefValue) {
	r.pushStackx(label, params,'B')
}

func (r *UnionStack) Pop() *LabelStack {
	if r.top < 0 || r.top >= len(r.Stack) {
		log.Panicf("stack pop error")
		return nil
	}

	d := r.Stack[r.top]
	r.top--
	if d.StackType == 'F' {
		r.FuncStackCount--
	}
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
	for i := r.top + 1; i < len(r.Params); i++ {
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
