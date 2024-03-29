package singleflight

import (
	"bytes"
	"fmt"
)

type stackError struct {
	value interface{}
	stack []byte
}

func newStackError(v interface{}, stack []byte) *stackError {
	if line := bytes.IndexByte(stack[:], '\n'); line >= 0 {
		stack = stack[line+1:]
	}
	return &stackError{value: v, stack: stack}
}

func (e *stackError) Error() string {
	return fmt.Sprintf("\n%v\n%s", e.value, e.stack)
}
