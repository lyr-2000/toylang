package ast

import (
	"testing"
	"github.com/lyr-2000/toylang/base/lexer"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_parseElseIfOrElse(t *testing.T) {
	Convey("test_if_else_if", t, func() {
		var source = `
			if(true) {
				var a = 33333;
			}else if(false ) {
				var t=2;
			}else {
				var a = c;
			}

		
		`
		l := lexer.NewStringLexer(source)
		ts := l.ReadTokens()
		t.Logf("tokens\n%+v\n", ts)
		s := &Tokens{i: 0, tokens: ts}
		node := parseIf(s)
		t.Logf("%+v\n", toDfsPatternStringNode(node))
		// t.Logf("%+v\n", toDfsPatternStringNode(parseIf(s)))

		// t.Logf("%+v\n", parseIf(s))
	})
}

func Test_condition_list(t *testing.T) {
	Convey("test_condition_list", t, func() {
		var source = `
		 
		if ( a!=b ) {
			b = a++
		}else if(a==4444) {
			c = 1
		} else if a=111 {

		}else if b ==2 {
			
		}
		
		
		else {

			cur = "aaa";	 
		}
		
		`
		l := lexer.NewStringLexer(source)
		ts := l.ReadTokens()
		t.Logf("tokens\n%+v\n", ts)
		s := &Tokens{i: 0, tokens: ts}
		node := parseIf(s)
		t.Logf("%+v\n", toDfsPatternStringNode(node))
		// fmt.Printf("curtoken %v, nextToken %v \n", s.tokens[s.i].Value, s.tokens[s.i+1].Value)
		// node = parseIf(s)
		// t.Logf("%+v\n", toDfsPatternStringNode(node))
	})
}
