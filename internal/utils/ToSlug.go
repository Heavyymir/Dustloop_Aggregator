package utils

import(
	"strings"
	"unicode"
)

// To slug converts strings like Iron Man or Chun li to url-friendly varients

func ToSlug(input string) string {
	clean := strings.ToLower(strings.TrimSpace(input))

	var b strings.Builder
	for _, r := range clean {
		switch {
			case unicode.IsLetter(r) || unicode.IsDigit(r):
				b.WriteRune(r)
			case r == ' ' || r == '_' || r == '-':
				// Collapse spaces, hyphens or underscores into a single underscore
				if b.Len() > 0 && !strings.HasSuffix(b.String(), "-") {
					b.WriteRune('_')
				}
		}
	}

	return b.String()
}
