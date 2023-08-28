/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package internal

import (
	"fmt"
	"log"
	"os/user"

	"golang.org/x/sys/unix"
	"gopkg.in/yaml.v3"
)

// SHOULD BE MOVED TO STHINGS-BASE!

func ValidateOSUser() (username string) {

	user, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}

	username = user.Username

	return

}

// SHOULD BE MOVED TO STHINGS-BASE!

func CheckForUnixWritePermissions(filePath string) bool {

	return unix.Access(filePath, unix.W_OK) == nil
}

// SHOULD BE MOVED TO STHINGS-CLI!
func ReadYamlFile(yamlFileContent []byte) (yamlStructure map[string]interface{}) {

	yamlStructure = make(map[string]interface{})
	data := make(map[interface{}]interface{})

	err := yaml.Unmarshal(yamlFileContent, &data)
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range data {
		yamlStructure[fmt.Sprintf("%v", k)] = fmt.Sprintf("%v", v)
	}

	return
}
