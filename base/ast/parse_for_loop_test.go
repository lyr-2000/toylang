package ast

import (
	"log"
	"testing"
	"github.com/lyr-2000/toylang/base/lexer"
)

func Test_parseForStmt(t *testing.T) {
	s := `
	 
	 for(i=0;i<n;i++) {
		 a = 1+2
	 }
	 
	 `
	a := lexer.NewStringLexer(s)
	ts := a.ReadTokens()

	st := &Tokens{
		i:      0,
		tokens: ts,
	}

	tree := parseForStmt(st)
	log.Printf("%+v\n", ShowTree(tree))

}

func Test_parseForStmt_exprs(t *testing.T) {
	s := `
	 
	 for(i=0;i<n;i++) a++ ,b+=1

	 println("hello")
	 
	 `
	a := lexer.NewStringLexer(s)
	ts := a.ReadTokens()

	st := &Tokens{
		i:      0,
		tokens: ts,
	}

	tree := parseForStmt(st)
	log.Printf("%+v\n", ShowTree(tree))

}

func Test_parseForStmt_exprs111(t *testing.T) {
	s := `
	 
	 for i=0;i<n;i++ {
		 a++
		 b++
	 }

	 for i=0;i<n;i++  a++;

	 println("--")
	 
	 `
	a := lexer.NewStringLexer(s)
	ts := a.ReadTokens()

	st := &Tokens{
		i:      0,
		tokens: ts,
	}

	tree := parseForStmt(st)
	log.Printf("%+v\n", ShowTree(tree))
	tree = parseForStmt(st)
	log.Printf("%+v\n", ShowTree(tree))

}
