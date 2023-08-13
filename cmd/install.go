/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cmd

import (
	"fmt"

	"github.com/stuttgart-things/machineShop/internal"
	sthingsCli "github.com/stuttgart-things/sthingsCli"

	"github.com/spf13/cobra"
)

var config Profile

type Install struct {
	Url string `mapstructure:"url"`
	Bin string `mapstructure:"bin"`
}

type Script struct {
	Script string `mapstructure:"script"`
}

type Profile struct {
	BinaryProfile []map[string]Install `mapstructure:"binary"`
	ScriptProfile []map[string]Script  `mapstructure:"script"`
}

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

		// READ PROFILE FILE

		templatePath := "/home/sthings/projects/go/src/github/machineShop/tests/install.yaml"
		config = sthingsCli.ReadYamlToObject(templatePath, ".yaml", config).(Profile)

		// INSTALL BINARIES
		for _, binaryProfile := range config.BinaryProfile {

			fmt.Println(binaryProfile["argocd"].Url)
			fmt.Println(binaryProfile["argocd"].Bin)

		}

		// INSTALL SCRIPTS
		for _, scriptProfile := range config.ScriptProfile {

			fmt.Println(scriptProfile["argocd"].Script)

		}

	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().String("source", "git", "source of profile: git or local")
	installCmd.Flags().String("profile", "machineShop.yaml", "path to install profile")
}
