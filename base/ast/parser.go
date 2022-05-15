package ast

import (
	"fmt"
	"toylang/base/lexer"

	"github.com/xlab/treeprint"
)

type Parser struct {
}

type Token = lexer.Token

type NodeType uint8

//go:generate stringer -type NodeType

const (
	BLOCK       NodeType = iota
	BINARY_EXPR          //1+1 ,binary expr
	UNARY_EXPR           // ++i
	VARIABLE             //变量
	SCALAR               // 标量，
	IF_STMT              //if 语句
	WHILE_STMT           //while

	FOR_STMT    //for循环
	ASSIGN_STMT // = 赋值
	FUNC_STMT   //func 赋值语句

)

type AstNode struct {
	//节点类型 ？
	NodeType NodeType
	//子节点
	Children []Anode
	//父节点
	Parent Anode
	//单元
	Lexeme *Token
}

// type Node = AstNode
type Anode interface {
	GetNodeType() NodeType
	GetLexeme() *Token
	GetChildren() []Anode
	GetParent() Anode
	SetLexeme(t *Token)
	SetNodeType(NodeType)
}

// 设置 node type
func (u *AstNode) GetNodeType() NodeType {
	return u.NodeType
}
func (u *AstNode) SetNodeType(n NodeType) {
	u.NodeType = n

}
func (u *AstNode) GetLexeme() *Token {
	return u.Lexeme
}
func (u *AstNode) GetChildren() []Anode {
	return u.Children
}
func (u *AstNode) GetParent() Anode {
	return u.Parent
}

// 设置 token
func (u *AstNode) SetLexeme(t *Token) {
	u.Lexeme = t
}

type ILexer = lexer.Lexer

// func toStringNode(n Anode) string {
// 	var (
// 		t   string
// 		val string
// 	)
// 	t = fmt.Sprintf("%T,", n)
// 	val = fmt.Sprintf("%+v", n)
// 	return t + val
// }

func toDfsPatternStringNode(n Anode) string {
	// var buf strings.Builder
	tree := treeprint.New()
	// tree.Newb
	bh := tree.AddBranch("root")
	var dfs func(nnode Anode, bh treeprint.Tree)
	dfs = func(nnode Anode, bh treeprint.Tree) {
		if nnode == nil {
			return
		}
		if len(nnode.GetChildren()) == 0 {
			bh.AddNode(fmt.Sprintf("%+v", nnode.GetLexeme()))
			return
		}
		ch := bh.AddBranch(fmt.Sprintf("(%v)%+v", nnode.GetNodeType(), nnode.GetLexeme()))
		// ch.AddNode(fmt.Sprintf("%+v", nnode.GetLexeme()))
		// bh.AddNode(fmt.Sprintf("%+v", nnode))
		for _, c := range nnode.GetChildren() {
			dfs(c, ch)
		}
	}
	dfs(n, bh)
	return tree.String()

}

// func (*Parser) ParseNode(lx ILexer) Anode {
// 	var (
// 		// scalar Scalar
// 		expr Expr
// 	)

// 	//表达式
// 	return expr.AstNode
// }

type PeekTokenIterator struct {
	i      int
	tokens []*Token
}

func (p *PeekTokenIterator) HasNext() bool {
	return p.i < len(p.tokens)
}
func (p *PeekTokenIterator) Next() *Token {
	cur := p.tokens[p.i]
	p.i++
	return cur
}
func (p *PeekTokenIterator) NextMatch(v string) (*Token, error) {
	cur := p.tokens[p.i]
	p.i++
	if cur.Value != v {
		return cur, fmt.Errorf("not an correct token")
	}

	return cur, nil
}
