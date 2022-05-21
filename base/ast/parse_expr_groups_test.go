package ast

import (
	"fmt"
	"testing"
	"toylang/base/lexer"
	// "golang.org/x/tools/go/analysis/passes/printf"
)

func Test_parseExprGroups(t *testing.T) {
	var s = `
		a==1,b==2,3==3
	
	`
	a := lexer.NewStringLexer(s)

	ls := parseExprGroups(&Tokens{i: 0, tokens: a.ReadTokens()})

	fmt.Printf("%+v\n", ShowTree(ls))

}
