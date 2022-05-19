package ast

func parseCallFuncStmt(t *Tokens) Node {
	if !t.hasNext() {
		return nil
	}
	if t.current().Value != "." {
		panic("not an call statement")
		//
	}

	t.next()
	// .printf 1,2,3
	var stmt = new(CallFuncStmt)

	var exprNode []Anode
	var fnName = new(Variable)
	fnName.Lexeme = t.current()
	t.next()
	stmt.Lexeme = fnName.Lexeme
	//fn name
	exprNode = append(exprNode, fnName)
	if t.peek() == "(" {
		t.next()
		for t.hasNext() && t.peek() != ")" {
			if t.peek() == "," {
				t.next()
				continue
			}
			exprNode = append(exprNode, parseExpr(t))
		}
		t.next()

	} else {
		for {
			if !t.hasNext() {
				break
			}
			if t.current().Value == ";" {
				break
			}
			// if t.current().Value == "," {
			// 	t.next()
			// 	exprNode = append(exprNode, parseExpr(t))
			// 	paramCnt++
			// 	hasNext = true
			// 	continue
			// } else {
			// 	if paramCnt > 0 && !hasNext {
			// 		break
			// 	}
			// 	exprNode = append(exprNode, parseExpr(t))
			// 	paramCnt++
			// }
			i := 0
			// paramCnt := 0
			dotCnt := 0
			for {
				tk := t.peekNext(i)
				if tk == "" || tk == ";" {
					break
				}
				if tk == "," {
					dotCnt++
				}
				i++
			}

			if dotCnt > 0 {
				paramCnt := dotCnt + 1
				for paramCnt > 0 {
					exprNode = append(exprNode, parseExpr(t))
					paramCnt--
					if t.peek() == "," {
						t.next()
					}

				}
			}
			// endl

		}

	}

	stmt.Children = exprNode
	return stmt
}
