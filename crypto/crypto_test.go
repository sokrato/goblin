package crypto

import (
	"testing"
)

func TestBase64HMAC(t *testing.T) {
	salt := "salt"
	val := "val"
	sec := "sec"
	if Base64HMAC(salt, val, sec) != "y0s5N06mVW4OPd-xRSuB8wo85vs" {
		t.Fail()
	}
}

func TestDjangoSign(t *testing.T) {
	sec := "bwnefv)a@#$#6wi%180vqhqxj35e%+coc#a=-ii%br$j^q@g3o"
	val := "val"
	res := DjangoSign(val, sec)
	if res != "K-Sq-94NC3za1UZ9HKdM6Vs0Udo" {
		t.Fail()
	}
}

func TestDjangoCookie(t *testing.T) {
	sec := "ik1z@sb8mp(%=ja5q+4w4n3dwb2fmy*5_c$zqscw(%hgqw*p&9"
	cookieName := "userid"
	salt := ""
	val := "320:1Xw7su:9O9_X0ncMwf6ohdddE08WJKfOEw"
	if _, err := DjangoGetSignedCookie(cookieName, salt, sec, val); err != nil {
		t.Fail()
	}
}
