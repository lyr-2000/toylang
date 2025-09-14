package ast

import (
	"fmt"
	"testing"
	"github.com/lyr-2000/toylang/base/lexer"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_parseFn(t *testing.T) {
	Convey("test_paresefn", t, func() {
		var s = `
			fn app (username,age)   {
				var a = 1+1;
				if a==1 {
					return 1;
				}else if a==2 {
					return 2; 
				}
				if(a!=3) {
					if a == 4 {
						return a*88
					}
					return a++ *1;
					
				}
				
				return (a++)*3;
			}
		
		
		`
		reader := lexer.NewStringLexer(s)
		tokens := reader.ReadTokens()
		t.Logf("%+v\n", tokens)
		var state = &Tokens{
			tokens: tokens,
			i:      0,
		}
		defer func() {
			err := recover()
			fmt.Printf("%T current = %s\n", state, state.peek())
			fmt.Printf("%+v\n", err)
		}()
		tree := parseFn(state)
		t.Logf("%+v\n", toDfsPatternStringNode(tree))
	})
}

func Test_parse_call_fn(t *testing.T) {
	Convey("test_parese_callfn", t, func() {
		var s = `.printf 1,2,3 ; 
		`
		reader := lexer.NewStringLexer(s)
		tokens := reader.ReadTokens()
		t.Logf("%+v\n", tokens)
		var state = &Tokens{
			tokens: tokens,
			i:      0,
		}
		defer func() {
			err := recover()
			fmt.Printf("%T current = %s\n", state, state.peek())
			fmt.Printf("%+v\n", err)
			// debug.PrintStack()

		}()
		tree := parseStmt(state)

		t.Logf("%+v\n", toDfsPatternStringNode(tree))
	})
	Convey("test_parese_callfn2", t, func() {
		var s = `a = .scanf("%d %d",a,b); 
		`
		reader := lexer.NewStringLexer(s)
		tokens := reader.ReadTokens()
		t.Logf("%+v\n", tokens)
		var state = &Tokens{
			tokens: tokens,
			i:      0,
		}
		defer func() {
			err := recover()
			fmt.Printf("%T current = %s\n", state, state.peek())
			fmt.Printf("%+v\n", err)
			// debug.PrintStack()

		}()
		tree := parseStmt(state)

		t.Logf("%+v\n", toDfsPatternStringNode(tree))
	})
}
