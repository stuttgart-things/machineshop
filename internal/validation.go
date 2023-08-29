/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package internal

import (
	"fmt"
	"log"

	"gopkg.in/yaml.v3"
)

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
