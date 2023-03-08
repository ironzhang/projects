package httputils

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVerboseLoadStore(t *testing.T) {
	var i int64
	var v Verbose
	if v != 0 {
		t.Fatalf("v(%d) != 0", v)
	}
	if i = v.Load(); i != 0 {
		t.Fatalf("i(%d) != 0", i)
	}

	v.Store(1)
	if v != 1 {
		t.Fatalf("v(%d) != 1", v)
	}
	if i = v.Load(); i != 1 {
		t.Fatalf("i(%d) != 1", i)
	}
}

func TestVerboseEnabled(t *testing.T) {
	tests := []struct {
		v     Verbose
		input bool
		want  bool
	}{
		{v: Verbose(-1), input: false, want: false},
		{v: Verbose(-1), input: true, want: false},
		{v: Verbose(0), input: false, want: false},
		{v: Verbose(0), input: true, want: true},
		{v: Verbose(1), input: false, want: true},
		{v: Verbose(1), input: true, want: true},
	}
	for i, tt := range tests {
		if got, want := tt.v.enabled(tt.input), tt.want; got != want {
			t.Errorf("tests[%d]: got(%v) != want(%v)", i, got, want)
		}
	}
}

func TestParseTraceId(t *testing.T) {
	{
		casename := "NotSetTraceId"

		h := make(http.Header)
		got, want := parseTraceId(h), h.Get(X_TRACE_ID)
		if got != want {
			t.Errorf("%s got(%v) != want(%v)", casename, got, want)
		} else {
			t.Logf("%s got(%v) == want(%v)", casename, got, want)
		}
	}

	{
		casename := "SetTraceId"

		h := make(http.Header)
		h.Set(X_TRACE_ID, "1")
		got, want := parseTraceId(h), h.Get(X_TRACE_ID)
		if got != want || got != "1" {
			t.Errorf("%s got(%v) != want(%v)", casename, got, want)
		} else {
			t.Logf("%s got(%v) == want(%v)", casename, got, want)
		}
	}
}

func TestVerboseHandler(t *testing.T) {
	h := NewVerboseHandler(nil, nil, nil)

	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	r.Header.Set(X_VERBOSE, "1")

	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
}

func TestVerboseRoundTripper(t *testing.T) {
	rt := NewVerboseRoundTripper(nil, nil, nil)

	s := httptest.NewServer(nil)
	r, err := http.NewRequest("GET", s.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	r.Header.Set(X_VERBOSE, "1")

	if _, err = rt.RoundTrip(r); err != nil {
		t.Fatal(err)
	}
}
