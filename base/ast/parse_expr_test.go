package ast

import (
	"testing"
	"toylang/base/lexer"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_parse_ast_tree_expr(t *testing.T) {
	Convey("test_parse_ast_tree_expr", t, func() {
		var source = "1+(2+5)*8*7"

		var lc = lexer.NewStringLexer(source)
		ts := lc.ReadToken()
		t.Logf("%+v\n", ts)
		var root = ParseExpr_(ts)
		t.Logf("%v", toDfsPatternStringNode(root))
		t.Logf("%+v\n", SimpleEvalExpr(root))
		So(SimpleEvalExpr(root), ShouldEqual, 1+(2+5)*8*7)
	})

	Convey("test_parse_ast_tree_expr2", t, func() {
		var source = "(1+2)+(1+3)*3"

		var lc = lexer.NewStringLexer(source)
		ts := lc.ReadToken()
		t.Logf("%+v\n", ts)
		var root = ParseExpr_(ts)
		t.Logf("%v", toDfsPatternStringNode(root))
		t.Logf("eval result-> %+v\n", SimpleEvalExpr(root))

	})
}

func Test_node_unary(t *testing.T) {
	Convey("test_node_unary", t, func() {
		var source = "a++ +1"
		var lc = lexer.NewStringLexer(source)
		ts := lc.ReadToken()
		t.Logf("%+v\n", ts)
		var root = ParseExpr_(ts)
		So(len(ts), ShouldEqual, 4)
		t.Logf("%v", toDfsPatternStringNode(root))
	})
	Convey("test_node_unary", t, func() {
		var source = "(a++)-1"
		var lc = lexer.NewStringLexer(source)
		ts := lc.ReadToken()
		t.Logf("%+v\n", ts)
		var root = ParseExpr_(ts)
		t.Logf("%v", toDfsPatternStringNode(root))
	})

	Convey("test_node_unary", t, func() {
		var source = "(a+=1)+2"
		var lc = lexer.NewStringLexer(source)
		ts := lc.ReadToken()
		t.Logf("%+v\n", ts)
		var root = ParseExpr_(ts)
		t.Logf("%v", toDfsPatternStringNode(root))
	})
}
