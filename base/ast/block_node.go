package ast

//用來描述抽象语法树

type BlockNode struct {
	BaseNode
}

type IfStmt struct {
	BaseNode
}

type DeclareStmt struct {
	BaseNode
}

type ForStmt struct {
	BaseNode
}

type FuncStmt struct {
	BaseNode
}

type Factor struct {
	BaseNode
}

type Scalar struct {
	Factor
}

type Variable struct {
	Factor
}

type Expr struct {
	BaseNode
	NodeType
}

func (a *BaseNode) Store(t *Token) {
	a.Lexeme = t
}
