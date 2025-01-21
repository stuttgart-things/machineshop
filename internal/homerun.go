/*
Copyright © 2025 Patrick Hermann patrick.hermann@sva.de
*/

package internal

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

var (
	contentType = "application/json"
)

// Send a message to the homerun service
func SendToHomerun(destination, token string, renderedBody []byte) (answer []byte, resp *http.Response) {
	// CREATE HTTP-Request
	req, err := http.NewRequest("POST", destination, bytes.NewBuffer([]byte(renderedBody)))
	if err != nil {
		fmt.Println("FAIULURE AT CREATING REQUESTS:", err)
		return
	}

	// ADD HEADER
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("X-Auth-Token", token)

	// CREATE HTTP-Client + SEND REQUEST
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println("ERROR AT SENDING REQUEST:", err)
		return
	}
	defer resp.Body.Close()

	// READ THE ANSWER
	answer, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading answer:", err)
		return
	}

	return answer, resp
}
