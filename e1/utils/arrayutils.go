package utils

import (
	"reflect"
	. "e1/log"
)

func CopyArray(dest reflect.Value, src []byte) bool {
	defer func() {
		if x := recover(); x != nil {
			LogError("CopyArray failed:", x)
		}
	}()
	return reflect.Copy(dest.Elem(), reflect.ValueOf(src)) > 0
}
