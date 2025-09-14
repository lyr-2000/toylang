package lexer

import (
	"bytes"
	"fmt"
	"log"
	"testing"

	"regexp"

	. "github.com/smartystreets/goconvey/convey"
)

// import (
// 	"fmt"

func Test_regex_01(t *testing.T) {
	re := regexp.MustCompile("w+")
	// var aaa = 1
	res := re.Match([]byte("helloworld"))
	fmt.Println(res)

}

func Test_convey0(t *testing.T) {
	Convey("测试运行", t, func() {
		So(1, ShouldEqual, 1)
	})
}

func Test_token_string1(t *testing.T) {
	Convey("test string func", t, func() {
		So(true, ShouldEqual, IsNumber('1'))
		So(true, ShouldEqual, IsLiteral('_'))
		//So(false, ShouldEqual, IsOperator('-'))
		So(true, ShouldEqual, IsOperator('-'))
	})
}

func Test_token_readString_(t *testing.T) {
	Convey("test read_string_func", t, func() {
		//New(bytes.NewBufferString("'aaa'"))
		var s = "\"aaa\""
		//

		lexer := New(bytes.NewBufferString(s))
		token := lexer.readString_()
		//fmt.Println(str)
		So(token.Value, ShouldEqual, "aaa")

	})
	Convey("test read_string_func_with_quote", t, func() {
		//var a = "\"bbhelloworld\""
		var s = `"aaaaaa\"bbbb\""`

		lexer := New(bytes.NewBufferString(s))
		token := lexer.readString_()

		So(token.Value, ShouldEqual, "aaaaaa\"bbbb\"")

	})
}

//. "github.com/smartystreets/goconvey/convey"
func Test_lexer_readnumber(t *testing.T) {
	Convey("case read big float", t, func() {
		lexer := New(bytes.NewBufferString("1e9+7"))
		p := lexer.readNumberTok().Value
		So(p, ShouldEqual, "1e9+7")
		t.Logf("p: %v", p)
	})
	Convey("case space right", t, func() {
		lexer := New(bytes.NewBufferString("1.11112     "))
		So(lexer.readNumberTok().Value, ShouldEqual, "1.11112")
	})
	Convey("case space left", t, func() {
		lexer := New(bytes.NewBufferString(" 1.11112     "))
		lexer.Next()
		So(lexer.readNumberTok().Value, ShouldEqual, "1.11112")
	})
	//TODO: 这个地方 还要完善，目前不支持
	// return
	// Convey("case magic number 1e9+7", t, func() {
	// 	lexer := New(bytes.NewBufferString("1e9+7"))
	// 	//lexer.Next()
	// 	So(lexer.readNumber_().Value, ShouldEqual, "1e9+7")
	// })

}

//. "github.com/smartystreets/goconvey/convey"
func Test_string_equal(t *testing.T) {
	Convey("test_string_equal", t, func() {

		So(string('+'), ShouldEqual, "+")
		So(string('1'), ShouldEqual, "1")
		So(string('='), ShouldEqual, "=")

	})

}

//. "github.com/smartystreets/goconvey/convey"
func Test_read_operator(t *testing.T) {
	Convey("test_read_operator +=", t, func() {
		lexer := New(bytes.NewBufferString("+= "))
		So(lexer.readOperator_().Value, ShouldEqual, "+=")
	})

	Convey("test_read_operator //", t, func() {
		lexer := New(bytes.NewBufferString("//"))

		So(lexer.readOperator_().Value, ShouldEqual, "/")
		So(lexer.readOperator_().Value, ShouldEqual, "/")
	})

}

func Test_token_incr(t *testing.T) {
	Convey("test_token_incr", t, func() {
		var source = "a+++1"
		// a ,++ ,+ ,1
		var a = NewStringLexer(source)
		log.Printf("%+v\n", a.ReadTokens())
	})
}

func Test_token_plus(t *testing.T) {
	Convey("test_token_plus", t, func() {
		ts := NewStringLexer("+1-2+3")
		t.Logf("%+v\n", ts.ReadTokens())
	})
}
