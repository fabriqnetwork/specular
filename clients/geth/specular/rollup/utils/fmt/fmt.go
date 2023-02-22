package fmt

import (
	"fmt"
	"runtime"

	"github.com/go-errors/errors"
)

type Error struct {
	err  string
	file string
	line int
	name string
}

func Errorf(format string, args ...interface{}) *Error {
	pc, file, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	return &Error{
		err:  errors.Errorf(format, args...).ErrorStack(),
		file: file,
		line: line,
		name: fn.Name(),
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s | %s:%d (%s)", e.err, e.file, e.line, e.name)
}
