package errs

import (
	"fmt"
	"net/http"

	"github.com/ironzhang/matrix/codes"
	"github.com/ironzhang/matrix/restful"
)

type errorAt struct {
	method string
	err    error
}

func ErrorAt(method string, err error) error {
	return errorAt{method: method, err: err}
}

func (e errorAt) Error() string {
	return fmt.Sprintf("call of %s occur: %v", e.method, e.err)
}

func NotFound(name string, value interface{}) error {
	return restful.Errorf(http.StatusNotFound, codes.NotFound, "%s:%v not found", name, value)
}
