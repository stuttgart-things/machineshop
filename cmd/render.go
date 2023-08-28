/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	sthingsBase "github.com/stuttgart-things/sthingsBase"

	"github.com/stuttgart-things/machineShop/internal"
	sthingsCli "github.com/stuttgart-things/sthingsCli"

	"github.com/spf13/cobra"
)

var defaultsFile string

// renderCmd represents the render command
var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "render templates",
	Long:  `Render templates to files`,

	Run: func(cmd *cobra.Command, args []string) {

		source, _ := cmd.LocalFlags().GetString("source")
		gitPath, _ := cmd.LocalFlags().GetString("path")
		templatePath, _ := cmd.LocalFlags().GetString("template")
		defaultsPath, _ := cmd.LocalFlags().GetString("defaults")

		fmt.Println(source, templatePath, defaultsPath)

		// PRINT BANNER
		internal.PrintBanner(logFilePath, gitPath, gitRepository, version, date, "/RENDER")

		// GET REPO + READ PROFILE FILE
		repo, _ := sthingsCli.CloneGitRepository(gitRepository, gitBranch, gitCommitID, nil)
		templateFile := sthingsCli.ReadFileContentFromGitRepo(repo, templatePath)

		// INIT TEMPLATE VARIABLES
		templateVariables := make(map[string]interface{})

		// READ DEFAULTS IF DEFINED
		if defaultsPath != "" {
			defaultsFile = sthingsCli.ReadFileContentFromGitRepo(repo, defaultsPath)
			fmt.Println(defaultsFile)
			templateVariables = internal.ReadYamlFile([]byte(defaultsFile))
		}

		// templateVariables["chartName"] = "gude-chart"
		// templateVariables["namespace"] = "gude"
		// templateVariables["randomName"] = true

		rendered, err := sthingsBase.RenderTemplateInline(templateFile, "missingkey=zero", "{{", "}}", templateVariables)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(string(rendered))

	},
}

func init() {
	rootCmd.AddCommand(renderCmd)
	renderCmd.Flags().String("source", "git", "source of profile: git or local")
	renderCmd.Flags().String("template", "tests/template.yaml", "path to to be rendered template")
	renderCmd.Flags().String("defaults", "", "path to defaults template file")
}
