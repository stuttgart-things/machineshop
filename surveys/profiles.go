/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package surveys

import (
	"github.com/pterm/pterm"

	sthingsCli "github.com/stuttgart-things/sthingsCli"
)

// INSTALL
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

// RELEASE
var allReleaseProfiles ReleaseProfile

type ReleaseRepository struct {
	Url     string `mapstructure:"url"`
	Branch  string `mapstructure:"branch"`
	Version string `mapstructure:"version"`
}

type ReleaseProfile struct {
	RepositoryProfile []map[string]ReleaseRepository `mapstructure:"repositories"`
}

func SelectReleaseProfiles(yamlFile string) (selectedReleaseProfiles []string, allReleaseProfiles ReleaseProfile) {

	allKeys := []string{}
	repositoryTable := pterm.TableData{}
	header := []string{"NAME", "URL", "BRANCH", "VERSION"}
	repositoryTable = append(repositoryTable, header)

	// allKeys := []string{}
	allReleaseProfiles = sthingsCli.ReadInlineYamlToObject([]byte(yamlFile), allReleaseProfiles).(ReleaseProfile)

	// ITERATE OVER INSTALL BINARIES
	for _, repositoryProfile := range allReleaseProfiles.RepositoryProfile {

		for key, value := range repositoryProfile {

			values := []string{key, value.Url, value.Branch, value.Version}
			allKeys = append(allKeys, key)
			repositoryTable = append(repositoryTable, values)
		}

	}

	pterm.DefaultTable.WithHasHeader().WithBoxed().WithData(repositoryTable).Render()

	selectedReleaseProfiles = sthingsCli.AskMultiSelectQuestion("SELECT TO RELEASE:", allKeys)

	return

}
