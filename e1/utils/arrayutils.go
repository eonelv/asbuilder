package utils

import (
	"reflect"
	"fmt"
)

func CopyArray(dest reflect.Value, src []byte) bool {
	defer func() {
		if x := recover(); x != nil {
			fmt.Println("CopyArray failed:", x)
		}
	}()
	return reflect.Copy(dest.Elem(), reflect.ValueOf(src)) > 0
}
