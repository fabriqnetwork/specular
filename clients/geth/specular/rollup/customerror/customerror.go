package customerror

import (
	"fmt"
	"runtime"
	"strconv"
)

type CustomError struct {
	err error
	pc  uintptr
}

func Errorf(format string, args ...interface{}) error {
	pc, _, _, _ := runtime.Caller(1)
	return &CustomError{
		err: fmt.Errorf(format, args...),
		pc:  pc,
	}
}

func (ce *CustomError) Error() string {
	return fmt.Sprintf("%s | %s:%d", ce.err.Error(), getFunctionName(ce.pc), getLineNumber(ce.pc))
}

func getFunctionName(pc uintptr) string {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown"
	}
	return fn.Name()
}

func getLineNumber(pc uintptr) int {
	_, file, line, ok := runtime.Caller(0)
	if !ok {
		return 0
	}
	return getLineNumberFromString(file, strconv.Itoa(line))
}

func getLineNumberFromString(file string, lineStr string) int {
	line, err := strconv.Atoi(lineStr)
	if err != nil {
		return 0
	}
	return line
}
