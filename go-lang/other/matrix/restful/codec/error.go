package codec

import (
	"fmt"

	"github.com/ironzhang/matrix/codes"
)

type Error struct {
	Code  int    `json:"code"`
	Desc  string `json:"desc"`
	Cause string `json:"cause,omitempty" xml:",omitempty"`
}

func (e Error) Error() string {
	if e.Cause != "" {
		return fmt.Sprintf("%d: %s (%s)", e.Code, e.Desc, e.Cause)
	}
	return fmt.Sprintf("%d: %s", e.Code, e.Desc)
}

type ErrorCode interface {
	ErrorCode() codes.Code
}

type ErrorCause interface {
	ErrorCause() string
}

func ToError(err error) Error {
	if e, ok := err.(Error); ok {
		return e
	}
	code := codes.Unknown
	if e, ok := err.(ErrorCode); ok {
		code = e.ErrorCode()
	}
	cause := err.Error()
	if e, ok := err.(ErrorCause); ok {
		cause = e.ErrorCause()
	}
	return Error{
		Code:  int(code),
		Desc:  code.String(),
		Cause: cause,
	}
}
