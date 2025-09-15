package evaluator

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"reflect"
	"strings"

	"github.com/lyr-2000/toylang/base/list"
	"github.com/spf13/cast"
)

// v2 evaluator

type Interpreter struct {
	// Code [][]string
	CodeReader io.Reader

	handler0 map[string]Handler0
	handler1 map[string]Handler1
	handler2 map[string]Handler2
	handler3 map[string]Handler3
	Nop      map[string]Nop
	stopDo   bool
	stack    *list.Stack
}

func (r *Interpreter) Top() *RefValue {
	return r.stack.Top().(*RefValue)
}
func (r *Interpreter) Pop() *RefValue {
	return r.stack.Pop().(*RefValue)
}
func (r *Interpreter) Push(v *RefValue) {
	v.Interpreter = r
	r.stack.Push(v)
}

func (r *Interpreter) SetReader(buf io.Reader) {
	r.CodeReader = buf
}
func (r *Interpreter) hasHandler(name string,k int) bool {
	switch k {
	case 0:
		_, ok := r.handler0[name]
		return ok
	case 1:
		_, ok := r.handler1[name]
		return ok
	case 2:
		_, ok := r.handler2[name]
		return ok
	case 3:
		_, ok := r.handler3[name]
		return ok
	default:
		_, ok := r.Nop[name]
		if ok {
			return ok
		}
	}
	return false
}
func (r *Interpreter) do(line string) {
	allf := strings.Fields(line)
	if len(allf) == 0 {
		return
	}
	name, args := strings.ToUpper(allf[0]), allf[1:]
	log.Println("cmd: ", name, len(args), args)
	if !r.hasHandler(name, len(args) + 1) {
		r.stopDo = true
		log.Panicf("invalid args :%s --%v", line,name)
	}
	switch len(args) + 1 {
	case 1:
		r.handler1[name](name)
	case 2:
		r.handler2[name](name, args[0])
	case 3:
		r.handler3[name](name, args[0], args[1])
	default:
		f, ok := r.Nop[name]
		if ok {
			f()
			return
		}
		r.stopDo = true
		log.Panicf("invalid args :" + line)
	}
}

func (r *Interpreter) Handle() {
	scanner := bufio.NewReader(r.CodeReader)
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
		r.do(line)
	}
}

type Handler0 = func() (any, error)
type Handler1 = func(op string) (any, error)
type Handler2 = func(op string, arg1 string) (any, error)
type Handler3 = func(op string, arg1, arg2 string) (any, error)

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
	case Nop:
		r.Nop[handlerName] = handler
	default:
		log.Panicf("invalid handler %s", handlerName)
	}
}

type RefValue struct {
	Interpreter *Interpreter
	getter      func(b *Interpreter) any
	setter      func(b *Interpreter, v any)
	Symbol      string
	Type        string
}

func NewVar(symbol string) *RefValue {
	return &RefValue{
		getter: func(b *Interpreter) any {
			return nil
		},
		setter: func(b *Interpreter, v any) {
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
func (r *RefValue) F() float64 {
	return cast.ToFloat64(r.getter(r.Interpreter))
}
func (r *RefValue) S() string {
	return cast.ToString(r.getter(r.Interpreter))
}

func (r *RefValue) Plus(b *RefValue) *RefValue {
	switch r.Type {
	case "NUMBER":
		return NewNumber(r.F() + b.F())
	}
	return NewString(r.S()+b.S(), false)
}

type Nop = func()

func New() *Interpreter {
	b := &Interpreter{
		handler0: make(map[string]Handler0),
		handler1: make(map[string]Handler1),
		handler2: make(map[string]Handler2),
		handler3: make(map[string]Handler3),
		Nop:      make(map[string]Nop),
		stack:    list.NewStack(),
	}
	globalSet(b)
	return b
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
	f.Set("OP", func(opname,plus string) (any, error) {
		switch plus {
		case "+":
			num1 := f.Pop()
			num2 := f.Pop()
			f.Push(num1.Plus(num2))
		case "-":
			num1 := f.stack.Pop()
			num2 := f.stack.Pop()
			f.stack.Push(NewNumber(num1.(float64) - num2.(float64)))
		case "*":
			num1 := f.stack.Pop()
			num2 := f.stack.Pop()
			f.stack.Push(NewNumber(num1.(float64) * num2.(float64)))
		case "/":
			num1 := f.stack.Pop()
			num2 := f.stack.Pop()
			f.stack.Push(NewNumber(num1.(float64) / num2.(float64)))
		default:
			return nil, ErrInvalidOp
		}
		return 0, fmt.Errorf("%w: %s", ErrInvalidOp, opname)
	})
}

var (
	ErrInvalidOp = fmt.Errorf("invalid op")
)
