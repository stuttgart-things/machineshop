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
	globalValues   = make(map[string]interface{})
)

type Profile struct {
	Name       string   `mapstructure:"name"`
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
		templatesDir, _ := cmd.LocalFlags().GetString("templates")

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

			// STORE TEMPLATES ON FS
			allTemplateFiles, _ := sthingsCli.GetFileListFromGitRepository(templatesDir, repo)
			for _, file := range allTemplateFiles {
				templateFile := sthingsCli.ReadFileContentFromGitRepo(repo, templatesDir+file)
				sthingsBase.WriteDataToFile(workspaceDir+file, templateFile)
				log.Info("CREATED TEMPLATE FILE ON WORKSPACE: ", workspaceDir+file)
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

		// READ DEFAULTS FROM FILES
		for _, defaultsFile := range selectedDefaults {

			if source == "git" {
				defaultsFile = workspaceDir + filepath.Base(defaultsFile)
			}

			log.Info("READING DEFAULTS FROM: ", defaultsFile)

			defaults = sthingsCli.ReadYamlKeyValuesFromFile([]byte(sthingsBase.ReadFileToVariable(defaultsFile)))
			log.Info("DEFAULTS: ", defaults)
			allDefaults = sthingsBase.MergeMaps(allDefaults, defaults)
		}

		// READ DEFAULTS FROM WORKFLOW
		defaultsKeyWorkflow := sthingsCli.GetYamlStringKey("defaults", profilePath, ".yaml")
		defaultsWorkflow := sthingsCli.ReadYamlKeyValuesFromFile([]byte(defaultsKeyWorkflow))
		log.Info("INLINE DEFAULTS FROM WORKFLOW: ", defaultsWorkflow)

		allDefaults = sthingsBase.MergeMaps(allDefaults, defaultsWorkflow)
		log.Info("ALL DEFAULTS: ", allDefaults)

		log.Info("WORKFLOW: ", templateConfig.DefaultProfile.Name)

		// RENDER TEMPLATES
		for _, templateKeys := range templateConfig.TemplateProfile {

			for _, template := range templateKeys {
				log.Info("RENDERING TEMPLATE: ", template.TemplatePath)

				renderedTemplatePath, _ := sthingsCli.RenderTemplateSurvey(template.TemplatePath, allDefaults)
				template.TemplatePath = renderedTemplatePath

				// LOAD TEMPLATE
				templateKey := sthingsCli.GetYamlStringKey("template", template.TemplatePath, ".yaml")
				defaultsKey := sthingsCli.GetYamlStringKey("defaults", template.TemplatePath, ".yaml")

				// LOAD FILE DEFAULTS
				templateDefaults := sthingsCli.ReadYamlKeyValuesFromFile([]byte(defaultsKey))
				log.Info("INLINE DEFAULTS FROM TEMPLATE: ", templateDefaults)

				// MERGE DEFAULT FILES
				allDefaults = sthingsBase.MergeMaps(allDefaults, templateDefaults)
				log.Info("MERGED/ALL DEFAULTS: ", allDefaults)

				// RENDER TEMPLATE
				renderedTemplate, globalValues := sthingsCli.RenderTemplateSurvey(templateKey, allDefaults)
				log.Info("ALL GLOBAL VALUES: ", globalValues)
				fmt.Println(renderedTemplate)

				// MERGE DEFAULT FILES W/ CACHED GLOABLS FROM PREVIOUS RENDERING RUN
				allDefaults = sthingsBase.MergeMaps(allDefaults, globalValues)

				selectedOutputDir := sthingsCli.AskSingleSelectQuestion("SELECT OUTPUT DIR:", []string{"/tmp", workspaceDir})

				sthingsBase.WriteDataToFile(selectedOutputDir+"/"+"hello.yaml", string(renderedTemplate))

				log.Info("RENDERED FILE WAS WRITTEN TO: ", selectedOutputDir+"/"+sthingsCli.AskSingleInputQuestion("Filename", ""))

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
	flowCmd.Flags().String("templates", "machineShop/templates/", "templates dir")
}
