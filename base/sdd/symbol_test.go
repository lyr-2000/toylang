package sdd

import (
	"testing"
	"toylang/base/lexer"
)

func Test_symbol_add(t *testing.T) {
	var tab SymbolTab
	tab.AddVarSymbol()
	tab.AddVarSymbol()
	tab.AddSymbolByToken(&Token{Value: 1, Type: lexer.Number})
	tab.AddLabelSymbol("L0", nil)
	var tab2 SymbolTab
	tab.AddChild(&tab2)
	t.Logf("l=%v,size=%v\n", tab.String(), tab.Size())

	ok := tab2.ExistsToken(&Token{Value: 1, Type: lexer.Number})
	t.Logf("ok = %+v\n", ok)
}
