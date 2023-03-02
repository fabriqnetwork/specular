package fmt

import (
	"fmt"
	"runtime"

	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/errors"
)

type Error struct {
	err  *errors.Error
	file string
	line int
	fn   *runtime.Func
}

func Errorf(format string, args ...interface{}) *Error {
	pc, file, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	return &Error{
		err:  errors.Errorf(format, args...),
		file: file,
		line: line,
		fn:   fn,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s | %s:%d (%s)", e.err.ErrorStack(), e.file, e.line, e.fn.Name())
}
