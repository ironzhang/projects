package filewatch

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func WriteTestFile(path string, content string) error {
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, []byte(content), 0644)
}

type Loader struct {
	path string
	hit  int
}

func (p *Loader) Reload(path string) bool {
	p.hit++
	p.path = path
	return false
}

func TestFileObserver(t *testing.T) {
	err := WriteTestFile("./testdata/test_file_observer/exist.txt", "")
	if err != nil {
		t.Fatalf("write test file: %v", err)
	}

	tests := []struct {
		path string
		hit  int
	}{
		{
			path: "./testdata/test_file_observer/exist.txt",
			hit:  1,
		},
		{
			path: "./testdata/test_file_observer/not_exist.txt",
			hit:  0,
		},
	}
	for _, tt := range tests {
		ld := Loader{}
		fo := &fileObserver{
			path:      tt.path,
			watchFunc: ld.Reload,
		}
		fo.observe()
		if got, want := ld.hit, tt.hit; got != want {
			t.Fatalf("%q hit: got %v, want %v", tt.path, got, want)
		}
		if ld.hit == 0 {
			continue
		}
		if got, want := ld.path, tt.path; got != want {
			t.Fatalf("%q path: got %v, want %v", tt.path, got, want)
		}
		fo.observe()
		if got, want := ld.hit, tt.hit; got != want {
			t.Fatalf("%q second hit: got %v, want %v", tt.path, got, want)
		}
	}
}
