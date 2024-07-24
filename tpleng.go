package tpleng

import (
	"fmt"
	"unicode/utf8"
)

func Parse[V any](template string, vars map[string]V) (string, error) {

	parsed := ""
	parsing := false

	for i, stride := 0, 0; i < len(template); i += stride {
		phrase, width, isParsing := decodePhraseInString(template[i:], vars, parsing)
		stride = width
		parsing = isParsing
		fmt.Printf("parse: %s %d %t\n", phrase, stride, isParsing)
		parsed += phrase

	}

	return parsed, nil
}

func decodePhraseInString[V any](s string, vars map[string]V, parsing bool) (string, int, bool) {

	if parsing {
		return decodePhraseInStringParsing(s, vars)
	} else {
		return decodePhraseInStringEchoing(s)
	}
}

func decodePhraseInStringEchoing(s string) (string, int, bool) {

	decoded := ""
	stride := 0
	parsing := false

	for i, w := 0, 0; i < len(s); i += w {
		runeValue, width := utf8.DecodeRuneInString(s[i:])
		w = width
		stride += width

		switch runeValue {
		case '{':
			lookahead, width := utf8.DecodeRuneInString(s[i+width:])
			parsing = lookahead == '{'
			if parsing {
				w += width
				stride += width
			} else {
				decoded += string(runeValue)
			}
		default:
			decoded += string(runeValue)
		}

		fmt.Printf("echoing: %s %d %t\n", decoded, stride, parsing)

		if parsing {
			break
		}

	}

	return decoded, stride, parsing
}

func decodePhraseInStringParsing[V any](s string, vars map[string]V) (string, int, bool) {

	expression := ""
	stride := 0
	parsing := true

	for i, w := 0, 0; i < len(s); i += w {
		runeValue, width := utf8.DecodeRuneInString(s[i:])
		w = width
		stride += width

		switch runeValue {
		case ' ', '\n', '\t':
			// Do nothing (ignore white space).
		case '}':
			lookahead, _ := utf8.DecodeRuneInString(s[i+width:])
			parsing = lookahead == '}'
			// if !parsing {
			// 	w += width
			// 	stride += width
			// }
		default:
			expression += string(runeValue)
		}

		if !parsing {
			break
		}
	}

	return decodeExpression(expression, vars), stride, parsing
}

func decodeExpression[V any](s string, vars map[string]V) string {
	result := ""
	stride := 0

	for i, w := 0, 0; i < len(s); i += w {
		runeValue, width := utf8.DecodeRuneInString(s[i:])
		w = width
		stride += width

		switch runeValue {
		case '.':
			key, width := decodeWord(s[i+width:])
			w += width

			result += fmt.Sprintf("%v", vars[key])
		}
	}

	return result
}

func decodeWord(s string) (string, int) {
	word := ""
	stride := 0

	for i, w := 0, 0; i < len(s); i += w {
		runeValue, width := utf8.DecodeRuneInString(s[i:])
		w = width
		stride += width

		switch runeValue {
		case '|', ' ', '\n', '\t':
		default:
			word += string(runeValue)
		}
	}

	return word, stride
}
