/*
Copyright © 2025 Patrick Hermann patrick.hermann@sva.de
*/

package cmd

import (
	"os"

	"github.com/stuttgart-things/machineshop/internal"
	sthingsBase "github.com/stuttgart-things/sthingsBase"
	sthingsCli "github.com/stuttgart-things/sthingsCli"

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

		// LOAD PROFILE BASED ON SOURCE
		switch source {

		case "fetch":
			content, err := internal.DownloadFile(url)
			profileFile = content
			if err != nil {
				log.Error("ERROR BY DOWNLOADING: ", err)
				os.Exit(3)
			}

		case "git":
			// GET REPO + READ PROFILE FILE
			repo, _ := sthingsCli.CloneGitRepository(gitRepository, gitBranch, gitCommitID, nil)
			profileFile = sthingsCli.ReadFileContentFromGitRepo(repo, profile)

		case "local":
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

	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().String("tmp", "/tmp/machineShop", "temporary machineShop dir")
	runCmd.Flags().String("bin", "/usr/bin", "target dir for installing binary files")
	runCmd.Flags().String("source", "fetch", "source of profile: local or by url fetch. default: fetch.")
	runCmd.Flags().String("url", "https://raw.githubusercontent.com/stuttgart-things/stuttgart-things/refs/heads/main/machineShop/binaries.yaml", "source of url download")
	runCmd.Flags().String("profile", "machineShop/binaries.yaml", "path to install profile")
}
