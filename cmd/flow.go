/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cmd

import (
	"fmt"

	sthingsBase "github.com/stuttgart-things/sthingsBase"
	sthingsCli "github.com/stuttgart-things/sthingsCli"
	"golang.org/x/exp/maps"

	"github.com/spf13/cobra"
)

var (
	templatePath   = "tests/template.yaml"
	renderOption   = "missingkey=zero"
	templateConfig Default
	templateKeys   = make(map[string]int)
)

type Profile struct {
	Defaults    []string `mapstructure:"defaults"`
	Repository  string   `mapstructure:"repository"`
	Environment []string `mapstructure:"environment"`
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
		fmt.Println(templateConfig.DefaultProfile.Defaults)
		fmt.Println(templateConfig.DefaultProfile.Repository)

		for i, config := range templateConfig.TemplateProfile {

			for template := range config {
				templateKeys[template] = i
			}
		}

		log.Info("LOCAL PROFILE READ : ", profilePath)

		selectedTemplates := sthingsCli.AskMultiSelectQuestion("SELECT TO BE RENDERED TEMPLATE(S):", maps.Keys(templateKeys))
		fmt.Println(selectedTemplates)

		selectedDefaults := sthingsCli.AskMultiSelectQuestion("SELECT TO BE USED DEFAULT(S):", templateConfig.DefaultProfile.Environment)
		fmt.Println(selectedDefaults)

		readInProfile := sthingsBase.ReadFileToVariable(profilePath)

		defaultVariables := sthingsCli.ReadYamlKeyValuesFromFile([]byte(readInProfile))
		fmt.Println(defaultVariables)

		// renderedTemplate, err := sthingsBase.RenderTemplateInline(templateFile, renderOption, brackets["curly"].begin, brackets["curly"].end, defaultVariables)
		// if err != nil {
		// 	fmt.Println(err)
		// }

		// fmt.Println(string(renderedTemplate))

		// sthingsBase.WriteDataToFile("./hello.yaml", string(renderedTemplate))

	},
}

func init() {
	rootCmd.AddCommand(flowCmd)
	flowCmd.Flags().String("profile", "tests/workflow.yaml", "path to to be used profile")

}
