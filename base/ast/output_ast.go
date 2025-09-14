package ast

import (
	"io"
	"strings"
)

func write(w io.Writer, lines []string) {
	w.Write([]byte(strings.Join(lines, ",")))
	w.Write([]byte("\n"))
}

func (e *Expr) Output(w io.Writer) {
	if e == nil {
		return
	}
	write(w, []string{"exprBegin"})
	for _, v := range e.GetChildren() {
		v.Output(w)
	}
	if e.GetLexeme() != nil {
		write(w, []string{e.GetLexeme().Value.(string)})
	}
	defer write(w, []string{"exprEnd"})
}

func (s *Scalar) Output(w io.Writer) {
	if s == nil {
		return
	}
	write(w, []string{"scalarBegin"})
	defer write(w, []string{"scalarEnd"})
	write(w, []string{s.GetLexeme().Value.(string)})
}

func (v *Variable) Output(w io.Writer) {
	if v == nil {
		return
	}
	write(w, []string{"variableBegin"})
	write(w, []string{v.GetLexeme().Value.(string)})
	defer write(w, []string{"variableEnd"})
}

func (f *Factor) Output(w io.Writer) {
	if f == nil {
		return
	}
	write(w, []string{"factorBegin"})
	defer write(w, []string{"factorEnd"})
	write(w, []string{f.GetLexeme().Value.(string)})
}

func (f *FnParam) Output(w io.Writer) {
	if f == nil {
		return
	}
	write(w, []string{"fnParamBegin"})
	defer write(w, []string{"fnParamEnd"})
	write(w, []string{f.GetLexeme().Value.(string)})
	for _, v := range f.GetChildren() {
		if v == nil {
			continue
		}
		v.Output(w)
	}
}

func (f *ReturnStmt) Output(w io.Writer) {
	if f == nil {
		return
	}
	write(w, []string{"returnStmtBegin"})
	defer write(w, []string{"returnStmtEnd"})
	write(w, []string{f.GetLexeme().Value.(string)})
}

func (f *BreakFlagStmt) Output(w io.Writer) {
	if f == nil {
		return
	}
	write(w, []string{"breakFlagStmtBegin"})
	defer write(w, []string{"breakFlagStmtEnd"})
	write(w, []string{f.GetLexeme().Value.(string)})
}

func (f *DeclareStmt) Output(w io.Writer) {
	if f == nil {
		return
	}
	write(w, []string{"declareStmtBegin"})
	defer write(w, []string{"declareStmtEnd"})
	write(w, []string{f.GetLexeme().Value.(string)})
}

func (f *ForStmt) Output(w io.Writer) {
	if f == nil {
		return
	}
	write(w, []string{"forStmtBegin"})
	defer write(w, []string{"forStmtEnd"})
	write(w, []string{f.GetLexeme().Value.(string)})
}

func (f *FuncStmt) Output(w io.Writer) {
	if f == nil {
		return
	}
	write(w, []string{"funcStmtBegin"})
	defer write(w, []string{"funcStmtEnd"})
	write(w, []string{f.GetLexeme().Value.(string)})
}

func (f *IfStmt) Output(w io.Writer) {
	if f == nil {
		return
	}
	write(w, []string{"ifStmtBegin"})
	defer write(w, []string{"ifStmtEnd"})
	write(w, []string{f.GetLexeme().Value.(string)})
}
func (b *BaseNode) Output(w io.Writer) {
	if b == nil {
		return
	}
	write(w, []string{"baseNodeBegin"})

	defer write(w, []string{"baseNodeEnd"})

	if b.GetLexeme() != nil {
		write(w, []string{b.GetLexeme().Value.(string)})
	}

	for _, v := range b.GetChildren() {
		if v == nil {
			continue
		}
		v.Output(w)
	}
}

func (f *BlockNode) Output(w io.Writer) {
	write(w, []string{"blockNodeBegin"})
	defer write(w, []string{"blockNodeEnd"})
	if f.GetLexeme() != nil {
		write(w, []string{f.GetLexeme().Value.(string)})
	}
	f.BaseNode.Output(w)
}

func (f *CallFuncStmt) Output(w io.Writer) {
	write(w, []string{"callFuncStmtBegin"})
	defer write(w, []string{"callFuncStmtEnd"})
	write(w, []string{f.GetLexeme().Value.(string)})
}

func (f *MapIndexNode) Output(w io.Writer) {
	write(w, []string{"mapIndexNodeBegin"})
	defer write(w, []string{"mapIndexNodeEnd"})
	write(w, []string{f.GetLexeme().Value.(string)})
}

func (f *ExprGroups) Output(w io.Writer) {
	write(w, []string{"exprGroupsBegin"})
	defer write(w, []string{"exprGroupsEnd"})
	write(w, []string{f.GetLexeme().Value.(string)})
}
