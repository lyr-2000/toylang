// Code generated by "stringer -type TokenType"; DO NOT EDIT.

package lexer

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Keyword-1]
	_ = x[Variable-2]
	_ = x[Operator-3]
	_ = x[Brackets-4]
	_ = x[String-5]
	_ = x[Char-6]
	_ = x[Number-7]
	_ = x[Boolean-8]
	_ = x[EOF-9]
	_ = x[Illegal-10]
}

const _TokenType_name = "KeywordVariableOperatorBracketsStringCharNumberBooleanEOFIllegal"

var _TokenType_index = [...]uint8{0, 7, 15, 23, 31, 37, 41, 47, 54, 57, 64}

func (i TokenType) String() string {
	i -= 1
	if i < 0 || i >= TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}
