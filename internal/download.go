/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package internal

import (
	"io"
	"net/http"
	"os"
)

func DownloadFile(url string) (string, error) {

	// HTTP GET REQUEST
	resp, err := http.Get(url)
	if err != nil {
		log.Error("ERROR DOWNLOADING THE FILE: ", err)
		os.Exit(3)
	}
	defer resp.Body.Close()

	// CHECK HTTP STATUS
	if resp.StatusCode != http.StatusOK {
		log.Error("ERROR DOWNLOADING THE FILE: STATUS ", resp.StatusCode)
		os.Exit(3)
	}

	// READ ANSWER
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("ERROR DOWNLOADING THE FILE: ", err)
		os.Exit(3)
	}

	return string(body), err
}
