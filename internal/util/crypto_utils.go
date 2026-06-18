package util

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
)

func GenerateMD5(val string) [16]byte {
	valBytes := []byte(val)
	return md5.Sum(valBytes)
}

func GenerateSHA256(val string) [32]byte {
	valBytes := []byte(val)
	return sha256.Sum256(valBytes)
}

func byteArrayAsBase64(val []byte) string {
	return base64.StdEncoding.EncodeToString(val)
}

func GenerateMD5String(val string) string {
	valHash := GenerateMD5(val)
	return byteArrayAsBase64(valHash[:])
}

func GenerateSHA256String(val string) string {
	valHash := GenerateSHA256(val)
	return byteArrayAsBase64(valHash[:])
}
