package helpers

import (
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
