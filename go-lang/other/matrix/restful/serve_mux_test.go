package restful

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"

	"github.com/ironzhang/matrix/codes"
	"github.com/ironzhang/matrix/restful/codec"
	"github.com/ironzhang/matrix/tlog"
)

func ServeHTTP(h http.Handler, method, path string, b []byte) (*httptest.ResponseRecorder, error) {
	r, err := http.NewRequest(method, path, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	//r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w, nil
}

type Args struct {
	A, B int
}

type Reply struct {
	C int
}

type Arith int

func (t *Arith) Add(ctx context.Context, values url.Values, args Args, reply *Reply) error {
	reply.C = args.A + args.B
	return nil
}

func (t *Arith) Sub(ctx context.Context, values url.Values, args Args, reply *Reply) error {
	reply.C = args.A - args.B
	return nil
}

func (t *Arith) Mul(ctx context.Context, values url.Values, args *Args, reply *Reply) error {
	reply.C = args.A * args.B
	return nil
}

func (t *Arith) Div(ctx context.Context, values url.Values, args Args, reply *Reply) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	reply.C = args.A / args.B
	return nil
}

func NewArithServeMux() (m *ServeMux, err error) {
	var a Arith
	m = NewServeMux(nil)
	if err = m.Add("POST", "/add", a.Add); err != nil {
		return nil, err
	}
	if err = m.Add("POST", "/sub", a.Sub); err != nil {
		return nil, err
	}
	if err = m.Add("POST", "/mul", a.Mul); err != nil {
		return nil, err
	}
	if err = m.Add("POST", "/div", a.Div); err != nil {
		return nil, err
	}
	return m, nil
}

func CallArith(m *ServeMux, method, path string, a, b int) (c int, err error) {
	args := Args{a, b}

	var buf bytes.Buffer
	if err = m.codec.Encode(&buf, args); err != nil {
		return 0, err
	}

	r, err := ServeHTTP(m, method, path, buf.Bytes())
	if err != nil {
		return 0, err
	}
	if r.Code != http.StatusOK {
		var e codec.Error
		if err = m.codec.DecodeError(r.Body, &e); err != nil {
			return 0, err
		}
		return 0, Errorf(r.Code, codes.Code(e.Code), e.Cause)
	}

	var reply Reply
	if err = m.codec.Decode(r.Body, &reply); err != nil {
		return 0, err
	}
	return reply.C, nil
}

func TestArithServeMuxReturnNil(t *testing.T) {
	tlog.Reset()

	m, err := NewArithServeMux()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		method string
		path   string
		a      int
		b      int
		c      int
	}{
		{method: "POST", path: "/add", a: 1, b: 2, c: 3},
		{method: "Post", path: "/sub", a: 1, b: 2, c: -1},
		{method: "post", path: "/mul", a: 1, b: 2, c: 2},
		{method: "post", path: "/div", a: 1, b: 2, c: 0},
	}

	for i, tt := range tests {
		c, err := CallArith(m, tt.method, tt.path, tt.a, tt.b)
		if err != nil {
			t.Errorf("tests[%d]: call arith: %v", i, err)
			continue
		}
		if c != tt.c {
			t.Errorf("tests[%d]: got(%v) != want(%v)", i, c, tt.c)
		}
	}
}

func TestArithServeMuxReturnErr(t *testing.T) {
	tlog.Init(tlog.Config{DisableStderr: true})

	m, err := NewArithServeMux()
	if err != nil {
		t.Fatal(err)
	}

	_, err = CallArith(m, "post", "/div", 1, 0)
	if err == nil {
		t.Fatal("CallArith, expect error buf not")
	}
	e := err.(Error)
	if status := http.StatusInternalServerError; e.Status != status {
		t.Errorf("status: %v != %v", e.Status, status)
	}
	if code := codes.Unknown; e.Code != code {
		t.Errorf("code: %v != %v", e.Code, code)
	}
	if cause := "divide by zero"; e.Cause != cause {
		t.Errorf("cause: %q != %q", e.Cause, cause)
	}
}

type Tester struct{}

func (t Tester) ReturnNil(ctx context.Context, values url.Values, req interface{}, resp interface{}) error {
	return nil
}

func (t Tester) ReturnDecodeFailError(ctx context.Context, values url.Values, req int, resp interface{}) error {
	return nil
}

func (t Tester) ReturnUnknownError(ctx context.Context, values url.Values, req interface{}, resp interface{}) error {
	return errors.New("unknown")
}

func (t Tester) ReturnInternalError(ctx context.Context, values url.Values, req interface{}, resp interface{}) error {
	return Errorf(http.StatusInternalServerError, codes.Internal, "internal error")
}

func (t Tester) ReturnInvalidParamError(ctx context.Context, values url.Values, req interface{}, resp interface{}) error {
	return NewError(http.StatusBadRequest, codes.InvalidParam)
}

func (t Tester) ReturnOutOfRangeErrorWithCause(ctx context.Context, values url.Values, req interface{}, resp interface{}) error {
	return Errorf(http.StatusInternalServerError, codes.OutOfRange, "out of range")
}

func NewTesterServeMux() (m *ServeMux, err error) {
	var t Tester
	m = NewServeMux(nil)
	if err = m.Add("GET", "/ReturnNil", t.ReturnNil); err != nil {
		return nil, err
	}
	if err = m.Add("GET", "/ReturnDecodeFailError", t.ReturnDecodeFailError); err != nil {
		return nil, err
	}
	if err = m.Add("GET", "/ReturnUnknownError", t.ReturnUnknownError); err != nil {
		return nil, err
	}
	if err = m.Add("GET", "/ReturnInternalError", t.ReturnInternalError); err != nil {
		return nil, err
	}
	if err = m.Add("GET", "/ReturnInvalidParamError", t.ReturnInvalidParamError); err != nil {
		return nil, err
	}
	if err = m.Add("GET", "/ReturnOutOfRangeErrorWithCause", t.ReturnOutOfRangeErrorWithCause); err != nil {
		return nil, err
	}
	return m, nil
}

func TestServeMuxReturnNil(t *testing.T) {
	tlog.Reset()

	m, err := NewTesterServeMux()
	if err != nil {
		t.Fatal(err)
	}

	r, err := ServeHTTP(m, "GET", "/ReturnNil", nil)
	if err != nil {
		t.Fatal(err)
	}
	if status := http.StatusOK; r.Code != status {
		t.Fatalf("status: %v != %v", r.Code, status)
	}
}

func TestServeMuxReturnErr(t *testing.T) {
	tlog.Init(tlog.Config{DisableStderr: true})

	m, err := NewTesterServeMux()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		method string
		path   string
		status int
		code   codes.Code
		cause  string
	}{
		{"POST", "/NotFound", http.StatusNotFound, codes.NotFound, "page(/NotFound) not found"},
		{"POST", "/ReturnNil", http.StatusMethodNotAllowed, codes.NotAllowed, "method(POST) not allowed"},
		{"GET", "/ReturnDecodeFailError", http.StatusBadRequest, codes.DecodeFail, "EOF"},
		{"GET", "/ReturnUnknownError", http.StatusInternalServerError, codes.Unknown, "unknown"},
		{"GET", "/ReturnInternalError", http.StatusInternalServerError, codes.Internal, "internal error"},
		{"GET", "/ReturnInvalidParamError", http.StatusBadRequest, codes.InvalidParam, ""},
		{"GET", "/ReturnOutOfRangeErrorWithCause", http.StatusInternalServerError, codes.OutOfRange, "out of range"},
	}
	for i, tt := range tests {
		r, err := ServeHTTP(m, tt.method, tt.path, nil)
		if err != nil {
			t.Fatal(err)
		}
		var e codec.Error
		if m.codec.DecodeError(r.Body, &e); err != nil {
			t.Fatal(err)
		}

		if got, want := r.Code, tt.status; got != want {
			t.Errorf("tests[%d]: http status: %v != %v", i, got, want)
		}
		if got, want := e.Code, int(tt.code); got != want {
			t.Errorf("tests[%d]: error code: %v != %v", i, got, want)
		}
		if got, want := e.Desc, tt.code.String(); got != want {
			t.Errorf("tests[%d]: error desc: %q != %q", i, got, want)
		}
		if got, want := e.Cause, tt.cause; got != want {
			t.Errorf("tests[%d]: error cause: %q != %q", i, got, want)
		}
	}
}

func BenchmarkServerSerial(b *testing.B) {
	tlog.Reset()

	m, err := NewArithServeMux()
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		_, err = CallArith(m, "POST", "/add", 0, 0)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkServerParallel(b *testing.B) {
	tlog.Reset()

	m, err := NewArithServeMux()
	if err != nil {
		b.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(b.N)
	for i := 0; i < b.N; i++ {
		go func() {
			_, err := CallArith(m, "POST", "/add", 0, 0)
			if err != nil {
				b.Error(err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
