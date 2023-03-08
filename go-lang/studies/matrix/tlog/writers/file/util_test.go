package file

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func FileExist(name string) bool {
	_, err := os.Stat(name)
	if err == nil {
		return true
	}
	return false
}

func TestFileName(t *testing.T) {
	timestamp, err := time.Parse("2006-01-02_15:04:05.000000", "2016-01-02_15:04:05.123456")
	if err != nil {
		t.Fatal(err)
	}
	got, want := fileName("TestLogName", timestamp), "TestLogName.2016-01-02_15:04:05.123456.log"
	if got != want {
		t.Errorf("%s != %s", got, want)
	}
}

func TestCreateLogFile(t *testing.T) {
	f, file, link, err := createLogFile("./", "TestCreateLogFile", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Fprintf(f, "%s\n%s\n", file, link)

	if !FileExist(file) {
		t.Errorf("%q file is not exist", file)
	}
	if !FileExist(link) {
		t.Errorf("%q link is not exist", link)
	}

	os.Remove(file)
	os.Remove(link)
}
