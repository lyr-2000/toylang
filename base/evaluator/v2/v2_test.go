package evaluator

import (
	"fmt"
	"strings"
	"testing"

	"github.com/lyr-2000/toylang/base/ast"
	"github.com/lyr-2000/toylang/base/compiler"
	"github.com/lyr-2000/toylang/base/evaluator"
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
for(a<5) {
	a = a + 1
	print(a)
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
// errno = 333
	a = a + 1
	print(a)
	continue
}
print("a=",a)
print("errno=",errno)
	`
	sc := compiler.Compile(raw)
	t.Log("~~")
	t.Log(sc)

	b.SetReader(strings.NewReader(sc))
	b.Handle()
	fmt.Println(b.Top().Str())
	t.Log(b.Top().Str())
}

/*
   funcStmtBegin L2 getUser 3
        funcEntry getUser L2 3
        fn_arg_count 3
        fn_arg variable a 0
        fn_arg variable b 1
        fn_arg variable c 2
        blockstart L3
        call_arg L4 print 1
        var variable a
        call L4 print 1
        blockend L3
        funcStmtEnd L2 getUser
        call_arg L5 getUser 3
        expr number 2
        expr number 3
        expr number 3
        call L5 getUser 3
        printstack
        blockend L1


*/

func Test_func_decl(t *testing.T) {
	b := New()
	raw := `
func getUser(a,b,c) {
    if c == 3 {
	    print("C",c)
		return 44
	}
	return 233
}

d = getUser(2,3,3)

print(d)

for i=0;i<3;i++ {
	print(d)
}

throw( 12000,"error")
print(recover())
print("---")

printstack;
	`
	sc := compiler.Compile(raw)
	t.Log("~~")
	t.Log(sc)
	d := evaluator.ParseSourceTree(raw)
	t.Log(ast.ShowTree(d))
	b.SetReader(strings.NewReader(sc))
	b.Handle()
	// t.Log(b.Top().Str())
}

func Test_throw_recover(t *testing.T) {
	b := New()
	raw := `

fn call1(inputStr) {
    arr = array(1,2,3)
    print(arr[2])
    //print(arr[3])
    print(1,2,3)
    set(arr,0,2)
    print(arr)
    print("inputStr",inputStr);
    exit(0)
}

call1("~")


	`
	sc := compiler.Compile(raw)
	t.Log("byte codes :", sc)
	t.Log("~~")
	b.SetReader(strings.NewReader(sc))
	b.Handle()
}



func Test_parseAndRun(t *testing.T) {
	b := New()
	raw := `

b = "3" + 1
print("b=",b)

el = 0
el |= 1
print("el=",el)
el &= 0
print("el &=0 ",el)

fn call1(inputStr) {
    mapValue = map("a",1,"b",2,"c",3)
	mapValue["a"] = 666
	print(mapValue)
	sliceTest = array(1,2,3)
	sliceTest[1] = 8
	sliceTest = append(sliceTest,4)
	sliceTest[0] = 999
	print(sliceTest)
	sliceTest = remove(sliceTest,0)
	print(sliceTest)
	print("removeMapTest",remove(mapValue,"a"))
}

call1("~")

if errno > 0 {
	print(recover())
}else {
	print("success" )
}

print(typeof(sliceTest))
	`
	// b.ParseAndRun(raw)
	byteCode := compiler.Compile(raw)
	t.Log("byte codes :", string(byteCode))
	b.SetReader(strings.NewReader(string(byteCode)))
	b.Handle()
	t.Log(b.UnionStack.FuncStackCount)
}




func Test_parseStackoverflowCheck(t *testing.T) {
	b := New()
	raw := `

func call1(d) {
	if (d>100000) {
		return
	}
	print(d)
    call1(d+1)
}

call1(0);
if errno > 0 {
	print(recover())
}else {
	print("success" )
}
`
	// b.ParseAndRun(raw)
	byteCode := compiler.Compile(raw)
	t.Log("byte codes :", string(byteCode))
	b.SetReader(strings.NewReader(string(byteCode)))
	b.MaxStack = 100
	b.Handle()
	t.Log(b.UnionStack.FuncStackCount)
}
