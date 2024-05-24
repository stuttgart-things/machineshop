/*
Copyright © 2024 Patrick Hermann patrick.hermann@sva.de
*/

package cmd

import (
	"os"
	"strings"

	"github.com/stuttgart-things/machineshop/internal"
	sthingsCli "github.com/stuttgart-things/sthingsCli"

	"github.com/spf13/cobra"
)

var (
	minioLocation = "us-east-1"
	colors        = map[string]string{
		"red":    "#FF0000",
		"green":  "#00FF00",
		"blue":   "#0000FF",
		"orange": "#DF813D",
		"purple": "#726EAD",
		// Add more colors as needed
	}
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "push things",
	Long:  `push things to external systems`,
	Run: func(cmd *cobra.Command, args []string) {

		// FLAGS
		target, _ := cmd.LocalFlags().GetString("target")
		source, _ := cmd.LocalFlags().GetString("source")
		color, _ := cmd.LocalFlags().GetString("color")
		destination, _ := cmd.LocalFlags().GetString("destination")

		if destination != "" {

			switch target {

			case "teams":

				if colors[color] == "" {
					target = "orange"
				}

				log.Info("PUSHING TO MS TEAMS")
				log.Info("MESSAGE: ", source)
				log.Info("MS TEAMS URL: ", destination)
				log.Info("COLOR: ", color)

				webhook := sthingsCli.MsTeamsWebhook{Title: "Machineshop", Text: source, Color: colors[color], Url: destination}
				sthingsCli.SendWebhookToTeams(webhook)

			case "minio":

				log.Info("PUSHING TO MINIO S3")
				log.Info("MINIO URL: ", os.Getenv("MINIO_ADDR"))
				log.Info("SOURCE: ", source)
				log.Info("TARGET: ", destination)

				// VERIFY IF SOURCE FILE IS EXISTING
				internal.ValidateSourceFile(source)

				clientCreated, minioClient := sthingsCli.CreateMinioClient()

				if !clientCreated {
					log.Error("MINIO CLIENT CAN NOT BE CREATED")
					os.Exit(3)

				} else {
					log.Info("MINIO CLIENT CREATED")

					destination := strings.Split(destination, ":")
					bucket := destination[0]
					objectName := destination[1]

					log.Info("BUCKET: ", bucket)
					log.Info("OBJECT: ", objectName)

					sthingsCli.CreateMinioBucket(minioClient, bucket, minioLocation)
					uploaded, fileSize := sthingsCli.UploadObjectToMinioBucket(minioClient, bucket, source, objectName)

					if uploaded {
						log.Info("SUCCESSFULLY UPLOADED OF SIZE: ", fileSize)
					}
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
	pushCmd.Flags().String("target", "minio", "push target")
	pushCmd.Flags().String("color", "orange", "color for webhook message")
}
