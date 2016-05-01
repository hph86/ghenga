package db

import "unicode"

// ToSnakeCase returns a snake_case version of the CamelCase string s.
func ToSnakeCase(s string) string {
	var (
		runes      = []rune(s)
		out        []rune
		prev, next rune
	)

	for i, r := range runes {
		if i > 0 {
			prev = runes[i-1]
		}
		if i < len(runes)-1 {
			next = runes[i+1]
		}

		if len(out) > 0 && unicode.IsUpper(r) {
			if unicode.IsLower(prev) || unicode.IsLower(next) {
				out = append(out, '_')
			}
		}

		out = append(out, unicode.ToLower(r))
	}

	return string(out)
}
