package errs

type opError struct {
	svc string
	op  string
	err error
}
