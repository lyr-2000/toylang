package evaluator

import (
	"bufio"
	"encoding/base64"
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
func (r *Interpreter) Pop() *RefValue {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("panic: %v", err)
			r.PrintStack()
			panic(err)
		}
	}()
	return r.UnionStack.Top().Pop().(*RefValue)
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
	switch h := h.(type) {
	case Handler0:
		h()
	case Handler1:
		h(name)
	case Handler2:
		h(name, args[0])
	case Handler3:
		h(name, args[0], args[1])
	case Handler4:
		h(name, args[0], args[1], args[2])
	case Nop:
		h()
	default:
		log.Panicf("invalid handler %s ", name)
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
}

func (r *RefValue) String() string {
	return fmt.Sprintf("%s %v", r.Symbol, r.Any())
}

func NewVar(symbol string) *RefValue {
	return &RefValue{
		getter: func(b *Interpreter) any {
			return b.globalVar[symbol]
		},
		setter: func(b *Interpreter, v any) {
			b.globalVar[symbol] = v
		},
		Symbol: symbol,
		Type:   "VAR",
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

func (r *RefValue) Any() any {
	return r.getter(r.Interpreter)
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
		top:      0,
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
	}
	globalSet(b)
	return b
}
func (r *Interpreter) GetGlobalVar(name string) any {
	return r.globalVar[name]
}
func (r *Interpreter) SetGlobalVar(name string, v any) {
	log.Println("set global var", name, v)
	r.globalVar[name] = v
}

func globalSet(f *Interpreter) {
	if len(f.handler0) > 0 {
		return
	}
	f.Set("VAR", func(op string, arg string) (any, error) {
		f.Push(NewVar(arg))
		return arg, nil
	})
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
		}
		return val, fmt.Errorf("invalid token type: %s", tokenType)
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
		f.UnionStack.Push(label, []any{fnName, cnt})
		return nil, nil
	})
	f.Set("CALL", func(op, label, Name string) (any, error) {
		log.Printf("call %s %s", op, Name)
		return nil, nil
	})
	f.Set("OP", func(opname, plus string) (any, error) {
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
			f.PrintStack()
			num1 := f.Pop()
			variableVal := f.Pop()
			variableVal.setter(f, num1.Any())
			// f.UnionStack.SetVarAtTop(num2.SysName(), num1)
			f.Push(variableVal)
		default:
			return nil, ErrInvalidOp
		}
		return f.Top().Any(), nil
	})
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
