package util

import (
	"encoding/base64"
)

func ToBase64(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}

func FromBase64(b64 string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(b64)
}
