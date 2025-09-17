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
/*
#IF L2
	if L3
	var variable a
	endif L3 ELSEIF L2
	ifbodystart L3
	blockstart L4
	call_arg L5 print 1
	expr string 1
	call L5 print
	blockend L4
	ifbodyend L3
	GOTO #ENDIF L2
	#ELSEIF L2
	#IF L6
	if L7
	var variable a
	expr number 1
	OP >=
	endif L7 ELSEIF L6
	ifbodystart L7
	blockstart L8
	call_arg L9 print 1
	expr string 2
	call L9 print
	blockend L8
	ifbodyend L7
	GOTO #ENDIF L6
	#ELSEIF L6
	blockstart L10
	call_arg L11 print 1
	expr string 3
	call L11 print
	blockend L10
	GOTO #ENDIF L6
	#ENDIF L6
	GOTO #ENDIF L2
#ENDIF L2
call_arg L12 print 1
expr string end
call L12 print
blockend L1

*/
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
	if(b) {
		print(x)
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


func Test_simple_if(t *testing.T) {
	b := New()
	globalSet(b)
	raw := `
var a = 1
if(a==2) {
	print("1")
}else if (a < 0) {
	a = 444
}else {
	a = 888
}
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


func Test_simple_for(t *testing.T) {
	b := New()
	globalSet(b)
	raw := `
var a = 1
for(a<100) {
	a = a + 1
	continue
}
print(a)
printstack;
	`
	sc := compiler.Compile(raw)
	t.Log("~~")
	t.Log(sc)

	b.SetReader(strings.NewReader(sc))
	b.Handle()
	fmt.Println(b.Top().Str())
	t.Log(b.Top().Str())
}

func Test_complex_for(t *testing.T) {
	b := New()
	globalSet(b)
	raw := `
var a = 0
for(i=0;i<10;i++) {
	a = a + 1
	print(a)
	continue
}
	`
	sc := compiler.Compile(raw)
	t.Log("~~")
	t.Log(sc)

	b.SetReader(strings.NewReader(sc))
	b.Handle()
	fmt.Println(b.Top().Str())
	t.Log(b.Top().Str())
}