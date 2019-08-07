package restful

import (
	"fmt"
	"net/http"

	"github.com/ironzhang/matrix/codes"
)

type HTTPStatus interface {
	HTTPStatus() int
}

func NewError(status int, code codes.Code) Error {
	return Error{Status: status, Code: code}
}

func Errorf(status int, code codes.Code, format string, a ...interface{}) Error {
	return Error{
		Status: status,
		Code:   code,
		Cause:  fmt.Sprintf(format, a...),
	}
}

type Error struct {
	Status int
	Code   codes.Code
	Cause  string
}

func (e Error) HTTPStatus() int {
	return e.Status
}

func (e Error) ErrorCode() codes.Code {
	return e.Code
}

func (e Error) ErrorCause() string {
	return e.Cause
}

func (e Error) Error() string {
	if e.Cause != "" {
		return fmt.Sprintf("[%d: %s] %d: %s (%s)", e.Status, http.StatusText(e.Status), e.Code, e.Code.String(), e.Cause)
	}
	return fmt.Sprintf("[%d: %s] %d: %s", e.Status, http.StatusText(e.Status), e.Code, e.Code.String())
}
