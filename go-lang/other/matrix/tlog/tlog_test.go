package tlog_test

import (
	"context"
	"os"
	"testing"

	"go.uber.org/zap"

	"github.com/ironzhang/matrix/context-value"
	"github.com/ironzhang/matrix/tlog"
	"github.com/ironzhang/matrix/tlog/writers/file"
)

func TestTlog(t *testing.T) {
	tlog.Reset()
	cfg := tlog.Config{
		Level: zap.DebugLevel,
		//Development: true,
		//DisableCaller:     true,
		DisableStacktrace: true,
		DisableStderr:     true,
		EnableFile:        true,
		FileOptions: file.Options{
			Dir: os.TempDir(),
		},
	}
	if _, err := tlog.Init(cfg); err != nil {
		t.Fatal(err)
	}

	log := tlog.Std()
	defer log.Sync()

	log.Debug("debug message", zap.String("function", "TestTlog"))
	log.Info("info message", zap.String("function", "TestTlog"))
	log.Warn("warn message", zap.String("function", "TestTlog"))
	log.Error("error message", zap.String("function", "TestTlog"))
}

func TestWithContext(t *testing.T) {
	tlog.Reset()
	log := tlog.WithContext(context_value.WithTraceId(context.Background(), "4a32dca9-7e2f-4d09-955d-a0103c6f5912"))
	defer log.Sync()

	log.Debug("debug message", zap.String("function", "TestWithContext"))
	log.Info("info message", zap.String("function", "TestWithContext"))
}
