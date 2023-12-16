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

type CategorizedError[T comparable] struct {
	Cat T
	Err error
}

func NewCategorizedError[T comparable](cat T, err error) *CategorizedError[T] {
	return &CategorizedError[T]{Cat: cat, Err: err}
}

func (e *CategorizedError[T]) Unwrap() error { return e.Err }
func (e *CategorizedError[T]) Category() T   { return e.Cat }
func (e *CategorizedError[T]) Error() string { return e.Err.Error() }

// Shallow comparison: categories must be equal.
func (e *CategorizedError[T]) Is(target error) bool {
	err, ok := target.(*CategorizedError[T])
	if !ok {
		return false
	}
	return e.Cat == err.Cat
}
