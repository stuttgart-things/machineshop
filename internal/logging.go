/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package internal

import (
	"os"
)

func CreateFileLogger(filepath string) (filewWiter *os.File) {

	filewWiter, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		panic(err)
	}

	return
}
