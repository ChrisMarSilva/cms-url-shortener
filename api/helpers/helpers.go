package helpers

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

func EnforceHTTP(url string) (nova_url string) {
	nova_url = url
	if nova_url[:4] != "http" {
		nova_url = "http://" + nova_url
	}
	return
}

func RemoveDomainError(url string) (retorno bool) {

	retorno = true

	domain := os.Getenv("API_DOMAIN")

	if url == domain {
		retorno = false
		return
	}

	newURL := url
	newURL = strings.Replace(newURL, "http://", "", 1)
	newURL = strings.Replace(newURL, "https://", "", 1)
	newURL = strings.Replace(newURL, "www.", "", 1)
	newURL = strings.Split(newURL, "/")[0]

	if url == domain {
		retorno = false
		return
	}

	return
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func GetMD5Hash2(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func GetMD5Hash3(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func GetMD5Hash4(text string) string {
	data := []byte(text)
	return fmt.Sprintf("%x", md5.Sum(data))
}

func GetMD5Hash5(text string) string {
	hasher := md5.New()
	return hex.EncodeToString(hasher.Sum([]byte(text)))
}

func GetMD5Hash6(text string) []byte {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hasher.Sum(nil)
}
