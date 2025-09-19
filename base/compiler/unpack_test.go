package compiler

import (
	"bytes"
	"os"
	"testing"

	"github.com/lyr-2000/toylang/base/ast"
	"github.com/lyr-2000/toylang/base/evaluator"
)
func TestMain(m *testing.M) {
	ast.AstDebugger.SetOutput(os.Stdout)
	m.Run()
}



func TestUnpack2(t *testing.T) {
	code := `
a = array(1,2,3)
print(a[0])
print(a["hello"])
f = a[d]+1;
print(f)

d = "hello world\n"

for (1) {

	print(2)
	break
}

if b > 0 {
	return 1	
} else if b > 2 {
	return 2
} else {
	return 3+print(x)
}
return 

	`
	tree := evaluator.ParseTree(code)
	buf := &bytes.Buffer{}
	treeSrc := ast.ShowTree(tree)
	os.WriteFile("../../test.tree.txt", []byte(treeSrc), 0644)
	UnPackWithDebug(code,tree,buf)
	os.WriteFile("../../test.txt", buf.Bytes(), 0644)
}

func TestUnpack(t *testing.T) {
	code := `
var a = 3+2
b = a+1
b++
b&=1
b = a&1
print(b,44,(8+1),"hello")

fn Main(start,end,b) {
	if b > 0 {
		return 1	
	}
	print(x)
	return 0
}
Main(2,3,1)
for i=0;i<100;i++ {
    print(i,i+1)
}
	
	`
	tree := evaluator.ParseTree(code)
	buf := &bytes.Buffer{}
	treeSrc := ast.ShowTree(tree)
	os.WriteFile("../../test.tree.txt", []byte(treeSrc), 0644)
	UnPackWithDebug(code,tree,buf)
	t.Log(buf.String())
	os.WriteFile("../../test.txt", buf.Bytes(), 0644)
}

func Test_todo(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todo()
		})
	}
}
