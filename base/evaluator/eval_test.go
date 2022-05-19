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
	a = 1;
    b = 2;
	a = 1+b; 
	b = "hello world" + b;
	// 这个时候 a == 3 
	`
	tree := parse_source_tree(code)
	t.Logf("%+v\n", ast.ShowTree(tree))
	c.RunCode(tree)
	bs, _ := json.Marshal(c)

	t.Logf("%v\n", string(bs))
}
