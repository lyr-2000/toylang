package sdd

import (
	"github.com/lyr-2000/toylang/base/ast"
	"github.com/lyr-2000/toylang/base/lexer"
)

func tran_block(p *Program, node Node, parentTab *SymbolTab) {
	var newTab *SymbolTab = new(SymbolTab)
	parentTab.Childs = append(parentTab.Childs, newTab)
	newTab.Parent = parentTab

	sp := newTab.AddVarSymbol()
	sp.Lexeme = &Token{Type: lexer.Number, Value: parentTab.Size()}
	// pushState :=
	stackPush := &Cmd{STACKPOINT, nil, "", nil, nil}
	p.Cmd = append(p.Cmd, stackPush)
	for _, v := range node.GetChildren() {
		tran_stmt(p, v, newTab)
	}
	stackPop := &Cmd{STACKPOINT, nil, "", nil, nil}
	p.Cmd = append(p.Cmd, stackPop)
	// 压入站
	stackPush.Arg1 = -parentTab.Size()
	//弹出栈
	stackPop.Arg1 = parentTab.Size()
}

func tran_if(p *Program, node Node, tab *SymbolTab) {
	ifNode := node.(*ast.IfNode)
	// ifNode.
	expr := ifNode.GetCondition()
	exprToken := _tran_expr(p, expr, tab)
	ifCmd := &Cmd{
		IF,
		nil,
		"",
		exprToken,
		nil,
	}
	p.AddCmd(ifCmd)
	tran_block(p, ifNode.GetBody(), tab)
	// var gotoCmd *Cmd
	if ifNode.GetElseNode() != nil {
		// gotoCmd = &Cmd{GOTO, nil, "", nil, nil}
		// elseLabel := p.AddLabel()
		gotoCmd := &Cmd{
			Type: GOTO,
			Arg1: nil,
		}
		p.AddCmd(gotoCmd)
		elseEnterLabel := p.AddLabelGetCmd()
		// tab.AddSymbol(elseEnterLabel)
		ifCmd.Arg2 = elseEnterLabel
		tran_block(p, node, tab)
		gotoLabel := p.AddLabelGetCmd()
		gotoCmd.Arg1 = gotoLabel
		/*
			if {
			} else {
			}
			elseLeave:
			==>
			if expr endLabel
			do 1,2
			do 2,3
			elseLabel:
				do ...
			endLabel:
				....
		*/
	}

}
