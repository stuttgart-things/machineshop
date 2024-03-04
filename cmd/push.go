/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cmd

import (
	"os"
	"strings"

	"github.com/stuttgart-things/machineshop/internal"
	sthingsCli "github.com/stuttgart-things/sthingsCli"

	sthingsBase "github.com/stuttgart-things/sthingsBase"

	"github.com/spf13/cobra"
)

var (
	commitMessage = "pushed w/ machineShop CLI"
	minioLocation = "us-east-1"
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
		internal.ValidateSourceFile(sourceFile)

		if destinationPath != "" {

			switch target {

			case "minio":

				log.Info("PUSHING TO MINIO S3")
				log.Info("MINIO URL: ", os.Getenv("MINIO_ADDR"))
				log.Info("SOURCE: ", sourceFile)
				log.Info("TARGET: ", destinationPath)

				clientCreated, minioClient := sthingsCli.CreateMinioClient()

				if !clientCreated {
					log.Error("MINIO CLIENT CAN NOT BE CREATED")
					os.Exit(3)
				} else {
					log.Info("MINIO CLIENT CREATED")

					destination := strings.Split(destinationPath, ":")
					bucket := destination[0]
					objectName := destination[1]

					log.Info("BUCKET: ", bucket)
					log.Info("OBJECT: ", objectName)

					sthingsCli.CreateMinioBucket(minioClient, bucket, minioLocation)
					uploaded, fileSize := sthingsCli.UploadObjectToMinioBucket(minioClient, bucket, sourceFile, objectName)

					if uploaded {
						log.Info("SUCCESSFULLY UPLOADED OF SIZE: ", fileSize)
					}
				}

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

		} else {
			log.Error("DESTINATION PATH SEEMS SO BE EMPTY")
		}
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
	pushCmd.Flags().String("source", "", "source file path")
	pushCmd.Flags().String("destination", "", "destination path")
	pushCmd.Flags().String("target", "git", "push target")
}
