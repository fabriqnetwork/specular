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
	err := errors.Errorf(format, args...)
	return &wrapError{
		err:  err,
		msg: fmt.Sprintf("%s", err.ErrorStack()),
	}
}

func (e *wrapError) Error() string {
	return e.msg
}

