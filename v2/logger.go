package larkgin

import (
	"context"
	"log"
	"os"

	"github.com/go-lark/lark/v2"
)

// SetLogger set a new logger
func (opt *LarkMiddleware) SetLogger(logger lark.LogWrapper) {
	opt.logger = logger
}

const logPrefix = "[go-lark] "

func initDefaultLogger() lark.LogWrapper {
	// create a default std logger
	logger := stdLogger{
		log.New(os.Stderr, logPrefix, log.LstdFlags),
	}
	return logger
}

type stdLogger struct {
	*log.Logger
}

func (sl stdLogger) Log(_ context.Context, level lark.LogLevel, msg string) {
	sl.Printf("[%s] %s\n", level, msg)
}
