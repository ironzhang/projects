//log.go
package log

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Priority int

const (
	FATAL Priority = iota
	ALERT
	CRIT
	ERROR
	WARN
	NOTICE
	INFO
	DEBUG
)

var pristrs = []string{"[FATAL]", "[ALERT]", "[CRIT]", "[ERROR]", "[WARN]", "[NOTICE]", "[INFO]", "[DEBUG]"}

func (p Priority) String() string {
	return pristrs[p]
}

const (
	DateFlag      = 1 << iota                          //the date: 2015/01/22
	TimeFlag                                           //the time: 13:35:00
	MsecFlag                                           //microsecond resolution: 13:35:00.123123. assumes TimeFlag
	LongFileFlag                                       // full file name and line num: /a/b/c/d.go:20
	ShortFileFlag                                      // final file name and line number: d.go:20 overrides LongfileFlag
	FuncNameFlag                                       // function name: d.go:20 (TestFunc) assumes LongfileFlag
	StdFlags      = DateFlag | TimeFlag | LongFileFlag //standard flags
)

func New() *Logger {
	l, _ := NewLogger("", DEBUG, 0, StdFlags, os.Stdout)
	return l
}

func NewLogger(prefix string, priority Priority, calldepth int, flag int, out io.Writer) (*Logger, error) {
	if out == nil {
		return nil, errors.New("out is nil")
	}

	l := &Logger{
		prefix:    prefix,
		priority:  priority,
		calldepth: calldepth,
		flag:      flag,
		out:       out,
	}
	return l, nil
}

type Logger struct {
	prefix    string
	priority  Priority
	calldepth int
	flag      int
	mu        sync.Mutex
	out       io.Writer
}

func (l *Logger) SetPrefix(prefix string) {
	l.prefix = prefix
}

func (l *Logger) SetPriority(priority Priority) {
	l.priority = priority
}

func (l *Logger) SetCalldepth(calldepth int) {
	l.calldepth = calldepth
}

func (l *Logger) SetFlag(flag int) {
	l.flag = flag
}

func (l *Logger) SetOutput(out io.Writer) {
	if out != nil {
		l.out = out
	}
}

func (l *Logger) Format(priority Priority, calldepth int, s string) string {
	var lineno = 0
	var filename, funcname = "???", "???"
	if l.flag&(LongFileFlag|ShortFileFlag|FuncNameFlag) != 0 {
		if pc, file, line, ok := runtime.Caller(calldepth); ok {
			lineno = line
			filename = file
			if l.flag&FuncNameFlag != 0 {
				if f := runtime.FuncForPC(pc); f != nil {
					funcname = f.Name()
				}
			}
		}
	}

	header := ""
	if l.flag != 0 {
		header = formatHeader(l.flag, time.Now(), filename, lineno, funcname)
	}
	return l.prefix + header + priority.String() + " " + s + "\n"
}

func (l *Logger) Output(priority Priority, calldepth int, s string) {
	if priority <= l.priority {
		outstr := l.Format(priority, calldepth, s)
		l.outputString(outstr)
	}
}

func (l *Logger) Outputf(priority Priority, calldepth int, format string, v ...interface{}) {
	if priority <= l.priority {
		s := fmt.Sprintf(format, v...)
		outstr := l.Format(priority, calldepth, s)
		l.outputString(outstr)
	}
}

func (l *Logger) Printf(priority Priority, format string, v ...interface{}) {
	l.Outputf(priority, l.calldepth+3, format, v...)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Outputf(FATAL, l.calldepth+3, format, v...)
}

func (l *Logger) Alertf(format string, v ...interface{}) {
	l.Outputf(ALERT, l.calldepth+3, format, v...)
}

func (l *Logger) Critf(format string, v ...interface{}) {
	l.Outputf(CRIT, l.calldepth+3, format, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Outputf(ERROR, l.calldepth+3, format, v...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.Outputf(WARN, l.calldepth+3, format, v...)
}

func (l *Logger) Noticef(format string, v ...interface{}) {
	l.Outputf(NOTICE, l.calldepth+3, format, v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.Outputf(INFO, l.calldepth+3, format, v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.Outputf(DEBUG, l.calldepth+3, format, v...)
}

func (l *Logger) outputString(s string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out.Write([]byte(s))
}

func writespace(buf *bytes.Buffer, needspace bool) {
	if needspace {
		buf.WriteByte(' ')
	}
}

func shortfilename(name string) string {
	return name[strings.LastIndex(name, "/")+1:]
}

func shortfuncname(name string) string {
	return name[strings.LastIndex(name, ".")+1:]
}

func formatHeader(flag int, t time.Time, file string, line int, funcname string) string {
	var buf bytes.Buffer
	buf.WriteByte('[')
	needspace := false
	if flag&(DateFlag|TimeFlag|MsecFlag) != 0 {
		if flag&DateFlag != 0 {
			year, month, day := t.Date()
			writespace(&buf, needspace)
			fmt.Fprintf(&buf, "%04d/%02d/%02d", year, month, day)
			needspace = true
		}
		if flag&(TimeFlag|MsecFlag) != 0 {
			hour, min, sec := t.Clock()
			writespace(&buf, needspace)
			fmt.Fprintf(&buf, "%02d:%02d:%02d", hour, min, sec)
			if flag&MsecFlag != 0 {
				fmt.Fprintf(&buf, ".%06d", t.Nanosecond()/1e3)
			}
			needspace = true
		}
	}
	if flag&(LongFileFlag|ShortFileFlag|FuncNameFlag) != 0 {
		if flag&ShortFileFlag != 0 {
			file = shortfilename(file)
		}
		writespace(&buf, needspace)
		fmt.Fprintf(&buf, "%s:%d", file, line)
		if flag&FuncNameFlag != 0 {
			fmt.Fprintf(&buf, " (%s)", shortfuncname(funcname))
		}
		needspace = true
	}
	buf.WriteByte(']')
	return buf.String()
}
