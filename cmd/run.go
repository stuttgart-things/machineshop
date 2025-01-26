/*
Copyright © 2025 Patrick Hermann patrick.hermann@sva.de
*/

package cmd

import (
	"os"

	"github.com/stuttgart-things/machineshop/internal"
	"github.com/stuttgart-things/machineshop/surveys"
	sthingsBase "github.com/stuttgart-things/sthingsBase"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run commands",
	Long:  `Run commands based on profile templates`,
	Run: func(cmd *cobra.Command, args []string) {

		// GETTING FLAGS
		source, _ := cmd.LocalFlags().GetString("source")
		url, _ := cmd.LocalFlags().GetString("url")
		profile, _ := cmd.LocalFlags().GetString("profile")
		scripts, _ := cmd.Flags().GetStringSlice("scripts")
		variables, _ := cmd.Flags().GetStringSlice("variables")

		// LOAD PROFILE BASED ON SOURCE
		switch source {

		case "fetch":
			content, err := internal.DownloadFile(url)
			profileFile = content
			if err != nil {
				log.Error("ERROR BY DOWNLOADING: ", err)
				os.Exit(3)
			}

		case loadFromLocal:
			// GET LOCAL FILE
			profileExists, _ := sthingsBase.VerifyIfFileOrDirExists(profile, "file")
			log.Info("LOCAL PROFILE FOUND : ", profile)

			if profileExists {
				profileFile = sthingsBase.ReadFileToVariable(profile)
			} else {
				log.Error("LOCAL PROFILE NOT FOUND : ", profile)
				os.Exit(3)
			}

		default:
			log.Error("SOURCE: GIT OR LOCAL ONLY", source)
			os.Exit(3)
		}

		// READ VARIABLES (IF DEFINED)
		if len(variables) > 0 {
			flagVariables = internal.VerifyReadKeyValues(variables, log, enableVault)
			log.Info("VARIABLES", flagVariables)
		} else {
			log.Warn("NO VARIABLES DEFINED")
		}

		// PARSE PROFILE
		runConfig := internal.LoadRunConfig(profileFile)
		surveys.RenderInstallScriptNew(runConfig, scripts, flagVariables)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().String("tmp", "/tmp/machineShop", "temporary machineShop dir")
	runCmd.Flags().String("source", "fetch", "source of profile: local or by url fetch. default: fetch.")
	runCmd.Flags().String("url", "https://raw.githubusercontent.com/stuttgart-things/machineshop/refs/heads/main/profiles/run.yaml", "profile download url")
	runCmd.Flags().String("profile", "profiles/run.yaml", "path to run profile")
	runCmd.Flags().StringSlice("scripts", []string{}, "scripts to be run; survey will be skipped if defined")
	runCmd.Flags().StringSlice("variables", []string{}, "variables to be used in the scripts")
}
