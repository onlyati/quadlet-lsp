package util

import (
	"github.com/iancoleman/strcase"
)

// See: https://en.wikipedia.org/wiki/Letter_case#Special_case_styles

// See:
// https://github.com/iancoleman/strcase
// https://github.com/stoewer/go-strcase

func ToDromedaryCase(name string) string {
	return strcase.ToLowerCamel(name)

	/*
		runes := []rune(camelCase)
		length := len(runes)

		if (length > 0) && unicode.IsUpper(runes[0]) { // sanity check
			if (length > 1) && unicode.IsUpper(runes[1]) {
				// If the second rune is also uppercase we'll keep the name as is
				return camelCase
			}
			runes_ := make([]rune, 1, length-1)
			runes_[0] = unicode.ToLower(runes[0])
			return string(append(runes_, runes[1:]...))
		} else {
			return camelCase
		}
	*/
}

func ToSnakeCase(name string) string {
	return strcase.ToSnake(name)
}

func ToKebabCase(name string) string {
	return strcase.ToKebab(name)
}
