package compile

import (
	"path"
	"strings"
	"unicode"
)

var knownAbbreviations = map[string]string{
	"id":   "ID",
	"url":  "URL",
	"uri":  "URI",
	"api":  "API",
	"http": "HTTP",
	"uuid": "UUID",
	"sql":  "SQL",
	"html": "HTML",
	"css":  "CSS",
	"json": "JSON",
	"xml":  "XML",
	"ip":   "IP",
	"tcp":  "TCP",
	"udp":  "UDP",
	"dns":  "DNS",
	"ssh":  "SSH",
	"ssl":  "SSL",
	"tls":  "TLS",
	"cpu":  "CPU",
	"gpu":  "GPU",
	"ram":  "RAM",
	"os":   "OS",
	"db":   "DB",
	"ui":   "UI",
	"ux":   "UX",
	"fk":   "FK",
	"pk":   "PK",
}

// PropertyNameToFieldName converts a snake_case or camelCase OpenAPI
// property name to PascalCase suitable for a Morphe field name.
func PropertyNameToFieldName(name string) string {
	words := splitWords(name)
	var result strings.Builder
	for _, word := range words {
		lower := strings.ToLower(word)
		if abbr, ok := knownAbbreviations[lower]; ok {
			result.WriteString(abbr)
		} else {
			result.WriteString(strings.ToUpper(word[:1]) + strings.ToLower(word[1:]))
		}
	}
	return result.String()
}

// SchemaNameToFileName converts a PascalCase schema name to a
// snake_case file name (without extension).
func SchemaNameToFileName(name string) string {
	words := splitWords(name)
	for i, w := range words {
		words[i] = strings.ToLower(w)
	}
	return strings.Join(words, "_")
}

// RefToSchemaName extracts the schema name from a $ref string
// like "#/components/schemas/Customer".
func RefToSchemaName(ref string) string {
	return path.Base(ref)
}

func splitWords(s string) []string {
	var words []string
	var current strings.Builder

	flushCurrent := func() {
		if current.Len() > 0 {
			words = append(words, current.String())
			current.Reset()
		}
	}

	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		switch {
		case r == '_' || r == '-' || r == ' ':
			flushCurrent()
		case unicode.IsUpper(r):
			if current.Len() > 0 {
				// Split on upper following lower: "fooBar" → "foo", "Bar"
				if i > 0 && unicode.IsLower(runes[i-1]) {
					flushCurrent()
				} else if i+1 < len(runes) && unicode.IsLower(runes[i+1]) && current.Len() > 1 {
					// Split acronym from next word: "HTTPServer" → "HTTP", "Server"
					flushCurrent()
				}
			}
			current.WriteRune(r)
		default:
			current.WriteRune(r)
		}
	}
	flushCurrent()
	return words
}
