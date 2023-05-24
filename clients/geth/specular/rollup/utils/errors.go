package utils

import (
	"errors"
	"fmt"
)

func AsAny(err error, targets ...any) bool {
	for _, target := range targets {
		if errors.As(err, target) {
			return true
		}
	}
	return false
}

type L1ReorgDetectedError struct{ Msg string }

type L2ReorgDetectedError struct{ Msg string }

func (e L1ReorgDetectedError) Error() string {
	return fmt.Sprintf("L1 reorg detected: %s", e.Msg)
}

func (e L2ReorgDetectedError) Error() string {
	return fmt.Sprintf("L2 reorg detected: %s", e.Msg)
}
