package sdd

import (
	"fmt"
	"github.com/lyr-2000/toylang/base/ast"
	"github.com/lyr-2000/toylang/base/lexer"
)

// 实现三地址代码转换

type Program struct {
	Cmd           []*Cmd
	LabelCnt      int
	GlobalVarTabs *SymbolTab //全局变量
	NodeVar       map[Node]*Symbol
}

// 五元组实现方式
type Itype = uint8
type Instruction struct {
	Type     Itype   //运算类型
	Result   *Symbol //返回值
	Operator string  //操作符
	Arg1     interface{}
	Arg2     interface{}
}
type Cmd = Instruction

const (
	_ Itype = iota
	ASSIGN
	IF
	LABEL
	GOTO
	RETURN
	PUSH
	CALL //call function
	STACKPOINT
)

func (u *Cmd) String() string {

	switch u.Type {
	case ASSIGN:
		if u.Arg2 != nil {
			return fmt.Sprintf("%+v = %v %v %v", u.Result.Lexeme.Value, u.Arg1.(*Symbol).Lexeme.Value, u.Operator, u.Arg2.(*Symbol).Lexeme.Value)
		} else {
			// temp := u.Result
			// fmt.Printf("%v\n", temp)
			u1 := u.Result.Lexeme.Value
			u2 := u.Arg1.(*Symbol).Lexeme.Value
			return fmt.Sprintf("%v = %v", u1, u2)
		}
	case STACKPOINT:
		return fmt.Sprintf("SP %v", u.Arg1)
	case IF:
		o := u.Arg2
		var o1 interface{}
		if o != nil {
			o1 = o.(*Cmd).Arg1.(*Symbol).Label
		}
		return fmt.Sprintf("IF %v ELSE %v", (u.Arg1.(*Symbol)).Lexeme.Value, o1)
	case GOTO:
		return fmt.Sprintf("GOTO %v", u.Arg1)
	case LABEL:
		a := u.Arg1
		if u.Arg2 != nil {
			return fmt.Sprintf("%v@%v:", a.(*Symbol).Label, u.Arg2.(Node).GetLexeme().Value)
		}
		return fmt.Sprintf("%v:", a.(*Symbol).Label)
	case RETURN:
		return fmt.Sprintf("RETURN %v", (u.Arg1))
	case PUSH:
		//参数压栈
		return fmt.Sprintf("PUSH_PARAM %v %v", u.Arg1, u.Arg2)
	case CALL:
		return fmt.Sprintf("CALL %v", (u.Arg1))
		// default:
		// 	return fmt.Sprintf("%v = %v", u.Result, u.Arg1)
	}
	return ""
}
func (b *Program) AddCmd(u *Cmd) {
	b.Cmd = append(b.Cmd, u)
}

func ParseNode(s string) Node {

	var lx = lexer.NewStringLexer(s)

	tt := lx.ReadTokens()
	b := ast.NewTokens(tt)
	tree := ast.ParseStmt(b)
	return tree
}
func (u *Program) AddFunLabelGetCmd(t *Token) *Cmd {

	label := &Symbol{
		Label:  fmt.Sprintf("L%v", u.LabelCnt),
		Lexeme: t,
	}

	b := &Cmd{
		Type:   LABEL,
		Result: label,
		Arg1:   label,
	}
	u.LabelCnt++
	u.Cmd = append(u.Cmd, b)
	return b
}

func (u *Program) AddLabelGetCmd() *Cmd {
	// var label = fmt.Sprintf("L%v", u.LabelCnt)
	label := &Symbol{
		Label:  fmt.Sprintf("L%v", u.LabelCnt),
		Lexeme: nil,
	}
	b := &Cmd{
		Type:   LABEL,
		Result: label,
		Arg1:   label,
	}
	u.LabelCnt++
	u.Cmd = append(u.Cmd, b)
	return b
}
