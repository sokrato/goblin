package utils

import (
	"testing"
)

func TestBase64HMAC(t *testing.T) {
	salt := []byte("salt")
	val := []byte("val")
	sec := []byte("sec")
	if Base64HMAC(salt, val, sec) != "y0s5N06mVW4OPd-xRSuB8wo85vs" {
		t.Fail()
	}
}

func TestDjangoSign(t *testing.T) {
	sec := []byte("bwnefv)a@#$#6wi%180vqhqxj35e%+coc#a=-ii%br$j^q@g3o")
	val := []byte("val")
	res := DjangoSign(val, sec)
	if res != "K-Sq-94NC3za1UZ9HKdM6Vs0Udo" {
		t.Fail()
	}
}
