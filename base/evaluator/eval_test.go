package evaluator

import (
	"encoding/json"
	"fmt"
	"testing"
	"toylang/base/ast"

	"github.com/spf13/cast"
)

func Test_aaa(t *testing.T) {
	var b float64 = 1.3
	fmt.Printf("%+v\n", cast.ToInt64(b)&1)
	var a int = 1
	fmt.Printf("%+v\n", a&1)

}

func Test_run_code(t *testing.T) {
	c := NewCodeRunner()
	code := `
	var a = 1;
    b = 2; 
	a = 1+b; 
	b = "hello world" + b;
	fn app() {
		.print (a,b,"\nhello world")
		// a = 99999
		.print (a);
		 a = 88;
		.print("global a=",a)
		var a = 999; 
		.print("app -> stack a=",a);
		a = 888888
		.print("app call a=",a)
		return 666
	}

	fn b() {
		return "hello world bb"
	}
	// 这个时候 a == 3 
	 .println (.app() );
	 //.print ("global a=",a);
	 .println()
	 .println 1,2,3;
	 // .println;
	 //.println "i am str",66;
	// .println "e"
	// .println  .b()
	// .println (.b())
	 a = .b() 
	 .println("666")

	`
	tree := parse_source_tree(code)
	t.Logf("%+v\n", ast.ShowTree(tree))
	c.RunCode(tree)
	bs, _ := json.Marshal(c)

	t.Logf("%v\n", string(bs))
}

/*
规定函数调用语法：
如果 函数参数都是常量，或者 variable的名字，可以省略括号，
如果 存在 复杂的 函数调用表达式 ，则直接报错，
如果只是普通的加减的话，则可以省略括号

*/
func Test_run_code111(t *testing.T) {
	c := NewCodeRunner()
	code := `
	 
	fn b() {
		return "hello world bb"
	}
	 
	.println(.b() ,"aaa")

	.println (.b() )
	.println (.b() ) ;
	.println 1,2,3 ,"5555"
	.println;

	.println "for next"
	.println 1+2+3+4


	`
	tree := parse_source_tree(code)
	t.Logf("%+v\n", ast.ShowTree(tree))
	c.RunCode(tree)
	bs, _ := json.Marshal(c)

	t.Logf("%v\n", string(bs))
}
