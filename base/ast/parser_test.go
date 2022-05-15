package ast

import (
	"testing"
	"toylang/base/lexer"

	. "github.com/smartystreets/goconvey/convey"
)

//  hello world
func Test_ast_node_parse(t *testing.T) {
	Convey("test_ast_node_parse", t, func() {
		var souce = "var a = 1+1"
		var l = lexer.NewStringLexer(souce)

		t.Logf("%v\n", l.ReadToken())

	})
}

type Pf struct {
	PeekTokenIterator
}

func (ua *Pf) parse1() Anode {
	var (
		expr   Expr
		scalar Scalar
	)
	anode := ua.Next()
	scalar.Store(anode)
	if !ua.HasNext() {
		return &scalar
	}
	expr.Children = append(expr.Children, &scalar)
	// scalar.Store(ua.NextMatch())
	o, err := ua.NextMatch("+") //must
	if err != nil {
		panic(err)
	}
	expr.NodeType = BINARY_EXPR //二元表达式
	expr.Lexeme = o
	expr.Children = append(expr.Children, ua.parse1())
	return &expr

}

func Test_node_binary_plus_fun(t *testing.T) {
	Convey("test_node_binary_plus", t, func() {
		var s = "1+2+3+4"
		lx := lexer.NewStringLexer(s)
		var result = lx.ReadToken()
		t.Logf("%v\n", result)
		var pf = Pf{
			PeekTokenIterator: PeekTokenIterator{
				i:      0,
				tokens: result,
			},
		}
		node := pf.parse1()

		// t.Logf("result=%s\n", toDfsPatternStringNode(node))
		t.Logf("%+v\n", toDfsPatternStringNode(node))
		// So(node.GetChildren()[0].GetLexeme().Value, ShouldEqual, "+")

	})
}
