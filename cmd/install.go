/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cmd

import (
	"fmt"
	"os"

	"github.com/stuttgart-things/machineshop/internal"
	"github.com/stuttgart-things/machineshop/surveys"

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
		url, _ := cmd.LocalFlags().GetString("url")
		bin, _ := cmd.LocalFlags().GetString("bin")
		tmpDownloadDir, _ := cmd.LocalFlags().GetString("tmp")
		binaries, _ := cmd.Flags().GetStringSlice("binaries")

		// PRINT BANNER
		internal.PrintBanner(logFilePath, gitPath, gitRepository, version, date, "/INSTALL")

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

		// INSTALL BINARIES IF DEFINED
		if len(binaries) > 0 {

			var allConfig surveys.Profile
			allConfig = sthingsCli.ReadInlineYamlToObject([]byte(profileFile), allConfig).(surveys.Profile)

			fmt.Println("ALL CONFIG", allConfig)

			surveys.InstallBinaries(binaries, allConfig, tmpDownloadDir, bin)
		} else {
			// GET TO BE INSTALLED BINS + START INSTALL SURVEY
			selectedBinariesProfiles, selectedScriptProfiles, allConfig := surveys.SelectInstallProfiles(profileFile)

			fmt.Println("ALL CONFIG", allConfig.BinaryProfile)

			if len(selectedBinariesProfiles) > 0 {

				binDir := sthingsCli.AskSingleInputQuestion("BIN DIR:", bin)
				surveys.InstallBinaries(selectedBinariesProfiles, allConfig, tmpDownloadDir, binDir)
			}

			if len(selectedScriptProfiles) > 0 {
				surveys.RenderInstallScript(selectedScriptProfiles, allConfig, tmpDownloadDir)
			}

		}

	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().String("tmp", "/tmp/machineShop", "temporary machineShop dir")
	installCmd.Flags().String("bin", "/usr/bin", "target dir for installing binary files")
	installCmd.Flags().String("source", "fetch", "source of profile: local, git or by url fetch. default: fetch.")
	installCmd.Flags().String("url", "https://raw.githubusercontent.com/stuttgart-things/stuttgart-things/refs/heads/main/machineShop/binaries.yaml", "source of url download")
	installCmd.Flags().String("profile", "machineShop/binaries.yaml", "path to install profile")
	installCmd.Flags().StringSlice("binaries", []string{}, "files to be installed; survey will be skipped if defined")
}
