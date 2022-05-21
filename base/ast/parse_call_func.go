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
		// 原本要 在函数前面 添加 . 才能识别是函数调用，现在 补充功能，不需要加 . 就可以识别，但是需要带括号
		//panic("not an call statement")
		// currrent= println,  and next is (
		// println()
		if t.current().Type != lexer.Variable && t.current().Type != lexer.Keyword {
			panic("illegal statement on call parameter")
		}
		// current must be an key word
		ne := t.peekNextString(1)
		if ne != "(" {
			panic("cannot parse function statement ,need brackets \"(\" \")\" ")
		}
	}

	if t.current().Value == "." {
		t.next()
	}
	// .printf 1,2,3
	var stmt = new(CallFuncStmt)

	var exprNode []Anode
	var fnName = new(Variable)
	fnName.Lexeme = t.current()
	t.next()
	stmt.Lexeme = fnName.Lexeme
	//fn name
	exprNode = append(exprNode, fnName)
	if t.peek() == "(" { // 解析 带括号的 函数调用语句
		t.next()
		for t.hasNext() && t.peek() != ")" {
			if t.peek() == "," {
				t.next()
				continue
			}
			exprNode = append(exprNode, parseExpr(t))
		}
		t.next()

	} else { // parse like   print 1,2,3 , 这种没有括号的
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
			if cur.Value == "." {
				// 语法检查，及时发现问题
				panic("unsupport call func at . ,syntax error")
			}
			// check syntax
			{
				// check:
				i := 0
				for t.peekNextToken(i) != nil {
					w := t.peekNextToken(i)
					if w.Value == "." {
						panic("illegal syntax at dot .,please add \";\" on line end , or add brackets \"(\" \")\" to wrap params")

					}
					if w.Value == ";" {
						break
					}
					i++
				}
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
		}

	}

	stmt.Children = exprNode
	return stmt
}
