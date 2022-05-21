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

type CallFuncStmt struct {
	BaseNode
}

func (a *BaseNode) Store(t *Token) {
	a.Lexeme = t
}

type MapIndexNode struct {
	BaseNode
	// PropName   Anode // a[1] , a["username"]
	// Variable   *Variable
	// CacheToken *Token
}
type FnParam struct {
	BaseNode
}
type ReturnStmt struct {
	BaseNode
}

type ExprGroups struct {
	BaseNode
}

type BreakFlagStmt struct {
	BaseNode
}

// func (b *MapIndexNode) GetChildren() []Anode {
// 	//a[1] ,是一个变量，不是运算符，不可能有子节点
// 	return []Anode{}
// }

// func (b *MapIndexNode) SetLexeme(*Token) {
// 	panic("illegal state of MapIndexNode ")
// }

// func (b *MapIndexNode) GetLexeme() *Token {
// 	if b.CacheToken != nil {
// 		return b.CacheToken
// 	}
// 	tk := new(Token)
// 	tk.Type = lexer.Variable
// 	tk.Value = fmt.Sprintf("%v[%v]", b.Variable.Lexeme.Value, b.PropName.GetLexeme().Value)
// 	b.CacheToken = tk
// 	return tk
// }
