package lexer

import (
	"bytes"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// func Test_(t *testing.T) {
// 	Convey("test_", t, func() {

// 	})
// }

func Test_letter(t *testing.T) {

	//lexer := New(bytes.NewBufferString("a"))
	Convey("test_IsLiteral", t, func() {
		//$key :=1
		//fmt.Println($key)
		So(IsLiteral(' '), ShouldEqual, false)
		So(IsLiteral('!'), ShouldEqual, false)
		So(IsLiteral('^'), ShouldEqual, false)
		So(IsLiteral('A'), ShouldEqual, true)
		So(IsLiteral('A'), ShouldEqual, true)
		//So('a', ShouldEqual, lexer.Next())
		//So(-1, ShouldEqual, lexer.Peek())
		//So(-1, ShouldEqual, lexer.Next())
		//So(-1, ShouldEqual, lexer.Next())
		//So(-1, ShouldEqual, lexer.Peek())
		//fmt.Println(lexer.Peek())
	})
}
func Test_reader(t *testing.T) {
	lexer := New(bytes.NewBufferString("a"))
	Convey("test_lexer_reader_hasNext", t, func() {

		So('a', ShouldEqual, lexer.Peek())
		So('a', ShouldEqual, lexer.Peek())
		So('a', ShouldEqual, lexer.Next())
		So(-1, ShouldEqual, lexer.Peek())
		So(-1, ShouldEqual, lexer.Next())
		So(-1, ShouldEqual, lexer.Next())
		So(-1, ShouldEqual, lexer.Peek())
		//fmt.Println(lexer.Peek())
		t.Logf("lexer=> %v", lexer.que)
	})
}
func Test_readKeywordOrVariableKey_(t *testing.T) {
	lexer := New(bytes.NewBufferString("var a = true;var b = false;" +
		"\nvar c = null;"))

	Convey("test_lexer_readKeywordOrVariableKey_", t, func() {
		//var key string
		tk := lexer.readKeywordOrVariableKey()
		So(tk.Value, ShouldEqual, "var")

		lexer.Next() // eat space
		tk = lexer.readKeywordOrVariableKey()
		So(tk.Value, ShouldEqual, "a")

	})

}

//. "github.com/smartystreets/goconvey/convey"
func Test_delete_comments(t *testing.T) {
	Convey("token 1", t, func() {
		So('1', ShouldBeBetweenOrEqual, '1', '9')
	})
	Convey("test_delete_comments", t, func() {
		l := NewStringLexer("/**/ /* *******  fasjfas \nbcdefghijik 🐮" +
			" ***fu,abc/  */" +
			" /* */ " +
			"" +
			"// aaa \n" +
			"" +
			"var a=1 ;" +
			"" +
			"\n" +
			"var b = a+1+4*333" +
			"" +
			"")
		var token []*Token = nil
		/*
			defer func() {
				err := recover()
				fmt.Println(l.Scanner.Pos())
				fmt.Println(l.LastPosString())
				fmt.Println(err)
			}()
		*/

		token = l.ReadTokens()
		So(token[0].Value, ShouldEqual, "var")
		So(token[1].Value, ShouldEqual, "a")
		So(token[2].Value, ShouldEqual, "=")
		So(token[3].Value, ShouldEqual, "1")
		t.Logf("%v", token)
	})

}

//. "github.com/smartystreets/goconvey/convey"
func Test_complex_expression(t *testing.T) {

	Convey("test_complex_expression", t, func() {
		//So(1+1, ShouldEqual, 2)
		l := NewStringLexer("" +
			"//b+c 赋值给a" +
			"\n" +
			"/*注释666666*/" +
			"var a =b+   c   --;\n\t\t")

		token := l.ReadTokens()
		t.Logf("%v", token)
		So(token[0].Value, ShouldEqual, "var")
		So(token[1].Value, ShouldEqual, "a")
		So(token[2].Value, ShouldEqual, "=")
		So(token[3].Value, ShouldEqual, "b")
		So(token[4].Value, ShouldEqual, "+")
		So(token[5].Value, ShouldEqual, "c")
	})

	Convey("test_complex_2", t, func() {
		var s = "var a= 111+6666+444+(999/2+888*888+999+876+5555+   77777/223+(96666+  888/ 2))"
		lexer := NewStringLexer(s)
		token := lexer.ReadTokens()
		t.Logf("%v", token)

	})

}

/*
TODO:
1. 配置 ctrl e 快捷键


2. 配置 ctrl j 快捷键


3. 配置 ctrl  t 快捷键

4. 配置 ctrl  left right 快捷键

5. 配置 ctrl space 快捷键
*/
