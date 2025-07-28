package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
)

func AuthReq(username, password, port string) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%s/users", port), nil)

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

	return nil
}
