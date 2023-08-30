/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package internal

import (
	"fmt"
	"log"
	"os"
	"strings"

	sthingsBase "github.com/stuttgart-things/sthingsBase"

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

func VerifyReadKeyValues(templateValues []string, log *sthingsBase.Logger) map[string]interface{} {

	templateData := make(map[string]interface{})

	if len(templateValues) > 0 {

		if sthingsBase.VerifyIfStringIsBase64(templateValues[0]) {
			base64decodedValues := sthingsBase.DecodeBase64String(templateValues[0])
			templateValues = strings.Split(base64decodedValues, ",")
		}

		for _, v := range templateValues {
			values := strings.Split(v, "=")

			// IF VALUE IS BASE64 ENDCODED THIS IS NEEDED TO PATCH THE STRING BACK IF = ARE INCLUDED
			values = []string{values[0], strings.Join(values[1:], "=")}

			// CHECK FOR EMPTY KEY
			if strings.TrimSpace(values[0]) == "" {
				fmt.Println("no key for value", strings.TrimSpace(values[1]), "defined. exiting")
				log.Error("no key defined. exiting")
				os.Exit(3)
			}

			templateData[strings.TrimSpace(values[0])] = strings.TrimSpace(values[1])

		}
	} else {
		log.Warn("no values defined..")
	}

	return templateData
}

func MergeMaps[K comparable, V any](m1 map[K]V, m2 map[K]V) map[K]V {
	merged := make(map[K]V)
	for key, value := range m1 {
		merged[key] = value
	}
	for key, value := range m2 {
		merged[key] = value
	}
	return merged
}
