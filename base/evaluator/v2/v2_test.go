package evaluator

import (
	"fmt"
	"strings"
	"testing"
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
	fmt.Println(b.Top().S())
	t.Log(b.Top().S())
}