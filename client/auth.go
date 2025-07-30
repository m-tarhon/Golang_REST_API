package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
)

func AuthReq(username, password, port string) error {
	client := GetAuthClient(username, password)
	req, err := http.NewRequest("GET", fmt.Sprintf("https://localhost:%s/users", port), nil)

	if err != nil {
		return err
	}

	// adds Basic Auth header
	auth := username + ":" + password
	encoded := base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", "Basic "+encoded)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// If auth fails, return an error
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return fmt.Errorf("authentication failed: %s", resp.Status)
	}
	// If auth succeeds, return nil
	return nil
}

func Helper4Auth(method, url, username, password string, body io.Reader) (*http.Response, error) {
    req, err := http.NewRequest(method, url, body)
        if err != nil{
            return nil, err 
        }

        if method==http.MethodPost || method == http.MethodDelete{
            req.Header.Set("Content-Type", "application/json")
        }

		return SharedClient.Do(req)
}
