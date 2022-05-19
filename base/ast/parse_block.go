package ast

import "fmt"

func parseBlock(t *Tokens) Node {
	for t.i < len(t.tokens) {
		if t.tokens[t.i].Value == ";" {
			t.i++
		} else {
			break
		}
	}
	if t.i >= len(t.tokens) {
		return nil
	}
	var ch = t.peek()
	if ch != "{" {
		panic(fmt.Sprintf("illegal state  actual %v", ch))
	}
	if ch == "{" {
		var block = new(BlockNode)
		t.i++
		for {
			if t.peek() == ";" {
				t.i++
			}
			if !t.hasNext() {
				break
			}
			stmt := parseStmt(t)
			if stmt == nil {
				break
			}
			block.Children = append(block.Children, stmt)
			// return stmt
		}
		if !t.hasNext() {
			panic("cannot parse body at }")
		}

		if t.tokens[t.i].Value != "}" {
			panic(fmt.Sprintf("match block } fail ,actual %v", t.tokens[t.i].Value))
		}
		t.i++ //eat }
		return block
	}
	panic("is not a block")

}
