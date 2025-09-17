package evaluator

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"reflect"
	"strings"

	"github.com/spf13/cast"
)

// v2 evaluator

type Interpreter struct {
	// Code [][]string
	CodeReader io.Reader
	globalVar  map[string]any
	UnionStack *UnionStack
	Labels     map[string]int

	handler0 map[string]Handler0
	handler1 map[string]Handler1
	handler2 map[string]Handler2
	handler3 map[string]Handler3
	handler4 map[string]Handler4
	Nop      map[string]Nop
	stopDo   bool
	// stack        *list.Stack
	Program      [][]string
	ProgramIndex int
}

func (r *Interpreter) Top() *RefValue {
	return r.UnionStack.Top().Top().(*RefValue)
}

func (r *RefValue) CompareTo(w *RefValue) int {
	switch w.Type {
	case "NUMBER":
		if r.F() < w.F() {
			return -1
		}
		if r.F() > w.F() {
			return 1
		}
		return 0
	}
	if r.Any() == nil {
		if w.Any() == nil {
			return 0
		}
		return -1
	}
	if r.Any() == nil {
		if w.Any() == nil {
			return 0
		}
		return -1
	}
	return 0
}

func (r *Interpreter) Pop() *RefValue {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("panic: %v", err)
			r.PrintStack()
			panic(err)
		}
	}()
	b := r.UnionStack.Top()
	if b == nil {
		return nil
	}
	if b.Top() == nil {
		return nil
	}
	return b.Pop().(*RefValue)
}
func (r *Interpreter) Push(v *RefValue) {
	v.Interpreter = r
	r.UnionStack.Top().Push(v.Type, v)
}

func (r *Interpreter) SetReader(buf io.Reader) {
	r.CodeReader = buf
}
func (r *Interpreter) hasHandler(name string, k int) (any, bool) {
	var (
		f  any
		ok bool
	)
	switch k {
	case 0:
		f, ok = r.handler0[name]
	case 1:
		f, ok = r.handler1[name]
	case 2:
		f, ok = r.handler2[name]
	case 3:
		f, ok = r.handler3[name]
	case 4:
		f, ok = r.handler4[name]
	}
	if ok {
		return f, true
	}
	f, ok = r.Nop[name]
	if ok {
		return f, ok
	}
	return f, false
}

func (r *Interpreter) PrintStack() {
	for i := 0; i <= r.UnionStack.top; i++ {
		stack := r.UnionStack.Stack[i]
		var buf strings.Builder
		buf.WriteString(fmt.Sprintf("(%s) \n", stack.Label))
		for j := 0; j <= stack.top; j++ {
			buf.WriteString(fmt.Sprintf("\t %s\n", ParamsString(stack.Params[j])))
		}
		log.Printf("stack: %s\n", buf.String())
	}
}
func (r *Interpreter) do(allf []string) {
	if len(allf) == 0 {
		return
	}
	name, args := strings.ToUpper(allf[0]), allf[1:]
	log.Printf("cmd: %s,paramLen: %v, args: %v", name, len(args), args)
	h, ok := r.hasHandler(name, len(args)+1)
	if !ok {
		r.stopDo = true
		log.Panicf("invalid args :%s ", name)
	}
	var (
		out any
		err error
	)
	switch h := h.(type) {
	case Handler0:
		h()
	case Handler1:
		out, err = h(name)
	case Handler2:
		out, err = h(name, args[0])
	case Handler3:
		out, err = h(name, args[0], args[1])
	case Handler4:
		out, err = h(name, args[0], args[1], args[2])
	case Nop:
		h()
	default:
		log.Panicf("invalid handler %s ", name)
	}
	if err != nil {
		log.Panicf("error: %d %v %v ", r.ProgramIndex, err, out)
	}

}
func (r *Interpreter) Next() {
	r.ProgramIndex++
}
func (r *Interpreter) Index() int {
	return r.ProgramIndex
}

func (r *Interpreter) Handle() {
	scanner := bufio.NewReader(r.CodeReader)
	var all [][]string
	for !r.stopDo {
		line, err := scanner.ReadString('\n')
		if err != nil {
			break
		}
		if err == io.EOF {
			break
		}
		if line == "" {
			continue
		}
		b := strings.Fields(line)
		all = append(all, b)
	}
	r.Program = all
	r.ProgramIndex = 0
	for r.ProgramIndex < len(r.Program) {
		r.do(r.Program[r.ProgramIndex])
		r.ProgramIndex++
	}
}

type Handler0 = func() (any, error)
type Handler1 = func(op string) (any, error)
type Handler2 = func(op string, arg1 string) (any, error)
type Handler3 = func(op string, arg1, arg2 string) (any, error)
type Handler4 = func(op string, arg1, arg2, arg3 string) (any, error)

func (r *Interpreter) Set(handlerName string, handler any) {
	log.Println("set", handlerName, reflect.TypeOf(handler))
	handlerName = strings.ToUpper(handlerName)
	switch handler := handler.(type) {
	case Handler0:
		r.handler0[handlerName] = handler
	case Handler1:
		r.handler1[handlerName] = handler
	case Handler2:
		r.handler2[handlerName] = handler
	case Handler3:
		r.handler3[handlerName] = handler
	case Handler4:
		r.handler4[handlerName] = handler
	case Nop:
		r.Nop[handlerName] = handler
	default:
		log.Panicf("invalid handler %s ,onlyFound:  %T", handlerName, handler)
	}
}

type RefValue struct {
	Interpreter *Interpreter
	getter      func(b *Interpreter) any
	setter      func(b *Interpreter, v any)
	Symbol      string
	Type        string
	Value       any
}

func (r *RefValue) String() string {
	return fmt.Sprintf("%s %v", r.Symbol, r.Any())
}

func NewVar(symbol string) *RefValue {
	d := &RefValue{
		getter: func(b *Interpreter) any {
			return b.globalVar[symbol]
		},
		Symbol: symbol,
		Type:   "VAR",
	}
	d.setter = func(b *Interpreter, v any) {
		b.globalVar[symbol] = v
		d.Value = v
	}
	return d
}

func NewBool(num bool) *RefValue {
	return &RefValue{
		getter: func(b *Interpreter) any {
			return num
		},
		setter: func(b *Interpreter, v any) {
			num = cast.ToBool(v)
		},
		Symbol: fmt.Sprintf("%v", num),
		Type:   "BOOL",
	}
}
func NewNumber(num float64) *RefValue {
	return &RefValue{
		getter: func(b *Interpreter) any {
			return num
		},
		setter: func(b *Interpreter, v any) {
			num = v.(float64)
		},
		Symbol: fmt.Sprintf("%v", num),
		Type:   "NUMBER",
	}
}

func NewString(str string, b64 bool) *RefValue {
	if b64 {
		str0, err := base64.StdEncoding.DecodeString(str)
		if err != nil {
			panic(err)
		}
		str = string(str0)
	}
	return &RefValue{
		getter: func(b *Interpreter) any {
			return str
		},
	}
}
func (r *RefValue) SysName() string {
	switch r.Type {
	case "NUMBER":
		return fmt.Sprintf("%v", r.F())
	case "STRING":
		return fmt.Sprintf("%s", r.Str())
	case "VAR":
		return fmt.Sprintf("%s (%v)", r.Symbol, r.Any())
	}
	return ""
}
func (r *RefValue) Bool() bool {
	w := r.Any()
	switch e := w.(type) {
	case bool:
		return e
	case float64:
		return e != 0
	case int, int64, int32, int16, int8:
		return e != 0
	case string:
		return e != ""
	case nil:
		return false
	}
	return false
}
func (r *RefValue) Any() any {
	e := r.getter(r.Interpreter)
	return e
}
func (r *RefValue) F() float64 {
	return cast.ToFloat64(r.getter(r.Interpreter))
}
func (r *RefValue) Str() string {
	return cast.ToString(r.getter(r.Interpreter))
}

func (r *RefValue) Plus(b *RefValue) *RefValue {
	switch r.Type {
	case "NUMBER":
		return NewNumber(r.F() + b.F())
	}
	return NewString(r.Str()+b.Str(), false)
}

type Nop = func()

func New() *Interpreter {
	sb := &UnionStack{
		Stack:    make([]*LabelStack, 0),
		top:      -1,
		AssignId: 0,
	}
	sb.Push("GLOBAL", []any{})
	b := &Interpreter{
		handler0: make(map[string]Handler0),
		handler1: make(map[string]Handler1),
		handler2: make(map[string]Handler2),
		handler3: make(map[string]Handler3),
		handler4: make(map[string]Handler4),
		Nop:      make(map[string]Nop),
		// stack:     list.NewStack(),
		globalVar:  make(map[string]any),
		UnionStack: sb,
		Labels:     make(map[string]int),
	}
	globalSet(b)
	return b
}

func (r *Interpreter) LookupVar(name string) *RefValue {
	for i := r.UnionStack.top; i >= 0; i-- {
		stack := r.UnionStack.Stack[i]
		for j := stack.Index(); j >= 0; j-- {
			if stack.ParamsName[j] == "VAR" {
				el := stack.Params[j].(*RefValue)
				if el.Symbol == name {
					return stack.Params[j].(*RefValue)
				}
			}
		}
	}
	return nil
}

func (r *Interpreter) GetGlobalVar(name string) any {
	return r.globalVar[name]
}
func (r *Interpreter) SetGlobalVar(name string, v any) {
	log.Println("set global var", name, v)
	r.globalVar[name] = v
}

// TODO: check fast
func (f *Interpreter) GOTO(index int, prefix, arg string) {
	for i := index; i < len(f.Program); i++ {
		line := f.Program[i]
		if len(line) == 0 {
			continue
		}
		if strings.EqualFold(line[0], prefix) {
			if len(line) > 1 {
				if line[1] == arg {
					log.Printf("goto %s %s", prefix, arg)
					f.ProgramIndex = i - 1
					return
				}
			}
		}
	}
	log.Panicf("!cmd error: goto %s %s not found", prefix, arg)
}

func globalSet(f *Interpreter) {
	if len(f.Nop) > 0 || len(f.handler2) > 0 {
		return
	}
	f.Set("VAR", func(op string, arg string) (any, error) {
		f.Push(NewVar(arg))
		return arg, nil
	})
	f.Set("forStmtBegin", func(_, label string, nodeCount string) (any, error) {
		f.Labels[label] = f.Index()
		return nil, nil
	})
	f.Set("forStmtEnd", func(_, label string) (any, error) {
		delete(f.Labels, label)
		return nil, nil
	})
	f.Set("ExitFor", func(_, label string) (any, error) {
		f.GOTO(f.Labels[label], "forStmtEnd", label)
		return nil, nil
	})
	// continue keyword
	f.Set("continue", func(_, label string) (any, error) {
		f.GOTO(f.Labels[label], "for_continue", label)
		return nil, nil
	})
	// is condition stmt
	f.Set("for_continue", func(_, label string) (any, error) {
		return nil, nil
	})
	f.Set("for_continue_end", func(_, label string, _ string) (any, error) {
		d := f.Pop()
		if !d.Bool() {
			f.GOTO(f.Index(), "ExitFor", label)
			return false, nil
		}
		return "", nil
	})
	f.Set("for_body", func(_, op string) (any, error) {
		return nil, nil
	})
	f.Set("for_body_end", func(_, op string) (any, error) {
		return nil, nil
	})

	// GOTO #ENDIF L6
	f.Set("GOTO", func(op, syslabel, label2 string) (any, error) {
		f.GOTO(f.Index(), syslabel, label2)
		return nil, nil
	})
	f.Set("if", func() {})
	f.Set("#IF", func() {})
	f.Set("#ENDIF", func() {})

	f.Set("endif", func(op, label string) (any, error) {
		b := f.Pop().Bool()
		if !b {
			f.GOTO(f.Index(), "#ENDIF", label)
			return nil, nil
		}
		return nil, nil
	})
	f.Set("endif", func(op, label, cond, label2 string) (any, error) {
		b := f.Pop().Bool()
		if b {
			return nil, nil
		}
		f.GOTO(f.Index(), "#ELSEIF", label2)
		return nil, nil
	})
	f.Set("#ELSEIF", func() {})
	f.Set("ifbodystart", func(op, label string) (any, error) {

		return nil, nil
	})
	f.Set("ifbodyend", func() {})

	f.Set("blockstart", func() {})
	f.Set("blockend", func() {})
	f.Set("expr", func(op, tokenType, val string) (any, error) {
		switch tokenType {
		case "number":
			f.Push(NewNumber(cast.ToFloat64(val)))
		case "str":
			f.Push(NewString(val, false))
		case "variable":
			f.Push(NewVar(val))
		case "strb64":
			f.Push(NewString(val, true))
		default:
			return val, fmt.Errorf("invalid token type: %s", tokenType)
		}
		return val, nil
	})
	f.Set("for_init", func(_, op string) (any, error) {

		return nil, nil
	})
	f.Set("for_init_end", func(_, op string) (any, error) {
		return nil, nil
	})
	f.Set("for_step", func(_, op string) (any, error) {
		f.UnionStack.Push("for_step", []any{})
		return nil, nil
	})
	f.Set("for_step_end", func(_, op string) (any, error) {
		f.UnionStack.Pop()
		return nil, nil
	})

	f.Set("blockstart", func() {})
	f.Set("printstack", func() {
		f.PrintStack()
	})

	f.Set("var", func(op, _, arg string) (any, error) {
		// f.UnionStack.SetVarAtTop(arg, nil)
		f.Push(NewVar(arg))
		return nil, nil
	})
	f.Set("declare", func(op, arg string) (any, error) {
		// f.UnionStack.SetVarAtTop(arg, nil)
		f.Push(NewVar(arg))
		return nil, nil
	})
	f.Set("call_arg", func(op, label, fnName, cnt string) (any, error) {
		// call_arg L3 print 1
		return nil, nil
	})
	f.Set("CALL", func(op, label, Name, cnt string) (any, error) {
		N := cast.ToInt(cnt)
		rev := make([]any, N)
		for k := N - 1; k >= 0; k-- {
			ele := f.Pop()
			rev[k] = ele.Any()
		}
		log.Printf("FUNCTIONCALL %s %s %s %v", op, Name, cnt, Json(rev))
		return nil, nil
	})

	f.Set("OP", f.OP)
	f.Set("ASSIGN", func(op, arg string) (any, error) {
		num1 := f.Pop()
		variableVal := f.Pop()
		variableVal.setter(f, num1.Any())
		// f.UnionStack.SetVarAtTop(num2.SysName(), num1)
		f.Push(variableVal)
		return nil, nil
	})
}

var (
	ErrInvalidOp = fmt.Errorf("invalid op")
)

func (f *Interpreter) OP(op, plus string) (any, error) {
	switch plus {
	case "+":
		num1 := f.Pop()
		num2 := f.Pop()
		f.Push(NewNumber(num1.F() + num2.F()))
	case "-":
		num1 := f.Pop()
		num2 := f.Pop()
		f.Push(NewNumber(num1.F() - num2.F()))
	case "*":
		num1 := f.Pop()
		num2 := f.Pop()
		f.Push(NewNumber(num1.F() * num2.F()))
	case "/":
		num1 := f.Pop()
		num2 := f.Pop()
		f.Push(NewNumber(num1.F() / num2.F()))
	case "=":
		num1 := f.Pop()
		variableVal := f.Pop()
		variableVal.setter(f, num1.Any())
		if e := f.LookupVar(variableVal.Symbol); e != nil {
			e.setter(f, variableVal.Any())
			return nil, nil
		}
		f.Push(variableVal)
	case "<":
		num1 := f.Pop()
		vari := f.Pop()
		f.Push(NewBool(vari.CompareTo(num1) < 0))
		return f.Top(), nil
	case "++":
		num1 := f.Pop()
		num1.setter(f, num1.F()+1)
		f.Push(num1)
	case ">":
		num1 := f.Pop()
		vari := f.Pop()
		f.Push(NewBool(vari.CompareTo(num1) > 0))
		return f.Top(), nil
	case "<=":
		num1 := f.Pop()
		vari := f.Pop()
		f.Push(NewBool(vari.CompareTo(num1) <= 0))
		return f.Top(), nil
	case ">=":
		num1 := f.Pop()
		vari := f.Pop()
		f.Push(NewBool(vari.CompareTo(num1) >= 0))
		return f.Top(), nil
	case "!=":
		num1 := f.Pop()
		vari := f.Pop()
		f.Push(NewBool(vari.CompareTo(num1) != 0))
		return f.Top(), nil
	case "==":
		num1 := f.Pop()
		numvar := f.Pop()
		f.Push(NewBool(num1.Any() == numvar.Any()))
		return f.Top(), nil
	default:
		log.Printf("invalid op: %s", plus)
		return nil, ErrInvalidOp
	}
	return f.Top().Any(), nil
}

func Json(v any) string {
	bs, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(bs)
}
