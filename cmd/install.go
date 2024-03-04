/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cmd

import (
	"os"

	"github.com/stuttgart-things/machineShop/surveys"
	"github.com/stuttgart-things/machineshop/internal"

	sthingsBase "github.com/stuttgart-things/sthingsBase"
	sthingsCli "github.com/stuttgart-things/sthingsCli"

	"github.com/spf13/cobra"
)

var (
	profileFile string
)
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install binaries",
	Long:  `Install binaries`,
	Run: func(cmd *cobra.Command, args []string) {

		gitPath, _ := cmd.LocalFlags().GetString("path")
		profile, _ := cmd.LocalFlags().GetString("profile")
		source, _ := cmd.LocalFlags().GetString("source")
		bin, _ := cmd.LocalFlags().GetString("bin")
		tmpDownloadDir, _ := cmd.LocalFlags().GetString("tmp")

		// PRINT BANNER
		internal.PrintBanner(logFilePath, gitPath, gitRepository, version, date, "/INSTALL")

		// LOAD PROFILE BASED ON SOURCE
		if source == "git" {
			// GET REPO + READ PROFILE FILE
			repo, _ := sthingsCli.CloneGitRepository(gitRepository, gitBranch, gitCommitID, nil)
			profileFile = sthingsCli.ReadFileContentFromGitRepo(repo, profile)

		} else if source == "local" {

			profileExists, _ := sthingsBase.VerifyIfFileOrDirExists(profile, "file")
			log.Info("LOCAL PROFILE FOUND : ", profile)

			if profileExists {
				profileFile = sthingsBase.ReadFileToVariable(profile)
			} else {
				log.Error("LOCAL PROFILE NOT FOUND : ", profile)
				os.Exit(3)
			}

		} else {
			log.Error("SOURCE: GIT OR LOCAL ONLY", source)
			os.Exit(3)
		}

		// GET TO BE INSTALLED BINS + START INSTALL SURVEY
		selectedBinariesProfiles, selectedScriptProfiles, allConfig := surveys.SelectInstallProfiles(profileFile)

		if len(selectedBinariesProfiles) > 0 {
			surveys.InstallBinaries(selectedBinariesProfiles, allConfig, tmpDownloadDir, bin)
		}

		if len(selectedScriptProfiles) > 0 {
			surveys.RenderInstallScript(selectedScriptProfiles, allConfig, tmpDownloadDir)
		}

	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().String("tmp", "/tmp/machineShop", "temporary machineShop dir")
	installCmd.Flags().String("bin", "/usr/bin", "target dir for installing binary files")
	installCmd.Flags().String("source", "git", "source of profile: git or local")
	installCmd.Flags().String("profile", "tests/install.yaml", "path to install profile")
}
