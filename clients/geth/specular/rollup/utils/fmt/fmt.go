package fmt

import (
	"fmt"
	"runtime"
)

type Error struct {
	err error
	file string
	line int
	name string
}

func Errorf(format string, args ...interface{}) error {
	pc, file, line, _ := runtime.Caller(2)
	fn := runtime.FuncForPC(pc)
	return &Error{
		err: fmt.Errorf(format, args...),
		file: file,
		line: line,
		name: fn.Name(),
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s | %s:%d (%s)", e.err.Error(),e.file, e.line, e.name)
}