package ast

import (
	"fmt"
)

func parseIf(t *Tokens) Node {
	if t.tokens[t.i].Value != "if" {
		panic(fmt.Sprintf("match block if fail ,actual %v", t.tokens[t.i].Value))
	}
	var ifToken = t.tokens[t.i]
	t.i++
	var ifNode = new(IfStmt)
	ifNode.Lexeme = ifToken
	if !t.hasNext() {
		// eof
		return nil
	}

	//parse condition
	{
		//parse condition
		// if t.tokens[t.i].Value != "(" {
		// 	panic(fmt.Sprintf("match block } fail, actual %v", t.tokens[t.i]))
		// }
		// t.i++
		expr := parseExpr(t)
		// if t.tokens[t.i].Value != ")" {
		// 	panic(fmt.Sprintf("illegal match if condition on ),actual %v", t.tokens[t.i]))
		// }

		if expr == nil {
			panic(fmt.Sprintf("cannot match if condition expr actual "))
		}
		//expr
		ifNode.Children = append(ifNode.Children, expr)
	}
	{
		//parse body
		b := parseBlock(t)
		ifNode.Children = append(ifNode.Children, b)
	}
	{
		var elseNode = parseElseIfOrElse(t)
		//parse else or else if
		if elseNode != nil {

			ifNode.Children = append(ifNode.Children, elseNode)
		}
	}

	return ifNode
}

type IfNode = IfStmt

func (i *IfStmt) GetCondition() Node {
	if len(i.Children) <= 0 {
		return nil
	}
	return i.Children[0]
}
func (i *IfStmt) GetBody() Node {
	if len(i.Children) <= 1 {
		return nil
	}
	return i.Children[1]
}
func (i *IfStmt) GetElseNode() Node {
	if len(i.Children)-1 >= 2 {
		return i.Children[2]
	}

	return nil
}

func parseElseIfOrElse(t *Tokens) Node {
	if !t.hasNext() {
		return nil
	}

	var expectElse = t.peek()
	if expectElse == "" || expectElse != "else" {
		return nil
	}
	t.next()
	nextStr := t.peek()

	if nextStr == "{" {
		return parseBlock(t)
		//condition must be true, it is an default branch
	} else if nextStr == "if" {
		// match (condition) {body}
		ifNode := parseIf(t)
		return ifNode
	}
	/*
		if() {

		}else if(){

		}else {

		}
	*/
	return nil

}
