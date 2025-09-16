package compiler

import (
	"io"
	"strings"

	"github.com/lyr-2000/toylang/base/ast"
	"github.com/lyr-2000/toylang/base/evaluator"
)

func Compile(raw string) string {
	var buf strings.Builder
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
