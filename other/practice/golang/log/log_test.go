package log_test

import (
	"testing"

	"github.com/ironzhang/golang/log"
)

func TestLog(t *testing.T) {
	//log.SetFlag(log.StdFlags | log.MsecFlag | log.ShortFileFlag | log.FuncNameFlag)
	//log.SetFileOutput("./test.log", 1024, 2)
	log.Fatalf("test %d %s", 1, "hello")
	log.Errorf("test %d", 1)
	log.Warnf("test %d", 1)
	log.Infof("test %d", 1)
	log.Debugf("test %d", 1)
}
