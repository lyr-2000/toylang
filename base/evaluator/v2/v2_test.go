package evaluator

import (
	"fmt"
	"strings"
	"testing"

	"github.com/lyr-2000/toylang/base/compiler"
)

func Test_v2_1(t *testing.T) {
	b := New()
	globalSet(b)
	b.SetReader(strings.NewReader(`
expr number 1
expr number 2
OP +
	`))
	b.Handle()
	fmt.Println(b.Top().Str())
	t.Log(b.Top().Str())
}

func Test_v2_2(t *testing.T) {
	b := New()
	globalSet(b)
	raw := `
var a = 3+1*5;
print(a);

	`
	sc := compiler.Compile(raw)
	sc += "\nprintstack\n"
	t.Log("~~")
	t.Log(sc)

	b.SetReader(strings.NewReader(sc))
	b.Handle()
	fmt.Println(b.Top().Str())
	t.Log(b.Top().Str())
}

func Test_v2_IFcond(t *testing.T) {
	b := New()
	globalSet(b)
	raw := `
	if(a) {
		print("1")
	}else if a >=1{
		print("2")
	}else{
		print("3")
	}
	print("end")
	`
	sc := compiler.Compile(raw)
	sc += "\nprintstack\n"
	t.Log("~~")
	t.Log(sc)

	b.SetReader(strings.NewReader(sc))
	b.Handle()
	fmt.Println(b.Top().Str())
	t.Log(b.Top().Str())
}


func Test_v2_Forcond(t *testing.T) {
	b := New()
	globalSet(b)
	raw := `

	for j<10 {
		print(k)
	}
	for {
	print(2)
	}
	for j<1{
		print(3)
	}
	`
	sc := compiler.Compile(raw)
	sc += "\nprintstack\n"
	t.Log("~~")
	t.Log(sc)
	return

	b.SetReader(strings.NewReader(sc))
	b.Handle()
	fmt.Println(b.Top().Str())
	t.Log(b.Top().Str())
}




func Test_v2_FnStmt(t *testing.T) {
	b := New()
	globalSet(b)
	raw := `
	err = app()
	if err != nil {
		print("error")
	}
	
	`
	sc := compiler.Compile(raw)
	sc += "\nprintstack\n"
	t.Log("~~")
	t.Log(sc)
	return

	b.SetReader(strings.NewReader(sc))
	b.Handle()
	fmt.Println(b.Top().Str())
	t.Log(b.Top().Str())
}