package errs

import "fmt"

type Errno int

func (e Errno) Error() string {
	if s, ok := errors[e]; ok {
		return s
	}
	return fmt.Sprintf("errno %d", e)
}

const (
	ErrUnknown      = Errno(1)
	ErrInvalidParam = Errno(2)
)

var errors = map[Errno]string{}

func init() {
	errors[0] = "nil"
	errors[ErrUnknown] = "unknown"
	errors[ErrInvalidParam] = "invalid parameter"
}
