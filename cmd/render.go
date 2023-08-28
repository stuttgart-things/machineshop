/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

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

		if defaultsPath != "" {
			defaultsFile = sthingsCli.ReadFileContentFromGitRepo(repo, templatePath)
		}

		fmt.Println(templateFile)

	},
}

func init() {
	rootCmd.AddCommand(renderCmd)
	renderCmd.Flags().String("source", "git", "source of profile: git or local")
	renderCmd.Flags().String("template", "tests/template.yaml", "path to to be rendered template")
	renderCmd.Flags().String("defaults", "", "path to defaults template file")
}
