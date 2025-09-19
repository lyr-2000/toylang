package evaluator

import (
	"fmt"

	"github.com/spf13/cast"
)

type InlineFunc func(f *Interpreter, args ...*RefValue) any

var (
	inlineFunc = map[string]InlineFunc{
		"print":   printLn,
		"throw":   throwErr,
		"recover": recoverErr,
		"fatal":   fatalErr,
		"array":   NewArray,
		"map":     NewMap,
		"set":     exitCodeSet,
		"exit":    exitCodeSet,
		"append":  unwrapLang(AppendSlice),
		"remove":  unwrapLang(RemoveValueAtIndex),
		"typeof":  unwrapLang(StrTypeOf),
	}
)

func StrTypeOf(f *Interpreter, args ...any) any {
	return fmt.Sprintf("%T", args[0])
}

func AppendSlice(f *Interpreter, args ...any) any {
	ele, ok := args[0].([]any)
	if !ok {
		return args[0]
	}
	ele = append(ele, args[1:]...)
	return ele
}

func RemoveValueAtIndex(f *Interpreter, args ...any) any {
	ele, isSlice := args[0].([]any)
	if isSlice {
		idx := cast.ToInt(args[1])
		if idx >= 0 && idx < len(ele) {
			ele = append(ele[:idx], ele[idx+1:]...)
		}
		return ele
	}
	mapval, ok := args[0].(map[string]any)
	if ok && mapval != nil {
		delete(mapval, cast.ToString(args[1]))
	}
	return args[0]
}
