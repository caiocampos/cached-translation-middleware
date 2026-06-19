package util

import (
	"strings"
)

func BuildKey(source string, target string, text string) string {
	var textHash string
	if len(text) < 50 {
		textHash = GenerateMD5String(text)
	} else {
		textHash = GenerateSHA256String(text)
	}
	return strings.Join([]string{"tr", source, target, textHash}, ":")
}
