package evaluator

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"

	// "reflect"
	"strings"

	"github.com/spf13/cast"
)

// v2 evaluator

type Interpreter struct {
	Stdout io.Writer
	Logger *log.Logger
	Debug  bool
	// Code [][]string
	CodeReader io.Reader
	globalVar  map[string]any
	UnionStack *UnionStack
	Labels     map[string]int

	handler0       map[string]Handler0
	handler1       map[string]Handler1
	handler2       map[string]Handler2
	handler3       map[string]Handler3
	handler4       map[string]Handler4
	Nop            map[string]Nop
	stopDo         bool
	Program        [][]string
	ProgramIndex   int
	funcAt         map[string][]int
	funcParamCount map[string]int

	ErrCode  int
	ExitCode int
	ErrMsg   string
}

func (r *Interpreter) scanFuncLine() {
	if r.funcAt == nil {
		r.funcAt = make(map[string][]int)
	}
	if r.funcParamCount == nil {
		r.funcParamCount = make(map[string]int)
	}
	for i := 0; i < len(r.Program); i++ {
		line := r.Program[i]
		if strings.EqualFold(line[0], "funcStmtBegin") {
			at := i
			label := line[1]
			fnName := line[2]
			paramCount := cast.ToInt(line[3])
			r.funcParamCount[fnName] = paramCount
			for j := i; j < len(r.Program); j++ {
				if strings.EqualFold(r.Program[j][0], "funcStmtEnd") {
					label2 := r.Program[j][1]
					if label2 == label {
						r.funcAt[fnName] = []int{at, j}
						break
					}
				}
			}
			if r.funcAt[fnName] == nil {
				log.Panicf("func %s not found", fnName)

			}
		}

	}
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
	if r.Debug {
		log.Printf("cmd: %s,paramLen: %v, args: %v", name, len(args), args)
	}
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
	r.scanFuncLine()
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
	// log.Println("set", handlerName, reflect.TypeOf(handler))
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

func (r *RefValue) ReplaceAsVar(varName string) {
	r.Symbol = varName
	r.Type = "VAR"
	r.Value = r.getter(r.Interpreter)
}

func (r *RefValue) String() string {
	return fmt.Sprintf("%s %v", r.Symbol, r.Any())
}

func NewVar2(symbol string, f any) *RefValue {
	d := &RefValue{
		Symbol: symbol,
		Type:   "VAR",
		Value:  f,
	}
	d.setter = func(b *Interpreter, v any) {
		d.Value = v
	}
	d.getter = func(b *Interpreter) any {
		return d.Value
	}
	return d
}

func NewVar(symbol string) *RefValue {
	d := &RefValue{
		Symbol: symbol,
		Type:   "VAR",
	}
	d.setter = func(b *Interpreter, v any) {
		d.Value = v
	}
	d.getter = func(b *Interpreter) any {
		return d.Value
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
func (r *RefValue) I() int {
	return cast.ToInt(r.getter(r.Interpreter))
}
func (r *RefValue) Str() string {
	el := r.getter(r.Interpreter)
	d := reflect.ValueOf(el)
	switch d.Kind() {
	case reflect.Map, reflect.Struct:
		return fmt.Sprintf("%v", el)
	}
	return cast.ToString(el)
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
		Stdout:     os.Stdout,
		Logger:     log.New(os.Stdout, "SYS", log.LstdFlags),
		handler0:   make(map[string]Handler0),
		handler1:   make(map[string]Handler1),
		handler2:   make(map[string]Handler2),
		handler3:   make(map[string]Handler3),
		handler4:   make(map[string]Handler4),
		Nop:        make(map[string]Nop),
		globalVar:  make(map[string]any),
		UnionStack: sb,
		Labels:     make(map[string]int),
	}
	globalSet(b)
	return b
}

func (r *Interpreter) LookupRevIndex(idx int) *RefValue {
	for i := r.UnionStack.top; i >= 0; i-- {
		stack := r.UnionStack.Stack[i]
		for j := stack.Index(); j >= 0; j-- {
			if idx == 0 {
				return stack.Params[j].(*RefValue)
			}
			idx--
		}
	}
	return nil
}

func (r *Interpreter) LookupVar(name string) *RefValue {
	for i := r.UnionStack.top; i >= 0; i-- {
		stack := r.UnionStack.Stack[i]
		for j := stack.Index(); j >= 0; j-- {
			el := stack.Params[j].(*RefValue)
			if el.Type == "VAR" && el.Symbol == name {
				return stack.Params[j].(*RefValue)
			}
		}
	}
	return nil
}

func (r *Interpreter) GetGlobalVar(name string) any {
	return r.globalVar[name]
}
func (r *Interpreter) SetGlobalVar(name string, v any) {
	// log.Println("set global var", name, v)
	r.globalVar[name] = v
}

// TODO: check fast
func (f *Interpreter) GOTO(index int, prefix, arg string) bool {
	for i := index; i < len(f.Program); i++ {
		line := f.Program[i]
		if len(line) == 0 {
			continue
		}
		if strings.EqualFold(line[0], prefix) {
			if len(line) > 1 {
				if line[1] == arg {
					// log.Printf("goto %s %s", prefix, arg)
					f.ProgramIndex = i - 1
					return true
				}
			}
		}
	}

	// log.Panicf("!cmd error: goto %s %s not found", prefix, arg)
	return false
}

func globalSet(f *Interpreter) {
	f.scanFuncLine()
	if len(f.Nop) > 0 || len(f.handler2) > 0 {
		return
	}
	f.Set("@RETURNBEGIN", func() {})
	f.Set("@RETURN", func(op string, label string) (any, error) {
		p := f.UnionStack.Top()
		return popFn(f, p.Label)
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
		if !f.GOTO(f.Index(), syslabel, label2) {
			if !f.GOTO(0, syslabel, label2) {
				log.Panicf("!cmd error: goto %s %s not found", syslabel, label2)
			}
		}
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
		case "str", "string":
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
		p := f.LookupVar(arg)
		if p != nil {
			f.Push(p)
		} else {
			f.Push(NewVar(arg))
		}
		return f.Top(), nil
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
		rev := make([]*RefValue, N)
		for k := N - 1; k >= 0; k-- {
			ele := f.Pop()
			rev[k] = ele
		}
		if f.Debug {
			log.Printf("FUNCTIONCALL %s %s %s %v", op, Name, cnt, Json(rev))
		}
		if fn, ok := inlineFunc[Name]; ok {
			d := fn(f, rev...)
			f.Push(NewVar2("@return", d))
			return nil, nil
		}
		_, ok := f.funcAt[Name]
		if !ok {
			log.Panicf("function %s not found", Name)
		}
		return langFuncCall(f, Name, rev...)
	})
	f.Set("fn_arg", func(_, varType, varName, cnt string) (any, error) {
		idx := cast.ToInt(cnt)
		ele := f.LookupRevIndex(idx)
		ele.ReplaceAsVar(varName)
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
	f.Set("funcStmtBegin", func(op, label, funcName, cnt string) (any, error) {
		f.skipFuncLine(funcName)
		return nil, nil
	})
	f.Set("funcEntry", func(op, funcName, label, cnt string) (any, error) {
		return nil, nil
	})
	f.Set("FN_ARG_COUNT", func(op, cnt string) (any, error) {
		return nil, nil
	})
	f.Set("funcStmtEnd", func(op, _, funcName string) (any, error) {
		return popFn(f, funcName)
	})
	f.Set("@getvalue", func(op, varName, varType, keyvalue string) (any, error) {
		d := getArrayOrMap(f, varName, varType, keyvalue)
		f.Push(NewVar2("@getVal", d))
		return nil, nil
	})
}

func getArrayOrMap(b *Interpreter, varName, varType, keyvalue string) any {
	d := b.LookupVar(varName)

	switch varType {
	case "number":
		bb, ok := d.Any().([]any)
		idx := cast.ToInt(keyvalue)
		if ok && bb != nil && idx >= 0 && idx < len(bb) {
			return bb[cast.ToInt(keyvalue)]
		}
		fallthrough
	default:
		bb, ok := d.Any().(map[string]any)
		if ok && bb != nil {
			return bb[keyvalue]
		}
		return nil
	}

}

func popFn(f *Interpreter, fname string) (any, error) {
	ln := f.UnionStack.Top()
	f.UnionStack.Pop()
	// ln.Line

	f.ProgramIndex = ln.Line
	if ln.Top() != nil {
		f.UnionStack.Top().Push("@0", ln.Top())
	} else {
		f.UnionStack.Top().Push("@0", NewNumber(0))
	}
	return f.UnionStack.Top(), nil
}

var (
	ErrInvalidOp = fmt.Errorf("invalid op")
)

func (f *Interpreter) skipFuncLine(funcName string) {
	at, ok := f.funcAt[funcName]
	if !ok {
		log.Panicf("func %s not found", funcName)
	}
	f.ProgramIndex = at[1]
}

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

type InlineFunc func(f *Interpreter, args ...*RefValue) any

var (
	inlineFunc = map[string]InlineFunc{
		"print":   printLn,
		"throw":   throwErr,
		"recover": recoverErr,
		"fatal":   fatalErr,
		"array":   NewArray,
		"map":     NewMap,
		"set":     exitCodeSet,
		"exit":    exitCodeSet,
	}
)

func exitCodeSet(f *Interpreter, args ...*RefValue) any {
	varObj := args[0]
	el := varObj.I()
	f.ExitCode = el
	return f.ExitCode
}
func setArrayOrMap(f *Interpreter, args ...*RefValue) any {
	varObj := args[0]
	key := args[1].Str()
	value := args[2].Any()
	el := varObj.Any()
	switch arr := el.(type) {
	case []any:
		idx := cast.ToInt(key)
		if idx >= 0 && idx < len(arr) {
			arr[idx] = value
		}
	case map[string]any:
		arr[key] = value
	default:
	}

	return nil
}

func NewMap(f *Interpreter, args ...*RefValue) any {
	d := make(map[string]any)
	for i := 0; i+1 < len(args); i += 2 {
		d[args[i].Str()] = args[i+1].Any()
	}
	return d
}
func NewArray(f *Interpreter, args ...*RefValue) any {
	d := make([]any, len(args))
	for i, v := range args {
		d[i] = v.Any()
	}
	return d
}

func fatalErr(f *Interpreter, args ...*RefValue) any {
	var buf strings.Builder
	for _, v := range args {
		buf.WriteString(v.Str())
	}
	f.Logger.Fatal(buf.String())
	return nil
}
func recoverErr(f *Interpreter, args ...*RefValue) any {
	if f.ErrCode == 0 {
		return nil
	}
	d := make(map[string]any)
	d["code"] = f.ErrCode
	d["msg"] = f.ErrMsg
	f.ErrCode = 0
	f.ErrMsg = ""
	return d
}
func throwErr(f *Interpreter, args ...*RefValue) any {
	f.ErrCode = args[0].I()
	f.ErrMsg = args[0].Str()
	return nil
}

func printLn(f *Interpreter, args ...*RefValue) any {
	var buf strings.Builder
	for _, v := range args {
		buf.WriteString(v.Str())
	}
	fmt.Fprintln(f.Stdout, buf.String())
	f.Push(NewNumber(1))
	return nil
}

func langFuncCall(f *Interpreter, fnName string, args ...*RefValue) (any, error) {
	f.UnionStack.Push(fnName, []any{})
	top := f.UnionStack.Top()
	top.Line = f.Index()
	// should push stack
	// todo: push param stack
	// todo: defer pop param stack
	// f.Push(NewNumber(1))
	cnt := f.funcParamCount[fnName]
	if cnt < len(args) {
		log.Panicf("func %s param count mismatch, expected %d, got %d", fnName, cnt, len(args))
	}
	nilCnt := cnt - len(args)
	for i := 0; i < nilCnt; i++ {
		f.Push(NewNumber(0))
	}
	for i := len(args) - 1; i >= 0; i-- {
		f.Push(args[i])
	}
	d := f.funcAt[fnName]

	//TODO:  program call
	f.GOTO(d[0], "funcEntry", fnName)
	return nil, nil
}
