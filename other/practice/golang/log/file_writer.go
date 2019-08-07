//file_writer.go
package log

import (
	"errors"
	"fmt"
	"os"
)

type FileWriter struct {
	f         *os.File
	name      string
	sizelimit int
	maxrotate int
	rotate    int
}

func NewFileWriter(name string, sizelimit int, maxrotate int) (*FileWriter, error) {
	w := &FileWriter{name: name, sizelimit: sizelimit, maxrotate: maxrotate}
	err := w.findLastFile()
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (w *FileWriter) Write(p []byte) (n int, err error) {
	if st, err := w.f.Stat(); err == nil && st.Size() >= int64(w.sizelimit) {
		w.rotateFile()
	}
	return w.f.Write(p)
}

func openfile(name string, sizelimit int) (*os.File, error) {
	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	st, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if st.Size() >= int64(sizelimit) {
		return nil, errors.New("file size exceed limit")
	}
	return f, nil
}

func (w *FileWriter) findLastFile() error {
	for {
		if w.rotate >= w.maxrotate {
			break
		}
		filename := fmt.Sprintf("%s.%d", w.name, w.rotate%w.maxrotate)
		if f, err := openfile(filename, w.sizelimit); err == nil {
			w.f = f
			return nil
		}
		w.rotate++
	}

	filename := fmt.Sprintf("%s.%d", w.name, w.rotate%w.maxrotate)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	w.f = f
	return nil
}

func (w *FileWriter) rotateFile() {
	w.rotate++
	filename := fmt.Sprintf("%s.%d", w.name, w.rotate%w.maxrotate)
	f, err := os.Create(filename)
	if err != nil {
		//fmt.Printf("create file failed, filename:%s\n", filename)
		return
	}
	w.f.Close()
	w.f = f
}
