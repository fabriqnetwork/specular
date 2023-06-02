package utils

import "reflect"

// Note: limit use of reflection to this file.

func GetField[T any](val any, fieldName string) T {
	indirect := reflect.Indirect(reflect.ValueOf(val))
	return indirect.FieldByName(fieldName).Interface().(T)
}

func SetField[T any](val any, fieldName string, fieldValue T) {
	indirect := reflect.Indirect(reflect.ValueOf(val))
	indirect.FieldByName(fieldName).Set(reflect.ValueOf(fieldValue))
}
