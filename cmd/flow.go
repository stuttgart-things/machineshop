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
		// fmt.Println(templateConfig.DefaultProfile.Defaults)
		// fmt.Println(templateConfig.DefaultProfile.Repository)

		// READ TEMPLATE KEYS
		for i, config := range templateConfig.TemplateProfile {
			for template := range config {
				templateKeys[template] = i
			}
		}

		log.Info("LOCAL PROFILE READ : ", profilePath)

		// SELECT TEMPLATES
		// selectedTemplates := sthingsCli.AskMultiSelectQuestion("SELECT TO BE RENDERED TEMPLATE(S):", maps.Keys(templateKeys))
		// log.Info("SELECTED TO BE USED ACTIONS: ", selectedTemplates)
		selectedDefaults := sthingsCli.AskMultiSelectQuestion("SELECT TO BE USED DEFAULT(S):", templateConfig.DefaultProfile.Defaults)

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
				fmt.Println(i.TemplatePath)
			}
		}

		// LOAD TEMPLATE
		templateKey := sthingsCli.GetYamlStringKey("template", "tests/maverick", ".yaml")
		defaultsKey := sthingsCli.GetYamlStringKey("defaults", "tests/maverick", ".yaml")

		// fmt.Println(templateKey)

		// LOAD FILE DEFAULTS
		fmt.Println(defaultsKey)
		templateDefaults := sthingsCli.ReadYamlKeyValuesFromFile([]byte(defaultsKey))
		log.Info("DEFAULTS FROM FILE: ", templateDefaults)

		allDefaults = sthingsBase.MergeMaps(allDefaults, templateDefaults)
		log.Info("ALL DEFAULTS: ", allDefaults)

		renderedTemplate, globalValues := sthingsCli.RenderTemplateSurvey(templateKey, allDefaults)
		fmt.Println(renderedTemplate)
		fmt.Println(globalValues)

		// GetYamlStringKey
		// renderedTemplate, err := sthingsBase.RenderTemplateInline(templateFile, renderOption, brackets["curly"].begin, brackets["curly"].end, defaultVariables)
		// if err != nil {
		// 	fmt.Println(err)
		// }

		// fmt.Println(string(renderedTemplate))

		sthingsBase.WriteDataToFile("/tmp/hello.yaml", string(renderedTemplate))

	},
}

func init() {
	rootCmd.AddCommand(flowCmd)
	flowCmd.Flags().String("profile", "tests/workflow.yaml", "path to to be used profile")

}
