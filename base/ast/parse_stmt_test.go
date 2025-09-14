package ast

import (
	"testing"
	"github.com/lyr-2000/toylang/base/lexer"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_parseStmt(t *testing.T) {

	Convey("test_parseStmt", t, func() {
		lc := lexer.NewStringLexer("{ var a =1;var b= 2;var c = 3;var d = 4 }")
		ts := lc.ReadTokens()
		t.Logf("%+v\n", ts)
		p := &Tokens{ts, 0}
		st := parseBlock(p)
		t.Logf("%+v\n", st)
		t.Logf("%+v\n", toDfsPatternStringNode(st))
	})

}
