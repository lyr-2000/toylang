package ast

// 优先级表
type PriorityTab struct {
}

/*
产生式

e(k) -> e(k+1) ep(k)
 var e = new Expr(); e.left = E(k+1); e.op = op(k); e.right = E(k+1) E_(k)

ep(k) -> op(k) e(k+1) ep(k) | w


*/
