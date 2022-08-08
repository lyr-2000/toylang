package sdd

import (
	"fmt"
	"strings"
	"toylang/base/ast"
	"toylang/base/lexer"
)

type Symbol struct {
	Label   string     // 表示代码段标签
	Lexeme  *ast.Token //值，变量
	Parent  *SymbolTab
	Offset  int
	TabDeep int // 记录 符号表的作用域深度
}

func NewSymbol(t *ast.Token, o int) *Symbol {
	return &Symbol{
		Lexeme: t,
		Offset: o,
	}
}

type SymbolTab struct {
	Parent  *SymbolTab
	Childs  []*SymbolTab
	Symbols []*Symbol
	Index   int
	Offset  int
	Height  int //符号表的高度
}

func (u *SymbolTab) AddSymbol(b *Symbol) {
	u.Symbols = append(u.Symbols, b)
	b.Parent = u
}

// func (u
type Token = ast.Token

func (tab *SymbolTab) CurrentHasToken(u *Token) bool {
	for _, v := range tab.Symbols {
		if v != nil && v.Lexeme.Value == u.Value {
			return true
		}
	}
	return false
}
func (tab *SymbolTab) ExistsToken(u *Token) bool {
	for _, v := range tab.Symbols {
		if v.Lexeme == u {
			return true
		}
		if v.Lexeme.Value == u.Value {
			return true
		}
	}
	if tab.Parent != nil {
		return tab.Parent.ExistsToken(u)
	}
	return false
}

var (
	NullSymbol = &Symbol{Lexeme: &Token{Value: "null"}}
)

// 遍历符号表树，寻找到对应的 变量关键字 token
func findTokenFromSymbolTableTree(tab *SymbolTab, u *Token, deep int) *Symbol {
	if tab == nil || u == nil || u.Value == nil {
		return nil
	}
	for _, v := range tab.Symbols {
		if v != nil && v.Lexeme != nil && v.Lexeme.Value == u.Value {
			retv := *v
			//记录符号的高度【对应的符号表树的位置】
			retv.TabDeep = deep
			return &retv
		}
	}
	if tab.Parent != nil {
		return findTokenFromSymbolTableTree(tab.Parent, u, deep+1)
	}
	return nil
}
func (tab *SymbolTab) AddSymbolByToken(u *Token) *Symbol {
	if u == nil {
		return nil
	}
	var ret *Symbol = nil
	if u.Type == lexer.Number || u.Type == lexer.Char || u.Type == lexer.String {
		ret = NewSymbol(u, -1)
		tab.Symbols = append(tab.Symbols, ret)
		return ret
		// tab.Symbols = append(tab.Symbols, ret)
	} else {
		// 比如 b = a*1, 看看 a是否存在
		ret := findTokenFromSymbolTableTree(tab, u, 0)
		if ret == nil {
			ret = NewSymbol(u, tab.Offset)
			tab.Offset++
		}
		tab.Symbols = append(tab.Symbols, ret)

		return ret
	}
}
func (tab *SymbolTab) AddVarSymbol() *Symbol {
	var token = &Token{Type: lexer.Variable, Value: fmt.Sprintf("p%d", tab.Index)}
	tab.Index++
	ret := &Symbol{Lexeme: token, Offset: tab.Offset}
	tab.Offset++

	tab.Symbols = append(tab.Symbols, ret)
	return ret
}

func (tab *SymbolTab) AddLabelSymbol(labelName string, token *Token) *Symbol {
	// tab.Index++
	// var token = &Token{Type: lexer.Variable, Value: fmt.Sprintf("p%d", tab.Index)}
	// tab.Offset++
	ret := &Symbol{Lexeme: token /*  Offset: tab.Offset */, Label: labelName}

	tab.Symbols = append(tab.Symbols, ret)
	return ret
}

func (tab *SymbolTab) Size() int {
	return tab.Offset
}

func (tab *SymbolTab) String() string {
	var res strings.Builder
	for _, v := range tab.Symbols {
		res.WriteString(fmt.Sprintf("label=%v,value=%#v,index=%v\n", v.Label, v.Lexeme, v.Offset))
	}
	return res.String()
}
func (tab *SymbolTab) AddChild(u *SymbolTab) {
	u.Height = tab.Height + 1
	u.Parent = tab
	tab.Childs = append(tab.Childs, u)
}
