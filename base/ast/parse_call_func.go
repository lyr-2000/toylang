package ast

import "toylang/base/lexer"

func isExpr__(n *lexer.Token) bool {
	if n == nil {
		return false
	}
	if n.Value == "." {
		//call statement
		return false
	}
	switch n.Type {
	case lexer.Boolean, lexer.Char, lexer.String, lexer.Variable, lexer.Number:
		return true

	default:

	}
	return false
}

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
		var needDot = false
		for {
			if !t.hasNext() {
				break
			}
			if t.current().Value == ";" {
				break
			}

			hasParameter := false
			cur := t.current()
			if needDot && cur.Value != "," {
				break
			}
			if cur.Value == "," {
				hasParameter = true
				t.next()
				needDot = true
			} else if isExpr__(cur) {
				hasParameter = true
				needDot = true
			} else {
				panic("cannot explain the fn call statement")
			}
			if hasParameter {
				exprNode = append(exprNode, parseExpr(t))
			} else {
				break
			}

			// else {
			// 	break
			// }
			// if dotCnt > 0 {
			// 	paramCnt := dotCnt + 1
			// 	for paramCnt > 0 {
			// 		exprNode = append(exprNode, parseExpr(t))
			// 		paramCnt--
			// 		if t.peek() == "," {
			// 			t.next()
			// 		}

			// 	}
			// }
			// endl

		}

	}

	stmt.Children = exprNode
	return stmt
}
