package ast

import (
	"testing"
	"toylang/base/lexer"
)

func Test_define_astnode(t *testing.T) {
	var source = "var a = b;"
	tl := lexer.NewStringLexer(source)
	var ts = tl.ReadTokens()
	t.Logf("%+v\n", ts)
	tnode := parseDeclareStmt(&Tokens{ts, 0})
	t.Logf("%+v\n", toDfsPatternStringNode(tnode))
}
