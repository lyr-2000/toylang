package ast

import "strconv"

func SimpleEvalExpr(node Anode) float64 {
	if node == nil {
		return 0
	}

	lexeme := node.GetLexeme()
	ch := node.GetChildren()

	switch node.(type) {
	case *Expr:
		if len(ch) != 2 {
			panic("illegal syntax")
		}
		w := node.(*Expr)
		// 二元运算符 ，如 a+b 等
		if w.NodeType == BINARY_EXPR {
			switch lexeme.Value {
			case "+":
				return SimpleEvalExpr(ch[0]) + SimpleEvalExpr(ch[1])
			case "-":
				return SimpleEvalExpr(ch[0]) - SimpleEvalExpr(ch[1])
			case "*":
				return SimpleEvalExpr(ch[0]) * SimpleEvalExpr(ch[1])
			case "/":
				v := SimpleEvalExpr(ch[1])
				if v == 0 {
					panic("divide by zero")
				}
				return SimpleEvalExpr(ch[0]) / v
			default:

			}
		} else if w.NodeType == UNARY_EXPR {
			//一元运算符
		}
	// case :
	case *Scalar:
		v, err := strconv.ParseFloat(lexeme.Value.(string), 64)

		if err != nil {
			panic(err)
		}
		return v

	default:
		panic(lexeme)
	}

	return 0
}
