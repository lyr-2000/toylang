package ast

import "fmt"

func panick(s string, a ...interface{}) {
	// p

	panic(fmt.Sprintf(s, a...))

}

func read_for_init(t *Tokens) Anode {
	//
	return parseExprGroups(t)

}

func read_for_cond(t *Tokens) Anode {
	if t.peekNextString(0) == ";" {
		t.next()
	}
	//for( i=0;i<n;i++)
	return parseExprGroups(t)
}
func read_for_loop_end(t *Tokens) Anode {
	if t.peekNextString(0) == ";" {
		t.next()
	}
	return parseExprGroups(t)
}
func read_for_body(t *Tokens) Anode {
	if t.peekNextString(0) == ")" {
		t.next()
	}
	if t.peekNextString(0) == "{" {
		return parseBlock(t) // for() { do_some() }
	}
	//for (int i=0;i<n;i++) do_some()
	return parseExprGroups(t)
}

func parseForStmt(t *Tokens) Anode {
	if !t.hasNext() {
		return nil
	}

	cct := t.current().Value
	isLoopSt := cct == "for" || cct == "while"
	if !isLoopSt {
		panick("cannot explain loop statement %v", cct)
	}
	forStmt := new(ForStmt)
	forStmt.Lexeme = t.current()
	t.next()
	if t.peekNextString(0) == "{" {
		forStmt.Children = []Anode{read_for_body(t)}

		return forStmt
	}
	// expect  (conditoin) or  (start;condition;each_loop_end)

	if t.peekNextString(0) == "(" {
		t.next()
	}

	init_node := read_for_init(t)
	if t.peekNextString(0) == ")" {
		t.next()
		// if t.peekNextString(0) == ";" {
		// 	goto for3
		// }
	}

	if t.peekNextString(0) == "{" {
		// conditoin ,{ body }
		forStmt.Children = []Anode{init_node, read_for_body(t)}
	} else {
		// for3:
		// init,condition,loop_each { body }
		cond := read_for_cond(t)
		lp := read_for_loop_end(t)
		body := read_for_body(t)
		forStmt.Children = []Anode{init_node, cond, lp, body}
	}
	return forStmt
}

func (f *ForStmt) GetInitNode() Anode {
	n := len(f.Children)
	if n == 0 {
		return nil
	}
	if n == 1 {
		return nil
	}
	//for(0) {1} , for(0;1;2) {3} , for { 0 }
	if n != 2 && n != 4 {
		panick("cannot explain ast tree on for stmt %v", f.Children)

	}
	if n == 4 {
		return f.Children[0]
	}
	return nil
}

func (f *ForStmt) GetCondition() Anode {
	n := len(f.Children)
	if n == 1 {
		return nil
	}
	if n == 4 {
		return f.Children[1]
	}
	if n == 2 {
		return f.Children[0]
	}
	return nil

}

func (f *ForStmt) GetEachLoopAction() Anode {
	if len(f.Children) == 4 {
		return f.Children[2]
	}
	return nil
}

func (f *ForStmt) GetBody() Anode {
	if len(f.Children) == 4 {
		return f.Children[3]
	}
	if len(f.Children) == 1 {
		return f.Children[0]
	}

	if len(f.Children) == 2 {
		return f.Children[1]
	}

	panick("no for statement body , !! break down")
	return nil
}
