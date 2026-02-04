package util

// Returns true if value is an int64, int32, int16, int8, int,
// uint64, uint32, uint16, uint8, uint, float64, or float32.
func IsNumber(value any) bool {
	switch value.(type) {
	case int64, int32, int16, int8, int, uint64, uint32, uint16, uint8, uint, float64, float32:
		return true
	default:
		return false
	}
}

// Returns true if value is an int64, int32, int16, int8, int,
// uint64, uint32, uint16, uint8, or uint.
func IsInteger(value any) bool {
	switch value.(type) {
	case int64, int32, int16, int8, int, uint64, uint32, uint16, uint8, uint:
		return true
	}
	return false
}

// Returns true if value is a float64 or float32.
func IsFloat(value any) bool {
	switch value.(type) {
	case float64, float32:
		return true
	}
	return false
}

// Converts any number type to int64.
//
// Supported types are int64, int32, int16, int8, int,
// uint64, uint32, uint16, uint8, uint,
// float64, and float32.
//
// Precision may be lost.
func ToInt64(value any) (int64, bool) {
	switch value_ := value.(type) {
	case int64:
		return value_, true
	case int32:
		return int64(value_), true
	case int16:
		return int64(value_), true
	case int8:
		return int64(value_), true
	case int:
		return int64(value_), true
	case uint64:
		return int64(value_), true
	case uint32:
		return int64(value_), true
	case uint16:
		return int64(value_), true
	case uint8:
		return int64(value_), true
	case uint:
		return int64(value_), true
	case float64:
		return int64(value_), true
	case float32:
		return int64(value_), true
	default:
		return 0, false
	}
}

// Converts any number type to uint64.
//
// Support types are int64, int32, int16, int8, int,
// uint64, uint32, uint16, uint8, uint,
// float64, and float32.
//
// Precision may be lost.
func ToUInt64(value any) (uint64, bool) {
	switch value_ := value.(type) {
	case uint64:
		return value_, true
	case uint32:
		return uint64(value_), true
	case uint16:
		return uint64(value_), true
	case uint8:
		return uint64(value_), true
	case uint:
		return uint64(value_), true
	case int64:
		return uint64(value_), true
	case int32:
		return uint64(value_), true
	case int16:
		return uint64(value_), true
	case int8:
		return uint64(value_), true
	case int:
		return uint64(value_), true
	case float64:
		return uint64(value_), true
	case float32:
		return uint64(value_), true
	default:
		return 0, false
	}
}

// Converts any number type to float64.
//
// Support types are int64, int32, int16, int8, int,
// uint64, uint32, uint16, uint8, uint,
// float64, and float32.
//
// Precision may be lost.
func ToFloat64(value any) (float64, bool) {
	switch value_ := value.(type) {
	case float64:
		return value_, true
	case float32:
		return float64(value_), true
	case int64:
		return float64(value_), true
	case int32:
		return float64(value_), true
	case int16:
		return float64(value_), true
	case int8:
		return float64(value_), true
	case int:
		return float64(value_), true
	case uint64:
		return float64(value_), true
	case uint32:
		return float64(value_), true
	case uint16:
		return float64(value_), true
	case uint8:
		return float64(value_), true
	case uint:
		return float64(value_), true
	default:
		return 0.0, false
	}
}
