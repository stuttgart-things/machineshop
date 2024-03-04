/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cmd

import (
	"fmt"

	"github.com/stuttgart-things/machineShop/surveys"
	"github.com/stuttgart-things/machineshop/internal"

	"github.com/spf13/cobra"
	sthingsCli "github.com/stuttgart-things/sthingsCli"
)

// releaseCmd represents the release command
var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "release artifacts",
	Long:  `Release artifact archive files packed from git sources`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("release called")

		tmpPath, _ := cmd.LocalFlags().GetString("tmp")
		gitPath, _ := cmd.LocalFlags().GetString("path")
		profile, _ := cmd.LocalFlags().GetString("profile")

		// PRINT BANNER
		internal.PrintBanner(logFilePath, gitPath, gitRepository, version, date, "/RELEASE")

		// GET REPO + READ PROFILE FILE
		repo, _ := sthingsCli.CloneGitRepository(gitRepository, gitBranch, gitCommitID, nil)
		profileFile := sthingsCli.ReadFileContentFromGitRepo(repo, profile)

		selectedReleaseProfiles, allConfig := surveys.SelectReleaseProfiles(profileFile)
		fmt.Println(selectedReleaseProfiles, allConfig)

		surveys.CloneRepositories(selectedReleaseProfiles, allConfig, tmpPath)

	},
}

func init() {
	rootCmd.AddCommand(releaseCmd)
	releaseCmd.Flags().String("profile", "tests/release.yaml", "path to install profile")
	releaseCmd.Flags().String("tmp", "/tmp", "tmp dir for collection release files")

}
