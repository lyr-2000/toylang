package sdd

import "github.com/lyr-2000/toylang/base/ast"

func IsFunNode(n Node) bool {
	switch n.(type) {
	case *ast.FuncStmt:
		return true
	}

	return false
}

func IsCallFnNode(n Node) bool {
	switch n.(type) {
	case *ast.CallFuncStmt:
		return true
	}
	return false
}

// 函数注册
func tran_fn_decl_stmt(p *Program, u Node, parentTab *SymbolTab) {
	node := u.(*ast.FuncStmt)
	label := p.AddFunLabelGetCmd(node.Lexeme)
	var tab SymbolTab
	tab.Parent = parentTab
	tab.AddVarSymbol()
	parentTab.AddChild(&tab)
	// 创建 label
	label.Arg2 = u
	// for _, v := range node.GetChildren() {
	// tab.AddSymbolByToken(v.GetLexeme())
	// }
	n := len(node.Children)
	for i := 0; i < n-1; i++ {
		tab.AddSymbolByToken(node.Children[i].GetLexeme())
	}
	body := node.Children[n-1]
	for _, v := range body.GetChildren() {
		tran_stmt(p, v, &tab)
	}

	// parentTab.AddChild()

}
func tran_call_fn_stmt(p *Program, u Node, ptab *SymbolTab) *Symbol {
	node := u.(*ast.CallFuncStmt)
	//call get_byId ,1,2,3
	n := len(u.GetChildren())
	var fnName = u.GetChildren()[0]
	retValue := ptab.AddVarSymbol()
	for i := 1; i < n; i++ {
		_ = _tran_expr(p, node.GetChildren()[i], ptab)
		p.AddCmd(&Cmd{PUSH, nil, "", nil, i - 1})
	}
	var fnAddr = findTokenFromSymbolTableTree(ptab, fnName.GetLexeme(), 0)

	if fnAddr == nil {
		panic("no funtion!")
	}
	p.AddCmd(&Cmd{STACKPOINT, nil, "", -ptab.Size(), nil})
	p.AddCmd(&Cmd{CALL, nil, "", fnAddr, nil})
	p.AddCmd(&Cmd{STACKPOINT, nil, "", ptab.Size(), nil})
	return retValue
}
func IsReturnStmt(u Node) bool {
	switch u.(type) {
	case *ast.ReturnStmt:
		return true
	}
	return false
}

func tran_ret_stmt(p *Program, n Node, t *SymbolTab) {
	node := n.(*ast.ReturnStmt)
	var retVal interface{}
	if len(n.GetChildren()) >= 1 {
		retVal = _tran_expr(p, node.GetChildren()[0], t)
	}
	p.AddCmd(&Cmd{RETURN, nil, "", retVal, nil})
}
