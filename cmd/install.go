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

// var config Profile

// type Install struct {
// 	Url string `mapstructure:"url"`
// 	Bin string `mapstructure:"bin"`
// }

// type Script struct {
// 	Script string `mapstructure:"script"`
// }

// type Profile struct {
// 	BinaryProfile []map[string]Install `mapstructure:"binary"`
// 	ScriptProfile []map[string]Script  `mapstructure:"script"`
// }

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install binaries",
	Long:  `Install binaries`,
	Run: func(cmd *cobra.Command, args []string) {

		gitPath, _ := cmd.LocalFlags().GetString("path")
		profile, _ := cmd.LocalFlags().GetString("profile")
		source, _ := cmd.LocalFlags().GetString("source")

		internal.PrintBanner(logFilePath, gitPath, gitRepository, version, date, "/INSTALL")

		fmt.Println(profile)
		fmt.Println(source)
		repo, _ := sthingsCli.CloneGitRepository(gitRepository, gitBranch, gitCommitID, nil)

		fileList, _ := sthingsCli.GetFileListFromGitRepository(gitPath, repo)
		selectedProfile := sthingsCli.AskSingleSelectQuestion("SELECT PROFILE:", fileList)
		fmt.Println(selectedProfile)

		fileList, directoryList := sthingsCli.GetFileListFromGitRepository("", repo)
		fmt.Println(fileList, directoryList)

		// READ PROFILE FILE
		file := sthingsCli.ReadFileContentFromGitRepo(repo, "tests/install.yaml")
		fmt.Println("File content: %+v\n", file)

		// READ FROM SCRIPT HERE
		surveys.SelectInstallProfiles()

	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().String("source", "git", "source of profile: git or local")
	installCmd.Flags().String("profile", "machineShop.yaml", "path to install profile")
}
