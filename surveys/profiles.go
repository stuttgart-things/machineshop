/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package surveys

import (
	"fmt"

	sthingsCli "github.com/stuttgart-things/sthingsCli"
)

var config Profile

type Install struct {
	Url string `mapstructure:"url"`
	Bin string `mapstructure:"bin"`
}

type Script struct {
	Script string `mapstructure:"script"`
}

type Profile struct {
	BinaryProfile []map[string]Install `mapstructure:"binary"`
	ScriptProfile []map[string]Script  `mapstructure:"script"`
}

func SelectInstallProfiles(yamlFile string) {

	keys := []string{}
	config := sthingsCli.ReadInlineYamlToObject([]byte(yamlFile), config).(Profile)

	// INSTALL BINARIES
	for _, binaryProfile := range config.BinaryProfile {

		for key := range binaryProfile {
			keys = append(keys, key+"-binary")
		}

	}

	// INSTALL SCRIPTS
	for _, scriptProfile := range config.ScriptProfile {

		for key := range scriptProfile {
			keys = append(keys, key+"-script")
		}

	}

	selectedInstall := sthingsCli.AskMultiSelectQuestion("SELECT TO INSTALL:", keys)
	fmt.Println("SELECTED: ", selectedInstall)
}
