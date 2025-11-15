package util

import (
	"fmt"
	"strconv"
	stringspkg "strings"
	"unsafe"
)

// See:
// https://go101.org/article/unsafe.html
// https://github.com/golang/go/issues/25484
// https://github.com/golang/go/issues/19367
// https://golang.org/src/strings/builder.go

// This conversion *does not* copy data. Note that converting via "(string)([]byte)" *does* copy data.
// Also note that you *should not* change the byte slice after conversion, because Go strings
// are treated as immutable. This would cause a segmentation violation panic.
func BytesToString(bytes []byte) string {
	return unsafe.String(unsafe.SliceData(bytes), len(bytes))

	// return *(*string)(unsafe.Pointer(&bytes))
}

// This conversion *does not* copy data. Note that converting via "([]byte)(string)" *does* copy data.
// Also note that you *should not* change the byte slice after conversion, because Go strings
// are treated as immutable. This would cause a segmentation violation panic.
func StringToBytes(string_ string) (bytes []byte) {
	return unsafe.Slice(unsafe.StringData(string_), len(string_))

	/*
		// StringHeader and SliceHeader have been deprecated in Go 1.21
		stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&string_))
		sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
		sliceHeader.Data = stringHeader.Data
		sliceHeader.Cap = stringHeader.Len
		sliceHeader.Len = stringHeader.Len
		return
	*/
}

// Converts any value to a string.
//
// Will use the [fmt.Stringer] interface if implemented.
//
// Will use the [error] interface if implemented.
//
// nil will become "nil". []byte will become a string.
// Primitive types will use [strconv]. Any other type
// will use fmt.Sprintf("%+v").
func ToString(value any) string {
	if value == nil {
		return "nil"
	}
	switch value_ := value.(type) {
	case string:
		return value_
	case []byte:
		return string(value_)
	case bool:
		return strconv.FormatBool(value_)
	case int64:
		return strconv.FormatInt(value_, 10)
	case int32:
		return strconv.FormatInt(int64(value_), 10)
	case int8:
		return strconv.FormatInt(int64(value_), 10)
	case int:
		return strconv.FormatInt(int64(value_), 10)
	case uint64:
		return strconv.FormatUint(value_, 10)
	case uint32:
		return strconv.FormatUint(uint64(value_), 10)
	case uint8:
		return strconv.FormatUint(uint64(value_), 10)
	case uint:
		return strconv.FormatUint(uint64(value_), 10)
	case float64:
		return strconv.FormatFloat(value_, 'g', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(value_), 'g', -1, 32)
	case fmt.Stringer:
		return value_.String()
	case error:
		return value_.Error()
	default:
		return fmt.Sprintf("%+v", value_)
	}
}

// Calls [ToString] on all slice elements.
func ToStrings(values []any) []string {
	length := len(values)
	if length == 0 {
		return nil
	}

	strings := make([]string, length)
	for index, value := range values {
		strings[index] = ToString(value)
	}

	return strings
}

func JoinQuote(strings []string, separator string) string {
	var builder stringspkg.Builder

	ultimateIndex := len(strings) - 1
	for index, value := range strings {
		builder.WriteString(strconv.Quote(value))
		if index != ultimateIndex {
			builder.WriteString(separator)
		}
	}

	return builder.String()
}

func JoinQuoteL(strings []string, separator string, lastSeparator string, coupleSeparator string) string {
	var builder stringspkg.Builder

	if len(strings) == 2 {
		builder.WriteString(strconv.Quote(strings[0]))
		builder.WriteString(coupleSeparator)
		builder.WriteString(strconv.Quote(strings[1]))
	} else {
		ultimateIndex := len(strings) - 1
		penultimateIndex := ultimateIndex - 1

		for index, value := range strings {
			builder.WriteString(strconv.Quote(value))
			if index != ultimateIndex {
				if index == penultimateIndex {
					builder.WriteString(lastSeparator)
				} else {
					builder.WriteString(separator)
				}
			}
		}
	}

	return builder.String()
}
