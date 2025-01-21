/*
Copyright © 2025 Patrick Hermann patrick.hermann@sva.de
*/

package internal

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
)

var (
	contentType = "application/json"
)

// SendToHomerun sends a message to the Homerun service with optional insecure TLS settings.
func SendToHomerun(destination, token string, renderedBody []byte, insecure bool) (answer []byte, resp *http.Response) {
	// CREATE HTTP-Request
	req, err := http.NewRequest("POST", destination, bytes.NewBuffer(renderedBody))
	if err != nil {
		fmt.Println("FAILURE AT CREATING REQUESTS:", err)
		return nil, nil
	}

	// ADD HEADER
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("X-Auth-Token", token)

	// CREATE CUSTOM HTTP CLIENT WITH INSECURE OPTION
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
	}
	client := &http.Client{Transport: tr}

	// SEND REQUEST
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println("ERROR AT SENDING REQUEST:", err)
		return nil, nil
	}
	defer resp.Body.Close()

	// READ THE ANSWER
	answer, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ERROR READING ANSWER:", err)
		return nil, nil
	}

	return answer, resp
}