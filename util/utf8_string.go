package util

import (
	"strings"
	"unicode/utf8"
)

// TruncateToMaxChars safely truncates a string to a maximum number of characters
// without breaking UTF-8 encoding. If the string exceeds maxChars,
// it returns a substring containing at most maxChars complete runes.
func TruncateToMaxChars(s string, maxChars int) string {
	if utf8.RuneCountInString(s) <= maxChars {
		return s
	}
	// Convert to runes for proper character handling
	runes := []rune(s)
	return string(runes[:maxChars])
}

func SanitizeUTF8(s string) string {
	if !utf8.ValidString(s) {
		return strings.ToValidUTF8(s, "?")
	}
	return s
}

func FormatUTF8(s string) string {
	s = strings.ReplaceAll(s, "\x00", "") // for postgres db do not accept 0x00 as string char
	return SanitizeUTF8(s)
}
