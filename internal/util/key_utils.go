package util

import (
	"fmt"
)

func BuildKey(source string, target string, text string) string {
	var textHash string
	if len(text) < 50 {
		textHash = GenerateMD5String(text)
	} else {
		textHash = GenerateSHA256String(text)
	}
	return fmt.Sprintf("translation:%s:%s:%x", source, target, textHash)
}
