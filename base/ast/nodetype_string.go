// Code generated by "stringer -type NodeType"; DO NOT EDIT.

package ast

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[BLOCK-0]
	_ = x[BINARY_EXPR-1]
	_ = x[UNARY_EXPR-2]
	_ = x[VARIABLE-3]
	_ = x[SCALAR-4]
	_ = x[IF_STMT-5]
	_ = x[WHILE_STMT-6]
	_ = x[FOR_STMT-7]
	_ = x[ASSIGN_STMT-8]
	_ = x[FUNC_STMT-9]
}

const _NodeType_name = "BLOCKBINARY_EXPRUNARY_EXPRVARIABLESCALARIF_STMTWHILE_STMTFOR_STMTASSIGN_STMTFUNC_STMT"

var _NodeType_index = [...]uint8{0, 5, 16, 26, 34, 40, 47, 57, 65, 76, 85}

func (i NodeType) String() string {
	if i >= NodeType(len(_NodeType_index)-1) {
		return "NodeType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _NodeType_name[_NodeType_index[i]:_NodeType_index[i+1]]
}
