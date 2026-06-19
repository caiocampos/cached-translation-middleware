package util

import "regexp"

var RegexLink = regexp.MustCompile(`(?i)(https?://[.a-z0-9@/\-]+)`)

func GetTextWithoutLinks(text string) string {
	parts := RegexLink.Split(text, -1)
	return parts[0]
}
