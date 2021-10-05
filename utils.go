package behold

import (
	"unicode"
)

func exportedField(field string) bool {
	if field == "" {
		// use Key field
		return true
	}

	for _, r := range field {
		// test first rune
		return unicode.IsUpper(r)
	}
	return false
}
