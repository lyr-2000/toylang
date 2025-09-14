package sdd

import (
	"bytes"
	"testing"

	"github.com/lyr-2000/toylang/base/evaluator"
)

func TestUnpack(t *testing.T) {
	code := "a+b*c"
	tree := evaluator.ParseTree(code)
	buf := &bytes.Buffer{}
	Unpack(tree, buf)
	t.Logf("%v", buf.String())
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
