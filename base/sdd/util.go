package sdd

import "toylang/base/ast"

func IsBlock(node Node) bool {

	switch node.(type) {
	case *ast.BlockNode:
		return true
	}
	return false
}

func IsIfBlock(node Node) bool {

	switch node.(type) {
	case *ast.IfNode:
		return true
	}
	return false
}

func IsAssignDeclNode(node Node) bool {

	switch node.(type) {
	case *ast.DeclareStmt:
		{
			return true
		}
	}
	return false
}

func IsAssignExpr(node Node) bool {
	switch node.(type) {
	case *ast.Expr:
		{
			return true
		}
	}
	return false
}
