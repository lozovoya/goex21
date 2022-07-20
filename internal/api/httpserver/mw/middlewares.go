package mw

import (
	"crypto/tls"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func CountryChecker(country string, lg *logrus.Entry) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			address, err := realIP(request.RemoteAddr)
			if err != nil {
				lg.Error(err)
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			sourceCountry, err := ipAPI(address)
			if err != nil {
				lg.Error(err)
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			if sourceCountry != country {
				http.Error(writer, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}
			next.ServeHTTP(writer, request)
		})
	}

}

func realIP(remoteAddr string) (string, error) {
	ip, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return "", fmt.Errorf("mw.realIP: %w", err)
	}
	return ip, nil
}

func ipAPI(address string) (string, error) {
	tr := http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{
		Transport:     &tr,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       5 * time.Second,
	}
	//reqAddr := "https://ipapi.com/" + address + "/country"
	reqAddr := "https://ipapi.co/176.111.175.20/country_name" //sample Cyprus ip
	response, err := client.Get(reqAddr)
	if err != nil {
		return "", fmt.Errorf("mw.ipAPI: %w", err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("mw.ipAPI: %w", err)
	}
	return string(body), nil
}
