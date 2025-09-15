package compiler

import (
	"io"

	"github.com/lyr-2000/toylang/base/ast"
)

func GenerateIntermediateCode(p ast.Anode, writer io.Writer) {
	p.Output(&ast.Writer{Writer: writer})
}

func todo() {
	panic("TODO")
}
