package tlog

import (
	"context"
	"os"

	"github.com/ironzhang/matrix/context-value"
	"github.com/ironzhang/matrix/tlog/writers/file"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var level zap.AtomicLevel
var std *zap.Logger
var sugar *zap.SugaredLogger

func init() {
	if err := Reset(); err != nil {
		panic(err)
	}
}

type Config struct {
	Level             zapcore.Level
	Development       bool
	DisableCaller     bool
	DisableStacktrace bool
	DisableStderr     bool
	EnableFile        bool
	FileOptions       file.Options
}

func (cfg Config) openSinks() (zapcore.WriteSyncer, error) {
	var writers []zapcore.WriteSyncer

	if !cfg.DisableStderr {
		writers = append(writers, os.Stderr)
	}

	if cfg.EnableFile {
		f, err := file.Open(cfg.FileOptions)
		if err != nil {
			return nil, err
		}
		writers = append(writers, f)
	}

	return zap.CombineWriteSyncers(writers...), nil
}

func (cfg Config) buildOptions(sink zapcore.WriteSyncer) []zap.Option {
	opts := []zap.Option{zap.ErrorOutput(sink)}

	if cfg.Development {
		opts = append(opts, zap.Development())
	}

	if !cfg.DisableCaller {
		opts = append(opts, zap.AddCaller())
	}

	stackLevel := zapcore.ErrorLevel
	if cfg.Development {
		stackLevel = zapcore.WarnLevel
	}
	if !cfg.DisableStacktrace {
		opts = append(opts, zap.AddStacktrace(stackLevel))
	}

	return opts
}

func Reset() error {
	cfg := zap.NewDevelopmentConfig()
	logger, err := cfg.Build()
	if err != nil {
		return err
	}

	level = cfg.Level
	std = logger
	sugar = std.Sugar()
	return nil
}

func Init(cfg Config) (*zap.Logger, error) {
	sink, err := cfg.openSinks()
	if err != nil {
		return nil, err
	}
	level = zap.NewAtomicLevelAt(cfg.Level)
	opts := cfg.buildOptions(sink)
	enc := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	core := zapcore.NewCore(enc, sink, level)
	std = zap.New(core, opts...)
	sugar = std.Sugar()
	return std, nil
}

func Std() *zap.Logger {
	return std
}

func StdSugar() *zap.SugaredLogger {
	return sugar
}

func Level() zap.AtomicLevel {
	return level
}

func WithContext(ctx context.Context) *zap.Logger {
	log := std
	if traceId := context_value.ParseTraceId(ctx); traceId != "" {
		log = log.With(zap.String("traceId", traceId))
	}
	return log
}
