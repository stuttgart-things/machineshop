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

		// FLAGS
		source, _ := cmd.LocalFlags().GetString("source")
		gitPath, _ := cmd.LocalFlags().GetString("path")
		templatePath, _ := cmd.LocalFlags().GetString("template")
		defaultsPath, _ := cmd.LocalFlags().GetString("defaults")
		outputFormat, _ := cmd.LocalFlags().GetString("output")
		destinationPath, _ := cmd.LocalFlags().GetString("destination")

		// PRINT BANNER
		internal.PrintBanner(logFilePath, gitPath, gitRepository, version, date, "/RENDER")

		// fmt.Println(source, templatePath, defaultsPath)
		log.Info("SOURCE: ", source)
		log.Info("TEMPLATE-PATH: ", templatePath)
		log.Info("DEFAULTS: ", defaultsPath)
		log.Info("OUTPUT-FORMAT: ", outputFormat)
		log.Info("DESTINATION-PATH: ", destinationPath+"\n")

		// GET REPO + READ PROFILE FILE
		repo, _ := sthingsCli.CloneGitRepository(gitRepository, gitBranch, gitCommitID, nil)
		templateFile := sthingsCli.ReadFileContentFromGitRepo(repo, templatePath)

		// INIT TEMPLATE VARIABLES
		templateVariables := make(map[string]interface{})

		// READ DEFAULTS (IF DEFINED)
		if defaultsPath != "" {
			defaultsFile = sthingsCli.ReadFileContentFromGitRepo(repo, defaultsPath)
			fmt.Println(defaultsFile)
			templateVariables = internal.ReadYamlFile([]byte(defaultsFile))
		}

		// READ VALUES (IF GIVEN)
		// tbd!

		// RENDER TEMPLATE
		renderedTemplate, err := sthingsBase.RenderTemplateInline(templateFile, "missingkey=zero", "{{", "}}", templateVariables)
		if err != nil {
			fmt.Println(err)
		}

		// HANDLE OUTPUT
		internal.HandleRenderOutput(outputFormat, destinationPath, string(renderedTemplate), true)

	},
}

func init() {
	rootCmd.AddCommand(renderCmd)
	renderCmd.Flags().String("source", "git", "source of profile: git or local")
	renderCmd.Flags().String("template", "tests/template.yaml", "path to to be rendered template")
	renderCmd.Flags().String("defaults", "", "path to defaults template file")
	renderCmd.Flags().String("output", "stdout", "outputFormat stdout|file")
	renderCmd.Flags().String("destination", "", "path to output (if output file)")
}
