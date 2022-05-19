package ast

func parseReturnStmt(t *Tokens) Node {
	t.NextMust(0, "return")
	var ret = new(ReturnStmt)
	ret.Lexeme = t.tokens[t.i]
	t.next()
	c := parseExpr(t)
	if c == nil {
		return ret
	}
	ret.Children = []Anode{
		c,
	}
	return ret
}
