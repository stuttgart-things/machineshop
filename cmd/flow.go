/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	sthingsBase "github.com/stuttgart-things/sthingsBase"
	sthingsCli "github.com/stuttgart-things/sthingsCli"
)

var (
	templatePath   = "tests/template.yaml"
	renderOption   = "missingkey=zero"
	templateConfig Default
	defaults       = make(map[string]interface{})
	templateKeys   = make(map[string]int)
	allDefaults    map[string]interface{}
	globalValues   map[string]interface{}
)

type Profile struct {
	Defaults   []string `mapstructure:"defaults"`
	Repository string   `mapstructure:"repository"`
}

type Default struct {
	DefaultProfile  Profile               `mapstructure:"profile"`
	TemplateProfile []map[string]Template `mapstructure:"templates"`
}

type Template struct {
	TemplatePath     string `mapstructure:"template"`
	DefaultTargetDir string `mapstructure:"defaultTargetDir"`
}

// flowCmd represents the flow command
var flowCmd = &cobra.Command{
	Use:   "flow",
	Short: "run workflows",
	Long:  `run workflows based on profiles`,

	Run: func(cmd *cobra.Command, args []string) {

		// READ FLAGS
		profilePath, _ := cmd.LocalFlags().GetString("profile")
		source, _ := cmd.LocalFlags().GetString("source")
		workspaceDir, _ := cmd.LocalFlags().GetString("workspace")
		defaultsDir, _ := cmd.LocalFlags().GetString("defaults")

		// READ CONFIG FROM GIT REPO
		if source == "git" {

			// GET REPO + READ PROFILE FILE
			sthingsBase.CreateNestedDirectoryStructure(workspaceDir, 0777)
			log.Info("CREATED WORKSPACE DIR: ", workspaceDir)

			// CLONE REPOSITORY
			repo, _ := sthingsCli.CloneGitRepository(gitRepository, gitBranch, gitCommitID, nil)

			// STORE PROFILE ON FS
			profileFile := sthingsCli.ReadFileContentFromGitRepo(repo, profilePath)
			sthingsBase.WriteDataToFile(workspaceDir+filepath.Base(profilePath), string(profileFile))
			profilePath = workspaceDir + filepath.Base(profilePath)
			log.Info("CREATED PROFILE FILE ON WORKSPACE: ", profilePath)

			// STORE DEFAULTS ON FS
			allDefaultFiles, _ := sthingsCli.GetFileListFromGitRepository(defaultsDir, repo)
			for _, file := range allDefaultFiles {
				defaultFile := sthingsCli.ReadFileContentFromGitRepo(repo, defaultsDir+file)
				sthingsBase.WriteDataToFile(workspaceDir+file, defaultFile)
				log.Info("CREATED DEFAULT FILE ON WORKSPACE: ", workspaceDir+file)
			}
		}

		// READ PROFILE FILE
		templateConfig = sthingsCli.ReadYamlToObject(profilePath, ".yaml", templateConfig).(Default)
		log.Info("LOCAL PROFILE READ IN: ", profilePath)
		for i, config := range templateConfig.TemplateProfile {
			for template := range config {
				templateKeys[template] = i
			}
		}

		// SELECT DEFAULTS
		selectedDefaults := sthingsCli.AskMultiSelectQuestion("SELECT DEFAULT FILE(S):", templateConfig.DefaultProfile.Defaults)

		// READ DEFAULTS
		for _, defaultsFile := range selectedDefaults {

			if source == "git" {
				defaultsFile = workspaceDir + filepath.Base(defaultsFile)
			}

			log.Info("READING DEFAULTS FROM: ", defaultsFile)

			defaults = sthingsCli.ReadYamlKeyValuesFromFile([]byte(sthingsBase.ReadFileToVariable(defaultsFile)))
			log.Info("DEFAULTS: ", defaults)
			allDefaults = sthingsBase.MergeMaps(allDefaults, defaults)
		}
		log.Info("ALL DEFAULTS: ", allDefaults)

		// RENDER TEMPLATES
		for _, templateKeys := range templateConfig.TemplateProfile {

			for _, template := range templateKeys {
				log.Info("RENDERING TEMPLATE: ", template.TemplatePath)

				// LOAD TEMPLATE
				templateKey := sthingsCli.GetYamlStringKey("template", template.TemplatePath, ".yaml")
				defaultsKey := sthingsCli.GetYamlStringKey("defaults", template.TemplatePath, ".yaml")

				// LOAD FILE DEFAULTS
				templateDefaults := sthingsCli.ReadYamlKeyValuesFromFile([]byte(defaultsKey))
				log.Info("INLINE DEFAULTS FROM TEMPLATE: ", templateDefaults)

				allDefaults = sthingsBase.MergeMaps(allDefaults, templateDefaults)
				log.Info("MERGED/ALL DEFAULTS: ", allDefaults)

				renderedTemplate, globalValues := sthingsCli.RenderTemplateSurvey(templateKey, allDefaults)
				fmt.Println(renderedTemplate)

				allDefaults = sthingsBase.MergeMaps(allDefaults, globalValues)
				fmt.Println(allDefaults)

				selectedOuputDir := sthingsCli.AskMultiSelectQuestion("SELECT OUTPUT DIR:", []string{"/tmp", workspaceDir})
				sthingsBase.WriteDataToFile(selectedOuputDir+"/"+templateKeys+".yaml", string(renderedTemplate))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(flowCmd)
	flowCmd.Flags().String("profile", "machineShop/workflows/build-vm-provision.yaml", "path to to be used profile")
	flowCmd.Flags().String("source", "git", "source of profile: git or local")
	flowCmd.Flags().String("workspace", "/tmp/machineShopFlow/", "dir for tmp files")
	flowCmd.Flags().String("defaults", "machineShop/defaults/", "defaults dir")
}