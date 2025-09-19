package ast

func parseReturnStmt(t *Tokens) Node {
	t.NextMust(0, "return")
	var ret = new(ReturnStmt)
	ret.Lexeme = t.tokens[t.i]
	t.next()
	if t.i < len(t.tokens) {
		nextval := t.tokens[t.i]
		if nextval.Value == "}" {
			return ret
		}
	}
	c := parseExpr(t)
	if c == nil {
		return ret
	}
	ret.Children = []Anode{
		c,
	}
	return ret
}
