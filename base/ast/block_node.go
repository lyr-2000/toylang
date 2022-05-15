package ast

//用來描述抽象语法树

type BlockNode struct {
	AstNode
}

type IfStmt struct {
	AstNode
}

type DeclareStmt struct {
	AstNode
}

type ForStmt struct {
	AstNode
}

type FuncStmt struct {
	AstNode
}

type Factor struct {
	AstNode
}

type Scalar struct {
	Factor
}

type Variable struct {
	Factor
}

type Expr struct {
	AstNode
}

func (a *AstNode) Store(t *Token) {
	a.Lexeme = t
}
