package util

func ToBytes(content any) []byte {
	if bytes, ok := content.([]byte); ok {
		return bytes
	} else {
		return StringToBytes(ToString(content))
	}
}
