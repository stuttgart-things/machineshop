/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stuttgart-things/machineShop/internal"
	sthingsCli "github.com/stuttgart-things/sthingsCli"
)

// releaseCmd represents the release command
var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "release artifacts",
	Long:  `Release artifact archive files packed from git sources`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("release called")

		gitPath, _ := cmd.LocalFlags().GetString("path")
		profile, _ := cmd.LocalFlags().GetString("profile")

		// PRINT BANNER
		internal.PrintBanner(logFilePath, gitPath, gitRepository, version, date, "/RELEASE")

		// GET REPO + READ PROFILE FILE
		repo, _ := sthingsCli.CloneGitRepository(gitRepository, gitBranch, gitCommitID, nil)
		profileFile := sthingsCli.ReadFileContentFromGitRepo(repo, profile)

		fmt.Println(profileFile)

	},
}

func init() {
	rootCmd.AddCommand(releaseCmd)
	installCmd.Flags().String("profile", "tests/release.yaml", "path to install profile")
}
