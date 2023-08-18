/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cmd

import (
	"fmt"

	"github.com/stuttgart-things/machineShop/surveys"

	"github.com/stuttgart-things/machineShop/internal"
	sthingsCli "github.com/stuttgart-things/sthingsCli"

	"github.com/spf13/cobra"
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

		fmt.Println(source)

		// PRINT BANNER
		internal.PrintBanner(logFilePath, gitPath, gitRepository, version, date, "/INSTALL")

		// GET REPO + READ PROFILE FILE
		repo, _ := sthingsCli.CloneGitRepository(gitRepository, gitBranch, gitCommitID, nil)
		profileFile := sthingsCli.ReadFileContentFromGitRepo(repo, profile)

		// GET CONFIG
		selectedInstallProfiles, allConfig := surveys.SelectInstallProfiles(profileFile)
		fmt.Println(selectedInstallProfiles, allConfig)

		surveys.InstallBin(selectedInstallProfiles, allConfig, bin)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().String("bin", "/usr/bin/", "source of profile: git or local")
	installCmd.Flags().String("source", "git", "source of profile: git or local")
	installCmd.Flags().String("profile", "tests/install.yaml", "path to install profile")
}
