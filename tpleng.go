package tpleng

import (
	"unicode/utf8"
)

func Parse[V any](template string, vars map[string]V) (string, error) {

	parsed := ""

	for i, w := 0, 0; i < len(template); i += w {
		runeValue, width := utf8.DecodeRuneInString(template[i:])
		w = width
		switch runeValue {
		default:
			parsed += string(runeValue)
		}
	}

	return parsed, nil
}
