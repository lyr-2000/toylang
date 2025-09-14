package evaluator

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/lyr-2000/toylang/base/ast"
	"github.com/lyr-2000/toylang/base/lexer"

	"github.com/spf13/cast"
)

func Test_aaa(t *testing.T) {
	var b float64 = 1.3
	fmt.Printf("%+v\n", cast.ToInt64(b)&1)
	var a int = 1
	fmt.Printf("%+v\n", a&1)

}

func Test_assign(t *testing.T) {
	code := `
	a = map()
	set(a,"key",1)
	print(a)
	set(A,"Value","2")
	print(get(A,"Value"))
	print(typeof(A))

	`
	node := ParseTree(code)
	t.Logf("%+v\n", ast.ShowTree(node))
	type StructValue struct {
		Name string
		Value string
	}
	var a = &StructValue{
		Name: "hello",
		Value: "world",
	}
	runner := NewCodeRunner()
	runner.SetVar("A", a, true)
	runner.RunCode(node)
	t.Logf("refValue: %#v", a)
}
func Test_prevEval(t *testing.T) {
	t.Run("test-Prev_eval1-0", func(t *testing.T) {
		code := `
	A>B
	`
		node := ParseTree(code)
		runner := NewCodeRunner()
		runner.DebugLog = log.New(os.Stdout, "[Test_prevEval]", log.LstdFlags)
		runner.SetVar("A", 100, true)
		runner.SetVar("B", 2, true)
		runner.RunCode(node)
		t.Logf("pe %+v\n", runner.PrevEval)
		t.Logf("%+v\n", runner.ExitCode)

	})
	t.Run("test-Prev_eval21", func(t *testing.T) {
		code := `
	A<B
	`
		node := ParseTree(code)
		runner := NewCodeRunner()
		runner.SetVar("A", 100, true)
		runner.SetVar("B", 2, true)
		runner.RunCode(node)
		t.Logf("pe %+v\n", runner.PrevEval)
		t.Logf("%+v\n", runner.ExitCode)

	})
	t.Run("test-Prev_eval41", func(t *testing.T) {
		code := `
	contains A B && A>B
	`
		node := ParseTree(code)
		runner := NewCodeRunner()
		runner.SetFunc("contains", func(params []interface{}) interface{} {
			t.Logf("Contains Call! %+v\n", params)
			return strings.Contains(cast.ToString(params[0]), cast.ToString(params[1]))
		})
		t.Logf("%v\n", ast.ShowTree(node))

		runner.SetVar("A", 100, true)
		runner.SetVar("B", 1, true)
		runner.RunCode(node)
		t.Logf("prev_Eval %+v\n", runner.PrevEval)
		t.Logf("%+v\n", runner.ExitCode)

	})
	t.Run("test-Prev_eval31x", func(t *testing.T) {
		code := `
	contains(A,B) && A>B
	`
		node := ParseTree(code)
		t.Logf("%v\n", ast.ShowTree(node))
		runner := NewCodeRunner()
		runner.SetFunc("contains", func(params []interface{}) interface{} {
			t.Logf("Contains Call! %+v\n", params)
			return strings.Contains(cast.ToString(params[0]), cast.ToString(params[1]))
		})

		runner.SetVar("A", 100, true)
		runner.SetVar("B", 1, true)
		runner.RunCode(node)
		t.Logf("prev_Eval %+v\n", runner.PrevEval)
		t.Logf("%+v\n", runner.ExitCode)

	})
}

func Test_array(t *testing.T) {
	code := `
	var a = array(1,2,3)
	print("\n")
	print(a[0])
	var mpValue = map("a","hello\n","b",2)
	print(mpValue["a"])
	print(len(mpValue["a"]))
	print("===")
	`
	node := ParseTree(code)
	// t.Logf("%v\n", ast.ShowTree(node))
	runner := NewCodeRunner()
	runner.DebugLog.SetOutput(io.Discard)
	// runner.Logger = log.New(os.Stdout, "[Test_array]", log.LstdFlags)
	runner.RunCode(node)
}
func Test_eval(t *testing.T) {
	s := `
	if (A>B) &&  (C< D)   {
		exit(1)
	}else if 1==1{
		.exit(2)	
	}
	// if h == "g" {
	// 	.print(h)
	// 	.exit(3)
	// }
	`
	node := ParseTree(s)
	runner := NewCodeRunner()
	t.Logf("%v\n", ast.ShowTree(node))
	runner.SetVar("A", 100, true)
	runner.SetVar("B", 2, true)
	runner.SetVar("C", 3, true)
	runner.SetVar("D", 4, true)
	runner.RunCode(node)
	t.Logf("pe %+v\n", runner.PrevEval)
	t.Logf("%+v\n", runner.ExitCode)

}
func Test_run_code(t *testing.T) {
	c := NewCodeRunner()
	code := `
	var a = 1;
    b = 2; 
	a = 1+b; 
	b = "hello world" + b;
	fn app() {
		print (a,b,"\nhello world")
		// a = 99999
		print (a);
		 a = 88;
		.print("global a=",a)
		var a = 999; 
		print("app -> stack a=",a);
		a = 888888
		print("app call a=",a)
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
	tree := parseSourceTree(code)
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
	tree := parseSourceTree(code)
	t.Logf("%+v\n", ast.ShowTree(tree))
	c.RunCode(tree)
	bs, _ := json.Marshal(c)

	t.Logf("%v\n", string(bs))
}

func Test_run_plus_plus(t *testing.T) {
	c := NewCodeRunner()
	code := `
	 
	 var a = 1

	 a += 1 

	 a += 2

	 .println a   ; 

	 .println a+1  ;

	 .println a+1   ; 
	 .println 1,2,3 

//	 .println 3,4,5

	`
	ll := lexer.NewStringLexer(code)
	t.Logf("%+v\n", ll.ReadTokens())
	tree := parseSourceTree(code)
	t.Logf("%+v\n", ast.ShowTree(tree))
	c.RunCode(tree)
	bs, _ := json.Marshal(c)

	t.Logf("%v\n", string(bs))
}

func Test_run_if(t *testing.T) {
	c := NewCodeRunner()
	code := `
	  var a = true
	  if (a) {
		  println (true);

	  }else {
		  print(false);
	  }

	  if (a) {
		  println("ok it is true ")
	  }
	  fn app(username) {
		print (username);
	  }
	  app("hello world")
	`
	ll := lexer.NewStringLexer(code)
	t.Logf("%+v\n", ll.ReadTokens())
	tree := parseSourceTree(code)
	t.Logf("%+v\n", ast.ShowTree(tree))
	c.RunCode(tree)
	bs, _ := json.Marshal(c)

	t.Logf("%v\n", string(bs))
}

func Test_fib_stack_call(t *testing.T) {
	c := NewCodeRunner()
	//1 ,1 ,2 3
	code := `
	  fn fib(n) {
		 // println(n)
		  if n==1 || n ==2 {
			return 1
		  }
		  
		  
		  return fib(n-1) + fib(n-2)
	  }

	  var a = (2>=1)
	  println(a)
	  println("fib result=", fib(4) )
	  println("fib result = ",fib(12))
	  // 1 1 , 2 ,3

	`
	ll := lexer.NewStringLexer(code)
	t.Logf("%+v\n", ll.ReadTokens())
	tree := parseSourceTree(code)
	t.Logf("%+v\n", ast.ShowTree(tree))

	defer func() {
		// bs, _ := json.Marshal(&stru)

		// t.Logf("%v\n", string(bs))
		t.Logf("%+v\n", c.Vars)
		t.Logf("%+v\n", c.Stack.Queue)
	}()
	c.RunCode(tree)

}

func Test_for_each_stmt(t *testing.T) {
	c := NewCodeRunner()
	//1 ,1 ,2 3
	code := `
	   i = 0;
	   for (1)  {
		  println(i);
		  if i>=3 {
			  break 
		  } 
		  i+=1
	   }

	`
	ll := lexer.NewStringLexer(code)
	t.Logf("%+v\n", ll.ReadTokens())
	tree := parseSourceTree(code)
	t.Logf("%+v\n", ast.ShowTree(tree))

	defer func() {
		// bs, _ := json.Marshal(&stru)

		// t.Logf("%v\n", string(bs))
		t.Logf("%+v\n", c.Vars)
		t.Logf("%+v\n", c.Stack.Queue)
	}()
	c.RunCode(tree)

}

func Test_FIB(t *testing.T) {
	c := NewCodeRunner()
	//1 ,1 ,2 3
	code := `
	    fn fib(n) {
			var a=0
			var b=1
		//	var sum = 0;
			for i=0;i<n;i++ {
				sum = a+b;
				a = b;
				b  = sum;
			}
		//	println("斐波那锲数列答案：",a,"--",n);
			return sum
		}
		fn fib0(n) {
			if n <= 1{
				return 1
			}
			 
			return fib0(n-1)+fib0(n-2)
		}

		a = fib(10)
		b = fib0(10);
		println("get fib = ",a)
		println("get fib = ",b)
		println("sum= ",sum)
	`
	ll := lexer.NewStringLexer(code)
	t.Logf("%+v\n", ll.ReadTokens())
	tree := parseSourceTree(code)
	t.Logf("%+v\n", ast.ShowTree(tree))

	// defer func() {
	// 	// bs, _ := json.Marshal(&stru)

	// 	// t.Logf("%v\n", string(bs))
	// 	t.Logf("%+v\n", c.Vars)
	// 	t.Logf("%+v\n", c.Stack.Queue)
	// }()
	c.RunCode(tree)

}
