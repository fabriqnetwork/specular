package log

import (
	"fmt"

	"github.com/ethereum/go-ethereum/log"
)

func Root() log.Logger { return &logger{log.Root()} }

func Trace(msg string, args ...interface{}) { log.Trace(getLogPrefix()+msg, args...) }
func Debug(msg string, args ...interface{}) { log.Debug(getLogPrefix()+msg, args...) }
func Info(msg string, args ...interface{})  { log.Info(getLogPrefix()+msg, args...) }
func Warn(msg string, args ...interface{})  { log.Warn(getLogPrefix()+msg, args...) }
func Error(msg string, args ...interface{}) { log.Error(getFullLogPrefix()+msg, args...) }
func Crit(msg string, args ...interface{})  { log.Crit(getFullLogPrefix()+msg, args...) }

// Prettier error logging.
func Errorf(msg string, err error, args ...interface{}) {
	wrappedErr := fmt.Errorf(getFullLogPrefix()+msg, err)
	log.Error(wrappedErr.Error(), args...)
}
