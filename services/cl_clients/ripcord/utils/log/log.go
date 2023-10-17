package log

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/ethereum/go-ethereum/log"
)

type Logger = log.Logger

// Wraps geth logger to add a prefix to all log messages.
type logger struct{ l Logger }

func New(ctx ...interface{}) log.Logger { return &logger{log.New(ctx...)} }

func (l *logger) New(ctx ...interface{}) log.Logger { return &logger{l.l.New(ctx...)} }
func (l *logger) GetHandler() log.Handler           { return l.l.GetHandler() }
func (l *logger) SetHandler(h log.Handler)          { l.l.SetHandler(h) }

func (l *logger) Trace(msg string, ctx ...interface{}) { l.l.Trace(getLogPrefix()+" | "+msg, ctx...) }
func (l *logger) Debug(msg string, ctx ...interface{}) { l.l.Debug(getLogPrefix()+" | "+msg, ctx...) }
func (l *logger) Info(msg string, ctx ...interface{})  { l.l.Info(getLogPrefix()+" | "+msg, ctx...) }
func (l *logger) Warn(msg string, ctx ...interface{})  { l.l.Warn(getLogPrefix()+" | "+msg, ctx...) }
func (l *logger) Error(msg string, ctx ...interface{}) { l.l.Error(getLogPrefix()+" | "+msg, ctx...) }
func (l *logger) Crit(msg string, ctx ...interface{})  { l.l.Crit(getLogPrefix()+" | "+msg, ctx...) }

// Prettier error logging.
func (l *logger) Errorf(msg string, err error, args ...interface{}) {
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
