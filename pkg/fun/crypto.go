package fun

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"log"
	"strings"
)

//md5
func MD5(text string) string {
	hasher := md5.New()
	_, err := hasher.Write([]byte(text))
	if err != nil {
		log.Println("error:", err)
		return ""
	} else {
		return hex.EncodeToString(hasher.Sum(nil))
	}
}

//base64_encode
func Base64Encode(sDec string) string {
	sEnc := base64.StdEncoding.EncodeToString([]byte(sDec))
	return sEnc
}

//base64_decode
func Base64Decode(sEnc string) string {
	sDec, err := base64.StdEncoding.DecodeString(sEnc)
	if err != nil {
		log.Println("error:", err)
		return ""
	} else {
		return string(sDec)
	}

}

//urlencode
func UrlEncode(uDec string) string {
	uEnc := base64.URLEncoding.EncodeToString([]byte(uDec))
	return uEnc
}

//urldecode
func UrlDecode(uEnc string) string {
	uDec, err := base64.URLEncoding.DecodeString(uEnc)
	if err != nil {
		log.Println("error:", err)
		return ""
	} else {
		return string(uDec)
	}
}

//rawurlencode
func RawUrlEncode(str string) string {
	return strings.Replace(UrlEncode(str), "+", "%20", -1)
}

//rawurldecode
func RawUrlDecode(str string) string {
	return UrlDecode(strings.Replace(str, "%20", "+", -1))
}
