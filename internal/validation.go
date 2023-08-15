/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package internal

import (
	"log"
	"os/user"

	"golang.org/x/sys/unix"
)

func ValidateOSUser() (username string) {

	user, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}

	username = user.Username

	return

}

func CheckForUnixWritePermissions(filePath string) bool {

	return unix.Access(filePath, unix.W_OK) == nil
}
