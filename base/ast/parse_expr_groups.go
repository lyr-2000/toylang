package ast

// example:   e1,e2,e3,e4
func parseExprGroups(t *Tokens) Anode {
	res := new(ExprGroups)
	prv := -1
	repeat := 0
	for {
		ta := parseExpr(t)
		if t.i == prv {
			repeat++
		} else {
			prv = t.i
			repeat = 0
		}
		if repeat > 8 {
			// cannot  read expr
			panick("cannot read expr groups, code error %v", t.peekNextString(0))
		}
		res.Children = append(res.Children, ta)
		if t.peekNextString(0) == "," {
			t.next()
			continue
		}
		break
	}
	return res

}
