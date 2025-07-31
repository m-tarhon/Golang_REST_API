package main

import (
	"crypto/tls"
	"encoding/base64"
	"net/http"
)

type AuthTransport struct {
	Username string
	Password string	
	Rt http.RoundTripper
}

func (at *AuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	auth := base64.StdEncoding.EncodeToString([]byte(at.Username + ":" + at.Password))
	req.Header.Set("Authorization", "Basic "+auth)
	return at.Rt.RoundTrip(req)
}

func GetAuthClient(username, password string) *http.Client {
	return &http.Client{
		Transport: &AuthTransport{
			Username: username,
			Password: password,
			Rt:       &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}
}
