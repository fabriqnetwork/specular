package utils

import (
	"errors"
)

func AsAny(err error, targets ...any) bool {
	for _, target := range targets {
		if errors.As(err, target) {
			return true
		}
	}
	return false
}
