/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cmd

import (
	"fmt"

	sthingsBase "github.com/stuttgart-things/sthingsBase"
	sthingsCli "github.com/stuttgart-things/sthingsCli"

	"github.com/spf13/cobra"
)

var (
	templatePath = "tests/template.yaml"
	renderOption = "missingkey=zero"
)

type Profile struct {
	Defaults    string `mapstructure:"defaults"`
	Repository  string `mapstructure:"repository"`
	Environment string `mapstructure:"environment"`
}

type Default struct {
	DefaultProfile Profile `mapstructure:"profile"`
}

type Binary struct {
	Url string `mapstructure:"url"`
	Bin string `mapstructure:"bin"`
}

type Config struct {
	BinaryProfile []map[string]Binary `mapstructure:"binaries"`
}

// flowCmd represents the flow command
var flowCmd = &cobra.Command{
	Use:   "flow",
	Short: "run workflows",
	Long:  `run workflows based on profiles`,

	Run: func(cmd *cobra.Command, args []string) {

		// READ FLAGS
		profilePath, _ := cmd.LocalFlags().GetString("profile")
		var config Config
		var config1 Default
		m := make(map[string]int)

		// READ PROFILE
		log.Info("LOCAL PROFILE FOUND : ", profilePath)

		config1 = sthingsCli.ReadYamlToObject(profilePath, ".yaml", config1).(Default)
		fmt.Println(config1.DefaultProfile.Defaults)

		config = sthingsCli.ReadYamlToObject(profilePath, ".yaml", config).(Config)

		for i, config := range config.BinaryProfile {

			for binary := range config {
				m[binary] = i
			}
		}

		readInProfile := sthingsBase.ReadFileToVariable(profilePath)

		templateFile := sthingsBase.ReadFileToVariable(templatePath)

		defaultVariables := sthingsCli.ReadYamlKeyValuesFromFile([]byte(readInProfile))
		fmt.Println(defaultVariables)

		hello := defaultVariables["environment"]
		fmt.Println(hello)

		// sthingsCli.AskMultiSelectQuestion("SELECT BINARIES TO INSTALL:", bla)

		// bla = defaultVariables["environment"].(string)

		// s := strings.Split(hello.(string), "")
		// fmt.Println(s)

		// sthingsCli.AskMultiSelectQuestion("SELECT BINARIES TO INSTALL:", strings.Split(hello, ""))

		renderedTemplate, err := sthingsBase.RenderTemplateInline(templateFile, renderOption, brackets["curly"].begin, brackets["curly"].end, defaultVariables)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(string(renderedTemplate))

		sthingsBase.WriteDataToFile("./hello.yaml", string(renderedTemplate))

	},
}

func init() {
	rootCmd.AddCommand(flowCmd)
	flowCmd.Flags().String("profile", "tests/workflow.yaml", "path to to be used profile")

}
