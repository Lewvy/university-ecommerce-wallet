package domain

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
