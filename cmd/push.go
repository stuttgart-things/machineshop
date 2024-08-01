/*
Copyright © 2024 Patrick Hermann patrick.hermann@sva.de
*/

package cmd

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/stuttgart-things/machineshop/internal"
	sthingsCli "github.com/stuttgart-things/sthingsCli"

	"github.com/spf13/cobra"
)

type Message struct {
	Title     string   `json:"title,omitempty"`     // if empty: info
	Info      string   `json:"info,omitempty"`      // if empty: title
	Severity  string   `json:"severity,omitempty"`  // default: info
	Author    string   `json:"author,omitempty"`    // default: unknown
	Timestamp string   `json:"timestamp,omitempty"` // generate timestamp func
	System    string   `json:"system,omitempty"`    // default: unknown
	Tags      []string `json:"tags,omitempty"`      // empty
}

var (
	contentType = "application/json"
	// url           = "https://homerun.homerun-dev.sthings-vsphere.labul.sva.de/generic"
	minioLocation = "us-east-1"
	colors        = map[string]string{
		"red":    "#FF0000",
		"green":  "#00FF00",
		"blue":   "#0000FF",
		"orange": "#DF813D",
		"purple": "#726EAD",
		// Add more colors as needed
	}
	token = "IhrGeheimerToken"

	homeRunBodyData = `{
		"Title": "{{ .Title }}",
		"Info": "{{ .Info }}",
		"Severity": "{{ .Severity }}",
		"Author": "{{ .Author }}",
		"Timestamp": "{{ .Timestamp }}",
		"System": "{{ .System }}",
		"Tags": "{{ .Tags }}"
	}`
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
		title, _ := cmd.LocalFlags().GetString("title")
		body, _ := cmd.LocalFlags().GetString("info")
		author, _ := cmd.LocalFlags().GetString("author")
		severity, _ := cmd.LocalFlags().GetString("severity")
		system, _ := cmd.LocalFlags().GetString("system")
		destination, _ := cmd.LocalFlags().GetString("destination")
		tags, _ := cmd.Flags().GetStringSlice("tags")

		if destination != "" {

			switch target {

			case "homerun":

				// GETTING TIMESTAMP
				dt := time.Now()
				timestamp := dt.Format("01-02-2006 15:04:05")

				if author == "" {
					author = "machineshop"
				}

				if system == "" {
					system = "machineshop"
				}

				if title == "" || body == "" {
					log.Error("TITILE AND/OR BODY EMPTY. EXITING")
					os.Exit(3)
				}

				log.Info("PUSHING TO HOMERUN")
				messageBody := Message{
					Title:     title,
					Info:      body,
					Severity:  severity,
					Author:    author,
					Timestamp: timestamp,
					System:    system,
					Tags:      tags,
				}

				rendered := RenderBody(homeRunBodyData, messageBody)
				fmt.Println(rendered)

				// CREATE HTTP-Request
				req, err := http.NewRequest("POST", destination, bytes.NewBuffer([]byte(rendered)))
				if err != nil {
					fmt.Println("faiulure at creating requests:", err)
					return
				}

				// ADD HEADER
				req.Header.Set("Content-Type", contentType)
				req.Header.Set("X-Auth-Token", token)

				// CREATE HTTP-Client + SEND REQUEST
				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					fmt.Println("error at sending request:", err)
					return
				}
				defer resp.Body.Close()

				// READ THE ANSWER
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					fmt.Println("error reading answer:", err)
					return
				}

				fmt.Println("answer status:", resp.Status)
				fmt.Println("answer body:", string(body))

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
	pushCmd.Flags().String("title", "", "title of homerun message")
	pushCmd.Flags().String("info", "", "homerun message body")
	pushCmd.Flags().String("severity", "info", "homerun message severity")
	pushCmd.Flags().String("author", "machineShop", "homerun message author")
	pushCmd.Flags().String("system", "", "homerun message system")
	pushCmd.Flags().StringSlice("tags", []string{}, "homerun message tags")
}

func RenderBody(templateData string, object interface{}) string {

	tmpl, err := template.New("template").Parse(templateData)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, object)

	if err != nil {
		fmt.Println(err)
	}

	return buf.String()

}
