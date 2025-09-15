package compiler

import (
	"io"

	"github.com/lyr-2000/toylang/base/ast"
)

type Writer interface {
	io.Writer
	io.StringWriter
}

func UnPackWithDebug(src string, tree ast.Anode, buf Writer) {
	treeSrc := ast.ShowTree(tree)
	buf.WriteString("<syntaxtree>")
	buf.WriteString(treeSrc)
	buf.WriteString("\n")
	buf.WriteString("</syntaxtree>")
	buf.WriteString("\n")
	if src != "" {
		buf.WriteString("<rawcode>")
		buf.WriteString(src)
		buf.WriteString("\n")
		buf.WriteString("</rawcode>")
	}
	buf.WriteString("\n")
	Unpack(tree, buf)
}
