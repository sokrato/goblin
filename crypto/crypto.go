package crypto

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"hash"
	"strings"
)

const (
	DjangoSalt            = "django.core.signing.Signersigner"
	DjangoCookieSecPrefix = "django.http.cookies"
)

var (
	b64Pads = map[byte]string{
		0: "",
		1: "=",
		2: "==",
		3: "===",
	}
)

var (
	ErrBadSign = errors.New("签名不正确")
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

func Base64HMAC(salt, value, secret string) string {
	hs := SaltedHMAC([]byte(salt), []byte(value), []byte(secret))
	return b64_encode(hs.Sum(nil))
}

func DjangoSign(value, sec string) string {
	return Base64HMAC(DjangoSalt, value, sec)
}

func DjangoSignCookie(cookieName, cookieValue, salt, sec string) string {
	salt = cookieName + salt + "signer"
	sec = DjangoCookieSecPrefix + sec
	return Base64HMAC(salt, cookieValue, sec)
}

func DjangoGetSignedCookie(cookieName, salt, sec, value string) (string, error) {
	parts := strings.Split(value, ":")
	var val, sig string
	if len(parts) == 2 { // no timestamp
		val, sig = parts[0], parts[1]
	} else if len(parts) == 3 {
		val = parts[0] + ":" + parts[1]
		sig = parts[2]
	} else {
		return "", ErrBadSign
	}

	if DjangoSignCookie(cookieName, val, salt, sec) != sig {
		return "", ErrBadSign
	}
	return val, nil
}
