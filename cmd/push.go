/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cmd

import (
	"fmt"
	"os"

	"github.com/stuttgart-things/machineShop/internal"
	sthingsCli "github.com/stuttgart-things/sthingsCli"

	sthingsBase "github.com/stuttgart-things/sthingsBase"

	"github.com/spf13/cobra"
)

var (
	commitMessage = "pushed w/ machineShop CLI"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "push artifacts",
	Long:  `push artifacts target external systems`,
	Run: func(cmd *cobra.Command, args []string) {

		// FLAGS
		target, _ := cmd.LocalFlags().GetString("target")
		sourceFile, _ := cmd.LocalFlags().GetString("source")
		destinationPath, _ := cmd.LocalFlags().GetString("destination")

		// VERIFY IF SOURCE FILE IS EXISTING
		if sourceFile != "" {
			sourceExists, _ := sthingsBase.VerifyIfFileOrDirExists(sourceFile, "file")
			if sourceExists {
				log.Info("SOURCE FOUND : ", sourceFile)
			} else {
				log.Error("SOURCE NOT FOUND : ", sourceFile)
				os.Exit(3)
			}
		} else {
			log.Error("SOURCE UNDEFINED")
			os.Exit(3)
		}

		switch target {

		case "s3":
			fmt.Println("s3")
			// VERIFY S3 ENV VARS
			// MINIO_URL
			// ACCESS_KEY_ID
			// SECRET_ACCESS_KEY
			// SECURE

			// sourceFile
			// destinationPath e.g. bucket:filepath/objectname

		case "git":

			fileContent := sthingsBase.ReadFileToVariable(sourceFile)

			gitUser = internal.ValidateGetVaultSecretValue(gitUser, log)
			gitToken = internal.ValidateGetVaultSecretValue(gitToken, log)

			// GET SECRET VALUE
			gitAuth := sthingsCli.CreateGitAuth(gitUser, gitToken)

			if sthingsCli.AddCommitFileToGitRepository(gitRepository, gitBranch, gitAuth, []byte(fileContent), destinationPath, commitMessage) {
				log.Info("PUSH OF FILE ", sourceFile+" SUCCESSFUL")

			} else {
				log.Error("PUSH OF FILE ", sourceFile+" NOT SUCCESSFUL")
			}

		}

	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
	pushCmd.Flags().String("source", "", "source file path")
	pushCmd.Flags().String("destination", "", "destination path")
	pushCmd.Flags().String("target", "git", "push target")
}
