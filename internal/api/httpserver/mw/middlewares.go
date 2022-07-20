package mw

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
	"net/http"
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
			sourceCountry, err := ipAPI(request.Context(), address)
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

func ipAPI(ctx context.Context, address string) (string, error) {
	//reqAddr := "https://ipapi.com/" + address + "/country"
	tr := http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{
		Transport:     &tr,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       0,
	}
	reqAddr := "https://ipapi.co/176.111.175.20/country_name" //sample ip
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
