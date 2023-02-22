package log

import (
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"runtime"
)

func Error(msg string, args ...interface{}) {
	log.Error(fmt.Sprintf(msg, args...)+" | "+getFunctionDetail())
}

func Warn(msg string, args ...interface{}) {
	log.Warn(fmt.Sprintf(msg, args...)+" | "+getFunctionDetail())
}

func Crit(msg string, args ...interface{}) {
	log.Crit(fmt.Sprintf(msg, args...)+" | "+getFunctionDetail())
}

func Info(msg string, args ...interface{}) {
	log.Info(fmt.Sprintf(msg, args...)+" | "+getFunctionDetail())
}

func getFunctionDetail() string {
	pc, file, line, _ := runtime.Caller(2)
	fn := runtime.FuncForPC(pc)
	return fmt.Sprintf("%s:%d (%s)", file, line, fn.Name())
}
