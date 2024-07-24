package tpleng

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func Parse(template string, vars map[string]string) (string, error) {

	parsed := ""
	parsing := false

	for i, stride := 0, 0; i < len(template); i += stride {
		phrase, width, isParsing := decodePhraseInString(template[i:], vars, parsing)
		stride = width
		parsing = isParsing
		parsed += phrase

	}

	return parsed, nil
}

func decodePhraseInString(s string, vars map[string]string, parsing bool) (string, int, bool) {

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

		if parsing {
			break
		}

	}

	return decoded, stride, parsing
}

func decodePhraseInStringParsing(s string, vars map[string]string) (string, int, bool) {

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

func decodeExpression(s string, vars map[string]string) string {
	result := ""
	stride := 0

	fmt.Printf("expression: %s\n", s)

	for i, w := 0, 0; i < len(s); i += w {
		runeValue, width := utf8.DecodeRuneInString(s[i:])
		w = width
		stride += width

		switch runeValue {
		case ' ', '\n', '\t':
			// Do nothing (ignore white space).
		case '.':
			key, width := decodeWord(s[i+width:])
			w += width
			stride += width

			fmt.Printf("dot: %s\n", key)

			result = vars[key]

			if len(result) >= 4 && result[:2] == "{{" && result[len(result)-2:] == "}}" {
				result = decodeExpression(result[2:len(result)-2], vars)
			}
		case '|':
			transform, width := decodeWord(s[i+width:])
			fmt.Printf("transform: %s\n", transform)
			w += width

			result = performTransform(transform, result)
		}
	}

	return result
}

func performTransform(transform string, input string) string {

	output := input

	switch transform {
	case "uppercase":
		output = strings.ToUpper(input)
	case "lowercase":
		output = strings.ToLower(input)
	}

	return output
}

func decodeWord(s string) (string, int) {
	word := ""
	stride := 0
	ended := false

	fmt.Printf("word pre: %s\n", s)

	for i, w := 0, 0; i < len(s); i += w {
		runeValue, width := utf8.DecodeRuneInString(s[i:])

		fmt.Printf("%s\n", string(runeValue))

		switch runeValue {
		case ' ', '\n', '\t':
			w = width
			stride += width
			// Do nothing (ignore white space)
		case '|', '}':
			ended = true
		default:
			w = width
			stride += width
			word += string(runeValue)
		}

		if ended {
			break
		}
	}

	fmt.Printf("word post: %s, %d\n", word, stride)

	return word, stride
}
