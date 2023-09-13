package log

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/ethereum/go-ethereum/log"
)

func Trace(msg string, args ...interface{}) {
	log.Trace(getLogPrefix()+" | "+msg, args...)
}

func Debug(msg string, args ...interface{}) {
	log.Debug(getLogPrefix()+" | "+msg, args...)
}

func Info(msg string, args ...interface{}) {
	log.Info(getLogPrefix()+" | "+msg, args...)
}

func Warn(msg string, args ...interface{}) {
	log.Warn(getLogPrefix()+" | "+msg, args...)
}

func Error(msg string, args ...interface{}) {
	log.Error(getLogPrefix()+" | "+msg, args...)
}

func Crit(msg string, args ...interface{}) {
	log.Crit(getLogPrefix()+" | "+msg, args...)
}

// Prettier error logging.
func Errorf(msg string, err error, args ...interface{}) {
	wrappedErr := fmt.Errorf(getLogPrefix()+" | "+msg, err)
	log.Error(wrappedErr.Error(), args...)
}

func getLogPrefix() string {
	// Skip two call frames (from here to the caller of log.X)
	pc, _, line, _ := runtime.Caller(2)
	fn := runtime.FuncForPC(pc)
	fullFnName := strings.Split(fn.Name(), ".")
	return fmt.Sprintf("%s:%d", fullFnName[len(fullFnName)-1], line)
}
