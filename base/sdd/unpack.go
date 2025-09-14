package sdd

import (
	"io"

	"github.com/lyr-2000/toylang/base/ast"
)

func Unpack(p ast.Anode, writer io.Writer) {
	p.Output(writer)
}

func todo() {
	panic("TODO")
}
