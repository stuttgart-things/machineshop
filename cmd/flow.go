/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	sthingsBase "github.com/stuttgart-things/sthingsBase"
	sthingsCli "github.com/stuttgart-things/sthingsCli"
)

var (
	templatePath   = "tests/template.yaml"
	renderOption   = "missingkey=zero"
	templateConfig Default
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

		// READ PROFILE
		templateConfig = sthingsCli.ReadYamlToObject(profilePath, ".yaml", templateConfig).(Default)

		// READ TEMPLATE KEYS
		for i, config := range templateConfig.TemplateProfile {
			for template := range config {
				templateKeys[template] = i
			}
		}

		log.Info("LOCAL PROFILE READ : ", profilePath)

		selectedDefaults := sthingsCli.AskMultiSelectQuestion("SELECT DEFAULT FILE(S):", templateConfig.DefaultProfile.Defaults)

		// READ DEFAULTS
		for _, defaultsFile := range selectedDefaults {
			log.Info("READ DEFAULTS FROM: ", defaultsFile)
			defaults := sthingsCli.ReadYamlKeyValuesFromFile([]byte(sthingsBase.ReadFileToVariable(defaultsFile)))
			log.Info("DEFAULTS: ", defaults)
			allDefaults = sthingsBase.MergeMaps(allDefaults, defaults)
		}
		log.Info("ALL DEFAULTS: ", allDefaults)

		for _, templateKeys := range templateConfig.TemplateProfile {

			for _, i := range templateKeys {
				log.Info("RENDERING TEMPLATE: ", i.TemplatePath)

				// LOAD TEMPLATE
				templateKey := sthingsCli.GetYamlStringKey("template", i.TemplatePath, ".yaml")
				defaultsKey := sthingsCli.GetYamlStringKey("defaults", i.TemplatePath, ".yaml")

				// LOAD FILE DEFAULTS
				templateDefaults := sthingsCli.ReadYamlKeyValuesFromFile([]byte(defaultsKey))
				log.Info("INLINE DEFAULTS FROM TEMPLATE: ", templateDefaults)

				allDefaults = sthingsBase.MergeMaps(allDefaults, templateDefaults)
				log.Info("MERGED/ALL DEFAULTS: ", allDefaults)

				renderedTemplate, globalValues := sthingsCli.RenderTemplateSurvey(templateKey, allDefaults)
				fmt.Println(renderedTemplate)

				allDefaults = sthingsBase.MergeMaps(allDefaults, globalValues)
				fmt.Println(allDefaults)

				sthingsBase.WriteDataToFile("/tmp/hello.yaml", string(renderedTemplate))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(flowCmd)
	flowCmd.Flags().String("profile", "tests/workflow.yaml", "path to to be used profile")

}
