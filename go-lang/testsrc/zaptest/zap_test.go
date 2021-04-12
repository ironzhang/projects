package main

import (
	"testing"

	"go.uber.org/zap"
)

func TestZapLogger(t *testing.T) {
	logger := zap.NewExample()
	defer logger.Sync()

	log := logger.Named("test")
	log = log.Named("test2")
	log.Info("hello")
}
