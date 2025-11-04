package domain

import (
	"errors"
	"fmt"
)

type (
	Kind string // the kind of error (eg: not found)...
	Op   string // for function names to make the stack trace more clear(eg: getUserHandler)
)

const (
	KindNotFound     Kind = "not_found"
	KindInvalid      Kind = "invalid_input"
	KindUnauthorized Kind = "unauthorized"
	KindUnexpected   Kind = "unexpected"
	KindInternal     Kind = "internal"
)

type Error struct {
	Err  error
	Op   Op
	Kind Kind
}

func (e *Error) Error() string {
	return fmt.Sprintf("operation %s failed: %v", e.Op, e.Err)
}

func E(op Op, err error, kind ...Kind) error {
	e := &Error{
		Op:  op,
		Err: err,
	}
	if len(kind) > 0 {
		e.Kind = kind[0]
	} else {
		e.Kind = KindInternal
	}
	return e
}

func GetKind(err error) Kind {
	var e *Error
	if errors.As(err, &e) {
		return e.Kind
	}
	return KindInternal
}
