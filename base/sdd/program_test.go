package sdd

import (
	"testing"
	"toylang/base/ast"
)

func Test_parse_cmd(t *testing.T) {
	var s = "a= 1*2+3*4 + a*b"

	node := ParseNode(s)
	t.Logf("%+v\n", ast.ShowTree(node))
	d := Translate(node)
	// t.Logf("%+v\n", d.String())
	for _, v := range d.Cmd {
		if v != nil {
			t.Logf("%v\n", v)
		}
		// if i == 3 {
		// 	t.Logf("last=%v %+v\n", v.Operator, v.Arg1)
		// }
	}
}

func Test_parse_cmd22(t *testing.T) {
	var s = "var c=b*2+3*4 + a*b"

	node := ParseNode(s)
	t.Logf("%+v\n", ast.ShowTree(node))
	d := Translate(node)
	// t.Logf("%+v\n", d.String())
	for _, v := range d.Cmd {
		if v != nil {
			t.Logf("%v\n", v)
		}
		// if i == 3 {
		// 	t.Logf("last=%v %+v\n", v.Operator, v.Arg1)
		// }
	}
}

func Test_parse_cmd33(t *testing.T) {
	var s = ` // var a=1;
		{
			var b = a*2
		}	
			{
				var c = a*1000
			}
	`

	node := ParseNode(s)
	t.Logf("%+v\n", ast.ShowTree(node))
	d := Translate(node)
	// t.Logf("%+v\n", d.String())
	for _, v := range d.Cmd {
		if v != nil {
			t.Logf("%v\n", v)
		}
		// if i == 3 {
		// 	t.Logf("last=%v %+v\n", v.Operator, v.Arg1)
		// }
	}
}

func Test_parse_cmd_IF(t *testing.T) {
	var s = `
		if(a==1) {
			a =b+1;
		
		}  else {
		//	a = b-1;
		} 
		fn app(name,age) {
			var a = 1;
			a = a+1;
			return 1
		}
		
		var c = 666
		c = .app(1,2)
	`

	node := ParseNode(s)
	t.Logf("%+v\n", ast.ShowTree(node))
	d := Translate(node)
	// t.Logf("%+v\n", d.String())
	for _, v := range d.Cmd {
		if v != nil {
			t.Logf("%v\n", v)
		}
		// if i == 3 {
		// 	t.Logf("last=%v %+v\n", v.Operator, v.Arg1)
		// }
	}
}

func Test_parse_cmd_func(t *testing.T) {
	var s = `
		if(a==1) {
			a =b+1;
		
		}  else {
		//	a = b-1;
		} 
		fn app(name,age) {
			var a = 1;
			a = a+1;
			return 1
		}
		
		var c = 666
		c = .app(1,2)
	`

	node := ParseNode(s)
	t.Logf("%+v\n", ast.ShowTree(node))
	d := Translate(node)
	// t.Logf("%+v\n", d.String())
	for _, v := range d.Cmd {
		if v != nil {
			t.Logf("%v\n", v)
		}
		// if i == 3 {
		// 	t.Logf("last=%v %+v\n", v.Operator, v.Arg1)
		// }
	}
}
