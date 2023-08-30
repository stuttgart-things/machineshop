/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	sthingsBase "github.com/stuttgart-things/sthingsBase"

	billy "github.com/go-git/go-billy/v5"
	"github.com/stuttgart-things/machineShop/internal"
	sthingsCli "github.com/stuttgart-things/sthingsCli"

	"github.com/spf13/cobra"
)

var (
	templateFile     string
	defaultsFile     string
	repo             billy.Filesystem
	defaultVariables = make(map[string]interface{})
	flagVariables    = make(map[string]interface{})
)

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
		templateValues, _ := cmd.Flags().GetStringSlice("values")
		forceRenderOption, _ := cmd.LocalFlags().GetBool("force")

		// PRINT BANNER
		internal.PrintBanner(logFilePath, gitPath, gitRepository, version, date, "/RENDER")

		// fmt.Println(source, templatePath, defaultsPath)
		log.Info("TEMPLATE-PATH: ", templatePath)
		log.Info("DEFAULTS: ", defaultsPath)
		log.Info("OUTPUT-FORMAT: ", outputFormat)
		log.Info("DESTINATION-PATH: ", destinationPath+"\n")

		// GET REPO + READ TEMPLATE + DEFAULTS

		// HANDLE SOURCE:GIT FOR TEMPLATE + DEFAULTS
		if source == "git" {
			repo, _ = sthingsCli.CloneGitRepository(gitRepository, gitBranch, gitCommitID, nil)
			templateFile = sthingsCli.ReadFileContentFromGitRepo(repo, templatePath)

			if defaultsPath != "" {
				defaultsFile = sthingsCli.ReadFileContentFromGitRepo(repo, defaultsPath)
				log.Info("LOADED DEFAULTS FILE FROM: ", defaultsPath)
				fmt.Println(defaultsFile)
				defaultVariables = internal.ReadYamlFile([]byte(defaultsFile))
			} else {
				log.Info("NO DEFAULTS FILE FROM GIT DEFINED")
			}

			// HANDLE SOURCE:LOCAL FOR TEMPLATE + DEFAULTS
		} else if source == "local" {
			templateExists, _ := sthingsBase.VerifyIfFileOrDirExists(templatePath, "file")
			log.Info("LOCAL TEMPLATE FOUND : ", templatePath)

			if templateExists {
				templateFile = sthingsBase.ReadFileToVariable(templatePath)
			} else {
				log.Error("LOCAL TEMPLATE NOT FOUND : ", templatePath)
				os.Exit(3)
			}

			if defaultsPath != "" {

				defaultsFileExists, _ := sthingsBase.VerifyIfFileOrDirExists(defaultsPath, "file")

				if defaultsFileExists {
					defaultsFile = sthingsBase.ReadFileToVariable(defaultsPath)
					log.Info("LOADED DEFAULTS FILE FROM: ", defaultsPath)
					fmt.Println(defaultsFile)
					defaultVariables = internal.ReadYamlFile([]byte(defaultsFile))
				} else {
					log.Error("LOCAL DEFAULTS FILE NOT FOUND : ", templatePath)
					os.Exit(3)
				}
			} else {
				log.Info("NO DEFAULTS FILE DEFINED")
			}

		} else {
			log.Error("SOURCE COULD BE ONLY GIT OR LOCAL", source)
			os.Exit(3)
		}

		// READ DEFAULTS (IF DEFINED)

		// READ VALUES (IF DEFINED)
		if len(templateValues) > 0 {
			flagVariables = internal.VerifyReadKeyValues(templateValues, log)
			fmt.Println("VALUES", flagVariables)
			// log.SayWithField("reading values..", "values", flagVariables)
		} else {
			log.Warn("NO VALUES DEFINED")
		}

		variables := internal.MergeMaps(defaultVariables, flagVariables)

		// RENDER TEMPLATE
		renderOption := "missingkey=error"
		if forceRenderOption {
			renderOption = "missingkey=zero"
		}

		renderedTemplate, err := sthingsBase.RenderTemplateInline(templateFile, renderOption, "{{", "}}", variables)
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
	renderCmd.Flags().Bool("force", false, "force rendering by missing keys")
	renderCmd.Flags().StringSlice("values", []string{}, "templating values")
}
