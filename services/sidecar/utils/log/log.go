package log

import (
	"fmt"
	"io"
	"runtime"
	"strings"

	"github.com/ethereum/go-ethereum/log"
)

// Re-export geth logger types for convenience.
type (
	Logger = log.Logger
	Lvl    = log.Lvl
)

const (
	LvlCrit  = log.LvlCrit
	LvlError = log.LvlError
	LvlWarn  = log.LvlWarn
	LvlInfo  = log.LvlInfo
	LvlDebug = log.LvlDebug
	LvlTrace = log.LvlTrace
)

func LvlFromString(lvl string) (log.Lvl, error)           { return log.LvlFromString(lvl) }
func NewGlogHandler(h log.Handler) *log.GlogHandler       { return log.NewGlogHandler(h) }
func StreamHandler(w io.Writer, f log.Format) log.Handler { return log.StreamHandler(w, f) }
func TerminalFormat(color bool) log.Format                { return log.TerminalFormat(color) }

// Wraps geth logger to add a prefix to all log messages.
type logger struct{ l Logger }

func New(ctx ...interface{}) log.Logger { return &logger{log.New(ctx...)} }

func (l *logger) New(ctx ...interface{}) log.Logger { return &logger{l.l.New(ctx...)} }
func (l *logger) GetHandler() log.Handler           { return l.l.GetHandler() }
func (l *logger) SetHandler(h log.Handler)          { l.l.SetHandler(h) }

func (l *logger) Trace(msg string, ctx ...interface{}) { l.l.Trace(getLogPrefix()+msg, ctx...) }
func (l *logger) Debug(msg string, ctx ...interface{}) { l.l.Debug(getLogPrefix()+msg, ctx...) }
func (l *logger) Info(msg string, ctx ...interface{})  { l.l.Info(getLogPrefix()+msg, ctx...) }
func (l *logger) Warn(msg string, ctx ...interface{})  { l.l.Warn(getLogPrefix()+msg, ctx...) }
func (l *logger) Error(msg string, ctx ...interface{}) { l.l.Error(getFullLogPrefix()+msg, ctx...) }
func (l *logger) Crit(msg string, ctx ...interface{})  { l.l.Crit(getFullLogPrefix()+msg, ctx...) }

// Prettier error logging.
func (l *logger) Errorf(msg string, err error, args ...interface{}) {
	wrappedErr := fmt.Errorf(getLogPrefix()+msg, err)
	log.Error(wrappedErr.Error(), args...)
}

// Gets the log prefix.
func getLogPrefix() string {
	var (
		// Skip two call frames (from here to the caller of log.X)
		pc, _, line, _ = runtime.Caller(2)
		fn             = runtime.FuncForPC(pc)
		fullFnName     = strings.Split(fn.Name(), ".")
	)
	return fmt.Sprintf("%s:%d | ", fullFnName[len(fullFnName)-1], line)
}

// Gets the full log prefix, including the file name.
// Shouldn't be used everywhere to limit verbosity/string processing.
func getFullLogPrefix() string {
	var (
		// Skip two call frames (from here to the caller of log.X)
		pc, file, line, _ = runtime.Caller(2)
		fn                = runtime.FuncForPC(pc)
		filePath          = strings.Split(file, "/")
		fullFnName        = strings.Split(fn.Name(), ".")
	)
	return fmt.Sprintf("%s:%s:%d | ", filePath[len(filePath)-1], fullFnName[len(fullFnName)-1], line)
}
