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

func FormatWikiSlug(input string) string {
	words := strings.Fields(strings.ReplaceAll(input, "_", " "))

	for i, w := range words {
		if strings.Contains(w, "-") {
			parts := strings.Split(w, "-")
			for j, p := range parts {
				if len(p) > 0 {
					parts[j] = strings.ToUpper(p[:1]) + strings.ToLower(p[1:])
				}
			}

			words[i] = strings.Join(parts, "-")
		} else if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + strings.ToLower(w[1:])
		}
	}

	return strings.Join(words, "_")
}
