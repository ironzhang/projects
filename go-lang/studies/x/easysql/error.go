package easysql

import "fmt"

type opError struct {
	op  string
	err error
}

func (e *opError) Error() string {
	return fmt.Sprintf("%s: %v", e.op, e.err)
}

func opErrorf(op string, format string, a ...interface{}) *opError {
	return &opError{op: op, err: fmt.Errorf(format, a...)}
}
