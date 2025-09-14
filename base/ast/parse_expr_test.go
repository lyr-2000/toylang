package ast

import (
	"testing"
	"github.com/lyr-2000/toylang/base/lexer"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_parse_ast_tree_expr(t *testing.T) {
	Convey("test_parse_ast_tree_expr", t, func() {
		var source = "1+(2+5)*8*7;"

		var lc = lexer.NewStringLexer(source)
		ts := lc.ReadTokens()
		t.Logf("%+v\n", ts)
		var root = ParseExpr_(ts)
		t.Logf("%v", toDfsPatternStringNode(root))
		t.Logf("%+v\n", SimpleEvalExpr(root))
		So(SimpleEvalExpr(root), ShouldEqual, 1+(2+5)*8*7)
	})

	Convey("test_parse_ast_tree_expr2", t, func() {
		var source = "(1+2)+(1+3)*3"

		var lc = lexer.NewStringLexer(source)
		ts := lc.ReadTokens()
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
		ts := lc.ReadTokens()
		t.Logf("%+v\n", ts)
		var root = ParseExpr_(ts)
		So(len(ts), ShouldEqual, 4)
		t.Logf("%v", toDfsPatternStringNode(root))
	})
	Convey("test_node_unary", t, func() {
		var source = "(a++)+1"
		var lc = lexer.NewStringLexer(source)
		ts := lc.ReadTokens()
		t.Logf("%+v\n", ts)
		var root = ParseExpr_(ts)
		t.Logf("%v", toDfsPatternStringNode(root))
	})

	Convey("test_node_unary", t, func() {
		var source = "a+2"
		var lc = lexer.NewStringLexer(source)
		ts := lc.ReadTokens()
		t.Logf("%+v\n", ts)
		var root = ParseExpr_(ts)
		t.Logf("%v", toDfsPatternStringNode(root))
	})
}

func Test_array_index(t *testing.T) {
	Convey("test_array_index", t, func() {
		var a = "1+a[1] +  3 * b[\"age\"]"
		var lc = lexer.NewStringLexer(a)
		var ts = lc.ReadTokens()
		// So(len(ts), ShouldEqual, 6)
		t.Logf("%+v\n", ts)
		var astNode = ParseExpr_(ts)
		// t.Logf("%+v\n", astNode)
		// t.Logf("%#v \n,%T\n", astNode, astNode)

		t.Logf("%s\n", toDfsPatternStringNode(astNode))
	})

}

func Test_aa(t *testing.T) {
	Convey("test_!", t, func() {
		var a = "! a==b "
		var lc = lexer.NewStringLexer(a)
		var ts = lc.ReadTokens()
		So(len(ts), ShouldEqual, 4)
		t.Logf("%+v\n", ts)
		var astNode = ParseExpr_(ts)
		// t.Logf("%+v\n", astNode)
		// t.Logf("%#v \n,%T\n", astNode, astNode)

		t.Logf("%v\n", toDfsPatternStringNode(astNode))
	})
}

func Test_multi_line(t *testing.T) {
	Convey("test_multi_line", t, func() {
		var a = "a=b+10*33;b=2*7+99 "
		var lc = lexer.NewStringLexer(a)
		var ts = lc.ReadTokens()
		// So(len(ts), ShouldEqual, 4)
		t.Logf("%+v\n", ts)
		var s = &Tokens{i: 0, tokens: ts}
		var astNode = parseExpr(s)
		// t.Logf("%+v\n", astNode)
		// t.Logf("%#v \n,%T\n", astNode, astNode)

		t.Logf("case 1:%v\n", toDfsPatternStringNode(astNode))
		astNode = parseExpr(s)
		t.Logf("case 2:%+v\n", toDfsPatternStringNode(astNode))
	})
}
