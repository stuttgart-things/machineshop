/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package surveys

import (
	sthingsCli "github.com/stuttgart-things/sthingsCli"
)

var allConfig Profile

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

func SelectInstallProfiles(yamlFile string) (selectedInstallProfiles []string, allConfig Profile) {

	allKeys := []string{}
	allConfig = sthingsCli.ReadInlineYamlToObject([]byte(yamlFile), allConfig).(Profile)

	// ITERATE OVER INSTALL BINARIES
	for _, binaryProfile := range allConfig.BinaryProfile {

		for key := range binaryProfile {
			allKeys = append(allKeys, key)
		}

	}

	// ITERATE OVERINSTALL SCRIPTS
	for _, scriptProfile := range allConfig.ScriptProfile {

		for key := range scriptProfile {
			allKeys = append(allKeys, key)
		}

	}

	selectedInstallProfiles = sthingsCli.AskMultiSelectQuestion("SELECT TO INSTALL:", allKeys)
	return
}
