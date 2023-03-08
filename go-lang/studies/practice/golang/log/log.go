package log

import (
	"io"
)

var std *Logger

func init() {
	std = New()
	std.SetCalldepth(1)
	std.SetFlag(StdFlags | MsecFlag | ShortFileFlag | FuncNameFlag)
}

func SetPrefix(prefix string) {
	std.SetPrefix(prefix)
}

func SetPriority(priority Priority) {
	std.SetPriority(priority)
}

func SetCalldepth(calldepth int) {
	std.SetCalldepth(calldepth)
}

func SetFlag(flag int) {
	std.SetFlag(flag)
}

func SetOutput(out io.Writer) {
	std.SetOutput(out)
}

func SetFileOutput(name string, sizelimit, maxrotate int) error {
	out, err := NewFileWriter(name, sizelimit, maxrotate)
	if err != nil {
		return err
	}
	std.SetOutput(out)
	return nil
}

func Fatalf(format string, v ...interface{}) {
	std.Fatalf(format, v...)
}

func Errorf(format string, v ...interface{}) {
	std.Errorf(format, v...)
}

func Warnf(format string, v ...interface{}) {
	std.Warnf(format, v...)
}

func Noticef(format string, v ...interface{}) {
	std.Noticef(format, v...)
}

func Infof(format string, v ...interface{}) {
	std.Infof(format, v...)
}

func Debugf(format string, v ...interface{}) {
	std.Debugf(format, v...)
}
