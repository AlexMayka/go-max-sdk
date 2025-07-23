package utils

import (
	"fmt"
	"reflect"
)

func ShouldOmit(value reflect.Value) bool {
	if value.Kind() == reflect.Ptr && value.IsNil() {
		return true
	}
	if value.Kind() != reflect.Ptr && value.IsZero() {
		return true
	}
	return false
}

func GetInterfaceValue(value reflect.Value) interface{} {
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return nil
		}
		return value.Elem().Interface()
	}
	return value.Interface()
}

func GetStringValue(value reflect.Value) string {
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return ""
		}
		return fmt.Sprintf("%v", value.Elem().Interface())
	}
	return fmt.Sprintf("%v", value.Interface())
}