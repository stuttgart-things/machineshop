/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package surveys

import (
	"github.com/pterm/pterm"

	sthingsCli "github.com/stuttgart-things/sthingsCli"
)

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

func SelectInstallProfiles(yamlFile string) (selectedBinariesProfiles, selectedScriptProfiles []string, allConfig Profile) {

	// allKeys := []string{}
	allBinaries := []string{}
	allScripts := []string{}

	allConfig = sthingsCli.ReadInlineYamlToObject([]byte(yamlFile), allConfig).(Profile)

	// ITERATE OVER INSTALL BINARIES
	for _, binaryProfile := range allConfig.BinaryProfile {

		for key := range binaryProfile {
			// allKeys = append(allKeys, key)
			allBinaries = append(allBinaries, key)
		}

	}

	// ITERATE OVERINSTALL SCRIPTS
	for _, scriptProfile := range allConfig.ScriptProfile {

		for key := range scriptProfile {
			// allKeys = append(allKeys, key)
			allScripts = append(allScripts, key)
		}

	}

	log.Info("FOUND BINARY PROFILES", len(allBinaries))
	log.Info("FOUND SCRIPT PROFILES", len(allScripts))

	if len(allBinaries) >= 0 {
		selectedBinariesProfiles = sthingsCli.AskMultiSelectQuestion("SELECT BINARIES TO INSTALL:", allBinaries)
	}

	if len(allScripts) >= 0 {
		selectedScriptProfiles = sthingsCli.AskMultiSelectQuestion("SELECT SCRIPTS TO INSTALL:", allScripts)
	}

	return
}

// RELEASE
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

	rendered := pterm.DefaultTable.WithHasHeader().WithBoxed().WithData(repositoryTable).Render()
	if err := rendered; err != nil {
		log.Error(err)
	}

	selectedReleaseProfiles = sthingsCli.AskMultiSelectQuestion("SELECT TO RELEASE:", allKeys)

	return

}
