package utils

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"hash"
	"strings"
)

const (
	DjangoSalt = "django.core.signing.Signersigner"
)

var (
	b64Pads = map[byte]string{
		0: "",
		1: "=",
		2: "==",
		3: "===",
	}
)

func SaltedHMAC(key_salt, value, secret []byte) hash.Hash {
	key := sha1.Sum(append(key_salt, secret...))
	hs := hmac.New(sha1.New, key[:])
	hs.Write(value)
	return hs
}

func b64_encode(bs []byte) string {
	s := base64.URLEncoding.EncodeToString(bs)
	return strings.TrimRight(s, "=")
}

func b64_decode(s string) ([]byte, error) {
	n := len(s) % 4
	pad := b64Pads[byte(n)]
	return base64.URLEncoding.DecodeString(s + pad)
}

func Base64HMAC(salt, value, secret []byte) string {
	hs := SaltedHMAC(salt, value, secret)
	return b64_encode(hs.Sum(nil))
}

func DjangoSign(value, sec []byte) string {
	return Base64HMAC([]byte(DjangoSalt), value, sec)
}
