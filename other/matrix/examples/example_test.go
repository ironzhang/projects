package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func postPing(t *testing.T, h http.Handler) {
	r, _ := http.NewRequest("POST", "/ping", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("post /ping result is unexpect, Code[%d]", w.Code)
	}

	buf, _ := ioutil.ReadAll(w.Body)
	t.Logf("Status: %v", w.Code)
	t.Logf("Header: %v", w.Header())
	t.Logf("Body: %s", string(buf))
}

func getPing(t *testing.T, h http.Handler) {
	r, _ := http.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("get /ping result is unexpect, Code[%d]", w.Code)
	}

	buf, _ := ioutil.ReadAll(w.Body)
	t.Logf("Status: %v", w.Code)
	t.Logf("Header: %v", w.Header())
	t.Logf("Body: %s", string(buf))
}

func doTest(t *testing.T, h http.Handler, method, url string, status int) {
	r, _ := http.NewRequest(method, url, nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != status {
		t.Errorf("do [%s %s] test is unexpect, status[%d]", method, url, status)
	}

	buf, _ := ioutil.ReadAll(w.Body)
	t.Logf("Status: %v", w.Code)
	t.Logf("Header: %v", w.Header())
	t.Logf("Body: %s", string(buf))
}

func TestPing(t *testing.T) {
	h := newHandler()

	var testcases = []struct {
		method string
		url    string
		status int
	}{
		{"POST", "/ping", http.StatusOK},
		{"GET", "/ping", http.StatusMethodNotAllowed},
	}
	for _, tc := range testcases {
		doTest(t, h, tc.method, tc.url, tc.status)
	}
}
