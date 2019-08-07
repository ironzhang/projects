package httputils

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/ironzhang/matrix/tlog"
	"github.com/ironzhang/matrix/uuid"
)

func NewVerbose(i int64) *Verbose {
	v := Verbose(i)
	return &v
}

type Verbose int64

func (v *Verbose) Load() int64 {
	return atomic.LoadInt64((*int64)(v))
}

func (v *Verbose) Store(i int64) {
	atomic.StoreInt64((*int64)(v), i)
}

func (v *Verbose) enabled(verbose bool) bool {
	i := v.Load()
	if i == 0 {
		return verbose
	} else if i < 0 {
		return false
	} else {
		return true
	}
}

func parseTraceId(h http.Header) string {
	traceId := h.Get(X_TRACE_ID)
	if traceId == "" {
		traceId = uuid.New().String()
		h.Set(X_TRACE_ID, traceId)
	}
	return traceId
}

func NewVerboseHandler(verbose *Verbose, writer io.Writer, handler http.Handler) *VerboseHandler {
	if verbose == nil {
		verbose = NewVerbose(0)
	}
	if writer == nil {
		writer = os.Stdout
	}
	if handler == nil {
		handler = http.DefaultServeMux
	}
	return &VerboseHandler{
		verbose: verbose,
		writer:  writer,
		handler: handler,
	}
}

type VerboseHandler struct {
	verbose *Verbose
	writer  io.Writer
	handler http.Handler
}

func (h *VerboseHandler) Verbose() *Verbose {
	return h.verbose
}

func (h *VerboseHandler) printRequest(clientId, traceId string, r *http.Request) {
	b, err := httputil.DumpRequest(r, true)
	if err != nil {
		tlog.StdSugar().Errorw("dump request", "clientId", clientId, "traceId", traceId, "error", err)
		return
	}
	fmt.Fprintf(h.writer, "%s\tPROTO\tserver request\t{%q: %q, %q: %q}\n%s\n", time.Now(), "clientId", clientId, "traceId", traceId, b)
}

func (h *VerboseHandler) printResponse(clientId, traceId string, r *ResponseDumper) {
	b := r.Dump(true)
	fmt.Fprintf(h.writer, "%s\tPROTO\tserver response\t{%q: %q, %q: %q}\n%s\n", time.Now(), "clientId", clientId, "traceId", traceId, b)
}

func (h *VerboseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v, _ := strconv.ParseBool(r.Header.Get(X_VERBOSE))
	if h.verbose.enabled(v) {
		clientId := r.Header.Get(X_CLIENT_ID)
		traceId := parseTraceId(r.Header)

		h.printRequest(clientId, traceId, r)
		d := NewResponseDumper(w, r)
		defer h.printResponse(clientId, traceId, d)
		w = d
	}
	h.handler.ServeHTTP(w, r)
}

func NewVerboseRoundTripper(verbose *Verbose, writer io.Writer, transport http.RoundTripper) *VerboseRoundTripper {
	if verbose == nil {
		verbose = NewVerbose(0)
	}
	if writer == nil {
		writer = os.Stdout
	}
	if transport == nil {
		transport = http.DefaultTransport
	}
	return &VerboseRoundTripper{
		verbose:   verbose,
		writer:    writer,
		transport: transport,
	}
}

type VerboseRoundTripper struct {
	verbose   *Verbose
	writer    io.Writer
	transport http.RoundTripper
}

func (rt *VerboseRoundTripper) Verbose() *Verbose {
	return rt.verbose
}

func (rt *VerboseRoundTripper) printRequest(clientId, traceId string, r *http.Request) {
	b, err := httputil.DumpRequestOut(r, true)
	if err != nil {
		tlog.StdSugar().Errorw("dump request out", "clientId", clientId, "traceId", traceId, "error", err)
		return
	}
	fmt.Fprintf(rt.writer, "%s\tPROTO\tclient request\t{%q: %q, %q: %q}\n%s\n", time.Now(), "clientId", clientId, "traceId", traceId, b)
}

func (rt *VerboseRoundTripper) printResponse(clientId, traceId string, r *http.Response) {
	b, err := httputil.DumpResponse(r, true)
	if err != nil {
		tlog.StdSugar().Errorw("dump response", "clientId", clientId, "traceId", traceId, "error", err)
		return
	}
	fmt.Fprintf(rt.writer, "%s\tPROTO\tclient response\t{%q: %q, %q: %q}\n%s\n", time.Now(), "clientId", clientId, "traceId", traceId, b)
}

func (rt *VerboseRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	var clientId, traceId string

	v, _ := strconv.ParseBool(r.Header.Get(X_VERBOSE))
	verbose := rt.verbose.enabled(v)
	if verbose {
		clientId = r.Header.Get(X_CLIENT_ID)
		traceId = parseTraceId(r.Header)
		rt.printRequest(clientId, traceId, r)
	}
	resp, err := rt.transport.RoundTrip(r)
	if err != nil {
		return resp, err
	}
	if verbose {
		rt.printResponse(clientId, traceId, resp)
	}
	return resp, err
}
