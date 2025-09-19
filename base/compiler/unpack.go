package compiler

import (
	"io"
	"strings"

	"github.com/lyr-2000/toylang/base/ast"
	"github.com/lyr-2000/toylang/base/evaluator"
)

var (
	CompilerVersion = "0.0.1"
)

func Compile(raw string) string {
	var buf strings.Builder

	buf.WriteString("// @remarks ToyLang Compiler Version: "+CompilerVersion+"\n")
	bn := evaluator.ParseSourceTree(raw)
	GenerateIntermediateCode(bn, &buf)
	return buf.String()
}

func GenerateIntermediateCode(p ast.Anode, writer io.Writer) {
	p.Output(&ast.Writer{Writer: writer})
}

func todo() {
	panic("TODO")
}
