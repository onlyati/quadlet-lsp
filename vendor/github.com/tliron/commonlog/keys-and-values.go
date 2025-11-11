package commonlog

import (
	"fmt"

	"github.com/tliron/kutil/util"
)

// Calls [Message.Set] on a provided sequence of key-value pairs.
// Thus keysAndValues must have an even length.
//
// Non-string keys are converted to strings using [util.ToString].
func SetMessageKeysAndValues(message Message, keysAndValues ...any) {
	length := len(keysAndValues)

	if length == 0 {
		return
	}

	if length%2 != 0 {
		panic(fmt.Sprintf("CommonLog message keysAndValues does not have an even number of arguments: %d", length))
	}

	for index := 0; index < length; index += 2 {
		key := util.ToString(keysAndValues[index])
		value := keysAndValues[index+1]
		message.Set(key, value)
	}
}

func GetKeyValue(key any, keysAndValues ...any) (any, bool) {
	length := len(keysAndValues)

	if length%2 != 0 {
		panic(fmt.Sprintf("CommonLog keysAndValues does not have an even number of arguments: %d", length))
	}

	for index := 0; index < length; index += 2 {
		if keysAndValues[index] == key {
			return keysAndValues[index+1], true
		}
	}

	return "", false
}

func SetKeyValue(key any, value any, keysAndValues ...any) ([]any, bool) {
	length := len(keysAndValues)

	if length%2 != 0 {
		panic(fmt.Sprintf("CommonLog keysAndValues does not have an even number of arguments: %d", length))
	}

	for index := 0; index < length; index += 2 {
		if keysAndValues[index] == key {
			keysAndValues[index+1] = value
			return keysAndValues, false
		}
	}

	return append(keysAndValues, key, value), true
}

func MergeKeysAndValues(toKeysAndValues []any, fromKeysAndValues []any) ([]any, bool) {
	length := len(fromKeysAndValues)

	if length%2 != 0 {
		panic(fmt.Sprintf("CommonLog keysAndValues does not have an even number of arguments: %d", length))
	}

	var changed bool
	var changed_ bool

	for index := 0; index < length; index += 2 {
		toKeysAndValues, changed_ = SetKeyValue(fromKeysAndValues[index], fromKeysAndValues[index+1], toKeysAndValues...)
		if changed_ {
			changed = true
		}
	}

	return toKeysAndValues, changed
}
