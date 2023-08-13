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

func SelectInstallProfiles() {

	templatePath := "/home/sthings/projects/go/src/github/machineShop/tests/install.yaml"
	config = sthingsCli.ReadYamlToObject(templatePath, ".yaml", config).(Profile)

	// INSTALL BINARIES
	for _, binaryProfile := range config.BinaryProfile {

		fmt.Println(binaryProfile["argocd"].Url)
		fmt.Println(binaryProfile["argocd"].Bin)

	}

	// INSTALL SCRIPTS
	for _, scriptProfile := range config.ScriptProfile {

		fmt.Println(scriptProfile["argocd"].Script)

	}

	// local or git
}
