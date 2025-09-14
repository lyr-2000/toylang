package ast

import (
	"fmt"
	"github.com/lyr-2000/toylang/base/lexer"
)

func parseDeclareStmt(t *Tokens) Anode {
	for t.i < len(t.tokens) && t.tokens[t.i].Value == ";" {
		t.i++
	}
	if t.tokens[t.i].Value == "var" {
		t.i++
		// read variable
		nextWord := t.tokens[t.i]
		if nextWord.Type != lexer.Variable {
			panic("illegal match at define variable")
		}
		var variableNode = new(Variable)
		variableNode.Lexeme = nextWord

		var stmt = new(DeclareStmt)
		t.i++
		if t.tokens[t.i].Value == ";" {
			t.i++
			stmt.Children = []Anode{variableNode}
			return stmt
		}
		//if match equal
		if t.tokens[t.i].Value == "=" {
			t.i++
			expr := parseExpr(t)
			stmt.Children = []Anode{variableNode, expr}
			return stmt
		}
		panic(fmt.Sprintf("illegal state declare func, %v", t.tokens[t.i]))
		/*
			var a = 1
			      def  --> leftchid_a
				  def --> rightchild_1
		*/

		// return stmt
	} else {
		//is not declare stmt
		panic("illegal syntax")
	}
}
