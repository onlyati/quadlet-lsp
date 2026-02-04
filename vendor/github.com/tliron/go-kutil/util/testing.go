package util

import (
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func assertTransform[T comparable](t *testing.T, transform func(T) T, original T, transformed T) {
	if transformed_ := transform(original); transformed_ != transformed {
		t.Errorf("%s: %s", getFunctionName(transform), ToString(transformed_))
	}
}

func getFunctionName(fn any) string {
	if function := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()); function != nil {
		s := strings.Split(function.Name(), ".")
		return s[len(s)-1]
	} else {
		return "<unknown function>"
	}
}
