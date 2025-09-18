package ast

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/spf13/cast"
)

var (
	Debug = true
)

type Writer struct {
	io.Writer
	no     uint64
	indent uint16
	CtxForLabel string
}

func (r *Writer) Indent() string {
	return strings.Repeat(" ", int(r.indent)*2)
}
func (r *Writer) AddI(i int) int {
	if !Debug {
		return i
	}
	r.indent += uint16(i)
	return int(r.indent)
}

func (r *Writer) L() uint64 {
	r.no++
	return r.no
}
func (r *Writer) Ln() string {
	return fmt.Sprintf("L%d", r.L())
}

func write(w *Writer, lines []string) {
	w.Write([]byte(w.Indent() + strings.Join(lines, " ")))
}

func writeln(w *Writer, lines []string) {
	w.Write([]byte(w.Indent() + strings.Join(lines, " ")))
	w.Write([]byte("\n"))
}

func jsons(w any) string {
	json, err := json.Marshal(w)
	if err != nil {
		return ""
	}
	return string(json)
}
func (e *Expr) Output(w *Writer) {
	if e == nil {
		return
	}
	for _, v := range e.GetChildren() {
		v.Output(w)
	}
	if e.GetLexeme() != nil {
		writeln(w, []string{"OP", e.GetLexeme().Value.(string)})
	}
}

var (
	AstDebugger = log.New(io.Discard, "AstDebugger: ", log.LstdFlags)
)

func (s *Scalar) Output(w *Writer) {
	if s == nil {
		return
	}
	write(w, []string{"expr", s.GetLexeme().ToString()})
	writeln(w, nil)
}
func ToString(all ...Anode) string {
	var buf []string
	for _, v := range all {
		buf = append(buf, v.GetLexeme().ToString())
	}
	return strings.Join(buf, " ")
}

func (v *Variable) Output(w *Writer) {
	if v == nil {
		return
	}
	// write(w, []string{"variableBegin"})
	// write(w, []string{v.GetLexeme().Value.(string)})

	// defer write(w, []string{"variableEnd"})
	writeln(w, []string{"var", v.BaseNode.Lexeme.ToString()})
}

func (f *Factor) Output(w *Writer) {
	if f == nil {
		return
	}
	l := f.BaseNode.Lexeme
	if l == nil {
		return
	}
	writeln(w, []string{"varf", l.ToString()})
}

func (f *FnParam) Output(w *Writer) {
	if f == nil {
		return
	}
	writeln(w, []string{"fn_arg_count", cast.ToString(len(f.Children))})
	for i, v := range f.Children {
		if v == nil {
			continue
		}
		writeln(w, []string{"fn_arg", v.GetLexeme().ToString(), cast.ToString(i)})
	}
}

func (f *ReturnStmt) Output(w *Writer) {
	if f == nil {
		return
	}
	if len(f.Children) == 0 {
		writeln(w, []string{"@RETURN", w.Ln(), "@void"})
		return
	}
	k := w.Ln()
	writeln(w, []string{"@RETURNBEGIN", k})
	for _, v := range f.Children {
		if v == nil {
			continue
		}
		v.Output(w)
	}
	writeln(w, []string{"@RETURN", k})
}

func (f *BreakFlagStmt) Output(w *Writer) {
	if f == nil {
		return
	}
	writeln(w, []string{"ExitFor",w.CtxForLabel})
}

func (f *DeclareStmt) Output(w *Writer) {
	if f == nil {
		return
	}
	writeln(w, []string{"declare", f.Children[0].GetLexeme().VarName()})
	for i := 1; i < len(f.Children); i++ {
		op := f.Children[i]
		op.Output(w)
	}
	writeln(w, []string{"ASSIGN", f.Children[0].GetLexeme().VarName()})
}

func (f *ForStmt) Output(w *Writer) {
	if f == nil {
		return
	}
	entryLabel := w.Ln()
	cnt  := len(f.Children)-1
	writeln(w, []string{"forStmtBegin", entryLabel, cast.ToString(cnt)})
	defer writeln(w, []string{"forStmtEnd", entryLabel})
	backUp := w.CtxForLabel  
	defer func() {
		w.CtxForLabel = backUp
	}()
	w.CtxForLabel = entryLabel
	if len(f.Children) == 4 {

		first := f.Children[0].(*ExprGroups)
		continueExpr := f.Children[1].(*ExprGroups)
		stepExpr := f.Children[2].(*ExprGroups)
		bodyExpr := f.Children[3].(*BlockNode)
		writeln(w, []string{"for_init",entryLabel})
		first.Output(w)
		writeln(w, []string{"for_init_end",entryLabel})
		writeln(w, []string{"for_continue",entryLabel})
		continueExpr.Output(w)
		writeln(w, []string{"for_continue_end",entryLabel,"ExitFor"})
		writeln(w, []string{"for_step",entryLabel})
		stepExpr.Output(w)
		writeln(w, []string{"for_step_end",entryLabel})
		writeln(w, []string{"for_body",entryLabel})
		bodyExpr.Output(w)
		writeln(w, []string{"for_body_end",entryLabel})
		writeln(w, []string{"GOTO","for_continue",entryLabel})
		writeln(w, []string{"ExitFor",entryLabel})
	} else if len(f.Children) >= 1 {
		ln := 0
		if len(f.Children) == 1 {
			ln = 0
			writeln(w, []string{"@deadloop",entryLabel})
		} else {
			ln = 1
			writeln(w, []string{"for_continue",entryLabel})
			f.Children[0].Output(w)
			writeln(w, []string{"for_continue_end",entryLabel,"ExitFor"})
		}
		for _, v := range f.Children[ln:] {
			if v == nil {
				continue
			}
			v.Output(w)
		}
		if ln == 0 {
			writeln(w, []string{"@deadloop_end",entryLabel})
		}
		writeln(w, []string{"GOTO","for_continue",entryLabel})
		writeln(w, []string{"ExitFor",entryLabel})
	}
}

func (f *FuncStmt) Output(w *Writer) {
	if f == nil {
		return
	}
	labelName := w.Ln()
	fnName := f.GetLexeme().Value.(string)
	paramCount := 0
	if param, ok := f.Children[0].(*FnParam);ok {
		paramCount = len(param.Children)
	}
	writeln(w, []string{"funcStmtBegin", labelName, fnName, cast.ToString(paramCount)})
	writeln(w, []string{"funcEntry", fnName,labelName, cast.ToString(paramCount)})
	defer writeln(w, []string{"funcStmtEnd", labelName, fnName})
	if len(f.Children) == 0 {
		writeln(w, []string{"func_param_empty", labelName})
		return
	}
	param, ok := f.Children[0].(*FnParam)
	if ok {
		param.Output(w)
	}
	for _, v := range f.Children[1:] {
		if v == nil {
			continue
		}
		v.Output(w)
	}
}

func (f *IfStmt) Output(w *Writer) {
	if f == nil {
		return
	}
	entryLabel := w.Ln()
	writeln(w, []string{"#IF", entryLabel})
	defer writeln(w, []string{"#ENDIF", entryLabel})
	bodyln := w.Ln()
	w.AddI(1)
	defer w.AddI(-1)
	writeln(w, []string{"if", entryLabel})
	cond := f.Children[0]
	// variable or expr 
	// todo: check cond or fatal

	cond.Output(w)
	skipCommand := []string{"endif",entryLabel}
	if len(f.Children) > 2 {
		// has else condition
		skipCommand = append(skipCommand, []string{"ELSEIF", entryLabel}...)
	}
	writeln(w, skipCommand)
	writeln(w, []string{"ifbodystart", bodyln})
	body, ok := f.Children[1].(*BlockNode)
	if !ok {
		writeln(w, []string{"@fatal", "if stmtbody syntax"})
		return
	}
	body.Output(w)
	writeln(w, []string{"ifbodyend", bodyln})
	writeln(w,[]string{"GOTO","#ENDIF",entryLabel})
	if len(f.Children) > 2 {
		writeln(w, []string{"#ELSEIF", entryLabel})
		for _, v := range f.Children[2:] {
			if v == nil {
				continue
			}
			v.Output(w)
			writeln(w, []string{"GOTO", "#ENDIF", entryLabel})
		}
		// writeln(w, []string{"#ENDELSEIF", entryLabel})
	}

}
func (b *BaseNode) Output(w *Writer) {
	if b == nil {
		return
	}
	if b.GetLexeme() != nil {
		writeln(w, []string{b.GetLexeme().ToString()})
	}
	for _, v := range b.GetChildren() {
		if v == nil {
			continue
		}
		v.Output(w)
	}
}

func (f *BlockNode) Output(w *Writer) {
	ln := w.Ln()
	writeln(w, []string{"blockstart", ln})
	defer writeln(w, []string{"blockend", ln})
	if f.GetLexeme() != nil {
		writeln(w, []string{f.GetLexeme().ToString()})
	}
	f.BaseNode.Output(w)
}

func (f *CallFuncStmt) Output(w *Writer) {
	e := w.Ln()
	paramCnt := len(f.Children[1:])
	writeln(w, []string{"call_arg", e, f.GetLexeme().Value.(string), cast.ToString(paramCnt)})
	for _, v := range f.Children[1:] {
		if v == nil {
			continue
		}
		v.Output(w)
	}
	writeln(w, []string{"call", e, f.GetLexeme().Value.(string),cast.ToString(paramCnt)})
}

// TODO:
func (f *MapIndexNode) Output(w *Writer) {
	// writeln(w, []string{"mapIndexNodeBegin"})
	// defer writeln(w, []string{"mapIndexNodeEnd"})
	// writeln(w, []string{f.GetLexeme().Value.(string)})
	ch := f.Children
	varn := ch[0].GetLexeme().VarName()
	varkey := ch[1].GetLexeme().ToString()
	writeln(w, []string{"@getvalue", varn, varkey})
}

func (f *ExprGroups) Output(w *Writer) {
	if len(f.Children) != 1 {
		writeln(w, []string{"@TODO"})
		return
	}
	for _, v := range f.Children {
		v.Output(w)
	}
}


// Output implements Anode.
func (c *KeywordStmt) Output(w *Writer) {
	writeln(w, []string{c.GetLexeme().Value.(string),w.CtxForLabel})
}