package logs

import (
	"runtime"
	"strings"
)

func GetFunctionDetail() string {
    pc, _, _, _ := runtime.Caller(2)
	nameFull := runtime.FuncForPC(pc).Name()    
    splitInput := strings.Split(nameFull, "/")
    return splitInput[len(splitInput)-1]
}