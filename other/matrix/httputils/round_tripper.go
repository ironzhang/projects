package httputils

import (
	"io"
	"net/http"
)

type nopReadCloser struct {
	io.Reader
}

func (p *nopReadCloser) Close() error {
	return nil
}

type NopRoundTripper struct{}

func (rt NopRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       &nopReadCloser{r.Body},
	}, nil
}
