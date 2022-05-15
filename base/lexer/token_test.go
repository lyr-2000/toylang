package lexer

import "testing"

func Test_token_name(t *testing.T) {
	var tk TokenType = Keyword
	t.Logf("out1-> %+v\n", tk.String())
}
