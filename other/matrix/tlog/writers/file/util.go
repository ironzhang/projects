package file

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func linkName(program string) string {
	return fmt.Sprintf("%s.log", program)
}

func fileName(program string, t time.Time) string {
	return fmt.Sprintf("%s.%s.log", program, t.Format("2006-01-02_15:04:05.000000"))
}

func createLogFile(dir, program string, t time.Time) (f *os.File, file, link string, err error) {
	name := fileName(program, t)
	file = filepath.Join(dir, name)
	link = filepath.Join(dir, linkName(program))

	f, err = os.Create(file)
	if err == nil {
		os.Remove(link)
		os.Symlink(name, link)
		return f, file, link, nil
	}
	return nil, "", "", err
}
