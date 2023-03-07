package fmt

import (
	"fmt"
	"runtime"
	
	"github.com/specularl2/specular/clients/geth/specular/rollup/utils/errors"
)

type wrapError struct {
	err  *errors.Error
	msg string
}

func Errorf(format string, args ...interface{}) *wrapError {
	pc, _, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	err := errors.Errorf(format, args...)
	return &wrapError{
		err:  err,
		msg: fmt.Sprintf("%s | %s:%d", err.ErrorStack(), fn.Name(), line),
	}
}

func (e *wrapError) Error() string {
	return e.msg
}

